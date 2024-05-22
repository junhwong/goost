package snowflakeid

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/junhwong/goost/apm"
)

// Builder 用于构造一个 snowflake 算法的唯一ID生成器。
type Builder struct {
	// ID布局结构.
	//
	// 注意: Layout 在应用上线前应当确定并不再对其修改, 否则可能造成ID混乱。
	Layout      Layout
	Period      time.Duration                                                            // 时间周期
	NoPanic     bool                                                                     // 出错或不合法时是否panic
	UseRandSeed bool                                                                     // 使用随机起始序号
	WorkerGen   func(topologyMax, workerMax uint64) (topology, worker uint64, err error) // 用于获取 workerID
	TimeGen     func(context.Context) int64                                              // 用于获取当前时间戳
	Log         apm.Logger
}

func (b *Builder) Build() (IDGen, error) {
	idMax, epochMax, topologyMax, workerMax, sequenceMax, idShift, topologyShift, err := b.Layout.Validate()
	if err != nil {
		return nil, err
	}

	log := b.Log
	if log == nil {
		log = apm.Default()
	}

	sequenceBits := int64(b.Layout.SequenceBits)

	var topology, worker int64

	if topologyMax > 0 || workerMax > 0 {
		if b.WorkerGen == nil {
			return nil, fmt.Errorf("%q cannot be nil", "WorkerIDGen")
		}
		d, w, err := b.WorkerGen(uint64(topologyMax), uint64(workerMax))
		if err != nil {
			return nil, err
		}
		topology = int64(d)
		if topology < 0 || topology > topologyMax {
			err = fmt.Errorf("datacenter ID is overflow max: %d/%d", d, topologyMax)
			return nil, err
		}
		worker = int64(w)
		if worker < 0 || worker > workerMax {
			err = fmt.Errorf("worker ID is overflow max: %d/%d", w, workerMax)
			return nil, err
		}

		// fmt.Printf("idworker: init worker-id: %d\n", w)
	}

	startEpoch := int64(b.Layout.Epoch)
	if startEpoch < 0 { // uint转换超过int大小变为负数
		err = fmt.Errorf("Epoch is too big %d", b.Layout.Epoch)
		return nil, err
	}

	timeGen := b.TimeGen
	if timeGen == nil {
		return nil, fmt.Errorf("%q cannot be nil", "TimeGen")
	}
	if n := timeGen(context.TODO()); n <= startEpoch {
		return nil, fmt.Errorf("time is overflow start %d/%d", n, startEpoch)
	} else if n > epochMax {
		return nil, fmt.Errorf("time is overflow max %d/%d", n, epochMax)
	}

	noPanic := b.NoPanic
	useRand := b.UseRandSeed

	idGen := func(epoch int64, seq int64) int64 {
		id := epoch - startEpoch
		if id < 1 || id > epochMax {
			err := fmt.Errorf("invalid epoch %v", epoch)
			if !noPanic {
				panic(err)
			}
			return -1
		}
		id &= idMax
		id <<= idShift
		id |= topology << topologyShift
		id |= worker << sequenceBits
		id |= seq

		if id >= idMax {
			err := fmt.Errorf("ID are exhausted. max id is %d", idMax)
			if !noPanic {
				panic(err)
			}
			log.Error(err)
			return -1

		}
		// fmt.Printf("gen: epoch:%d workerID:%d seq:%d\n", epoch, workerID, seq)
		return id
	}

	seed := int64(99)
	if seed >= sequenceMax {
		seed = sequenceMax
	}
	nextSeq := func() int64 {
		if !useRand {
			return -1
		}
		for {
			seq := rand.Int63n(seed)
			if seq < sequenceMax {
				return seq
			}
		}
	}

	period := b.Period
	if period <= 0 {
		period = time.Millisecond
	}
	if period > time.Minute {
		period = time.Minute
	}
	period /= 2 // 1半 减少等待

	lastSequence := nextSeq()

	var lastEpoch int64
	var mu sync.RWMutex

	return func(ctx context.Context) int64 {
		var cnt int
		for {
			if cnt > 60 { // TODO 传入尝试次数
				err := fmt.Errorf("clock moved backwards. time epoch difference is %d", lastEpoch)
				if !noPanic {
					panic(err)
				}
				log.Error(err)
				return -1
			}

			epoch := timeGen(ctx)                       // 复杂环境
			if epoch-atomic.LoadInt64(&lastEpoch) < 0 { // 时钟回拨或并发调用
				cnt++
				time.Sleep(period)
				continue
			}

			mu.Lock()

			if epoch-lastEpoch < 0 {
				cnt++
				mu.Unlock()
				time.Sleep(period)
				continue
			}

			if epoch != lastEpoch {
				atomic.SwapInt64(&lastEpoch, epoch)
				lastSequence = nextSeq()
			}

			seq := lastSequence + 1
			if seq > sequenceMax { // 如果当前时钟周期内序列号大于最大值，则等待下一个周期
				log.Debug("waiting for clock to change. retries: ", cnt) // 监测日志如果发现过多，则表示应该扩容
				mu.Unlock()
				time.Sleep(period)
				continue
			}
			lastSequence = seq
			mu.Unlock()

			return idGen(epoch, seq)
		}
	}, nil
}
