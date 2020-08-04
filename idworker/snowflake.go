package idworker

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync/atomic"
	"time"
)

// SnowflakeIdLayout 用于定义ID的二进制布局
type SnowflakeIdLayout struct {
	// TimeBits 表示时间的位数
	TimeBits uint64
	// WorkerBits 表示工作节点的位数
	WorkerBits uint64
	// SequenceBits 表示序号的位数
	SequenceBits uint64
}

func (layout *SnowflakeIdLayout) Validate() (idMax, workerMax, sequenceMax, epochMax uint64, err error) {
	timeBits, workerBist, sequenceBits := layout.TimeBits, layout.WorkerBits, layout.SequenceBits
	bits := timeBits + workerBist + sequenceBits
	idMax, workerMax, sequenceMax, epochMax = (uint64(1)<<bits)-1, (uint64(1)<<workerBist)-1, (uint64(1)<<sequenceBits)-1, (uint64(1)<<timeBits)-1
	if timeBits == 0 || sequenceBits == 0 || bits > 63 {
		err = fmt.Errorf("bad bits")
	}
	return
}

// SnowflakeIdBuilder 用于构造一个snowflake算法(Twitter)的 ID 生成器。
//
// 注意：所有参数需服务第一次上线时间确定, 设置后不允许改变，否则ID会错乱。
type SnowflakeIdBuilder struct {
	Layout      SnowflakeIdLayout                             // ID布局结构
	WorkerIDGen func(max uint64) (workerID uint64, err error) // 用于获取 workerID
	TimeGen     func() uint64                                 // 用于获取当前时间戳
}

func (config SnowflakeIdBuilder) Build() (IDGenerator, error) {
	idMax, workerMax, sequenceMax, epochMax, err := config.Layout.Validate()
	if err != nil {
		return nil, err
	}
	sequenceBits, workerBits := config.Layout.SequenceBits, config.Layout.WorkerBits
	var workerID, lastSequence, lastEpoch uint64
	if workerBits > 0 {
		if config.WorkerIDGen == nil {
			return nil, fmt.Errorf("%q cannot be nil", "WorkerIDGen")
		}
		var err error
		workerID, err = config.WorkerIDGen(workerMax)
		if err != nil {
			return nil, err
		}

		if workerID > workerMax {
			err = fmt.Errorf("worker ID is overflow max: %d/%d", workerID, workerMax)
			return nil, err
		}
		fmt.Printf("idworker: init worker-id: %d\n", workerID)
	}

	timeGen := config.TimeGen
	if timeGen == nil {
		return nil, fmt.Errorf("%q cannot be nil", "TimeGen")
	}
	if n := timeGen(); n > epochMax {
		return nil, fmt.Errorf("time is overflow max %d/%d", n, epochMax)
	}

	idGen := func(epoch, seq uint64) ID {
		id := epoch & idMax
		id <<= sequenceBits + workerBits
		id |= workerID << sequenceBits
		id |= seq
		if id >= idMax {
			// id耗尽
			panic(fmt.Errorf("idworker: ID are exhausted. max id is %d", idMax))
		}
		// fmt.Printf("gen: epoch:%d workerID:%d seq:%d\n", epoch, workerID, seq)
		return ID(id) // strconv.FormatUint(id, 10)
	}
	seed := 10
	if n := int(sequenceBits); n < seed {
		seed = n
	}
	return func() ID {
		var cnt int
		for {
			epoch := timeGen()
			// 如果时钟无变化
			if epoch == lastEpoch {
				seq := lastSequence + 1
				if seq > sequenceMax { // 如果当前时钟周期内序列号大于最大值，则等待下一个周期
					// 等待时钟变化
					cnt++
					fmt.Printf("idworker: waiting for clock to change. {now: %d, retry: %d}\n", epoch, cnt) // 监测日志如果发现过多，则表示应该扩容
					time.Sleep(time.Millisecond / 2)                                                        //半毫秒
					continue
				}
				if !atomic.CompareAndSwapUint64(&lastSequence, lastSequence, seq) {
					// 如果更新序列号不成功，等待并继续
					runtime.Gosched()
					continue
				}
				return idGen(epoch, seq)
			}
			// 如果发生时钟回拨
			if epoch < lastEpoch {
				step := time.Millisecond
				for {
					if step > time.Minute {
						// 等待循环15次, 共1分钟左右
						// 1分钟是应为tfd是43秒同步
						// 闰秒可以自动修复 https://baike.baidu.com/item/%E9%97%B0%E7%A7%92/696742?fr=aladdin
						panic(fmt.Errorf("idworker: clock moved backwards. time epoch difference is %d", lastEpoch-epoch))
					}
					cnt++
					fmt.Printf("idworker: waiting for clock to sync. {now: %d, last: %d}\n", epoch, lastEpoch)
					time.Sleep(step)
					epoch = timeGen()
					if epoch >= lastEpoch {
						// 时间恢复
						break
					}
					step <<= 1 // 每次翻倍
				}
			}
			// 如果时钟前进(多数时候)
			if epoch > lastEpoch {
				seq := uint64(rand.Intn(seed)) // 防止时间间隔过长导致尾数全部是0

				if !atomic.CompareAndSwapUint64(&lastEpoch, lastEpoch, epoch) || !atomic.CompareAndSwapUint64(&lastSequence, lastSequence, seq) { // FIXME: 同步存在不一致的情况
					// 如果更新序列号不成功，等待并继续
					runtime.Gosched()
					continue
				}
				return idGen(epoch, seq)
			}
		}
	}, nil
}

// DefaultLayout 符合 web(ecmascript) 数值安全的布局。
// see also: http://www.ecma-international.org/ecma-262/6.0/index.html#sec-ecmascript-language-types-number-type
var DefaultLayout = SnowflakeIdLayout{
	TimeBits:     42, // 默认最大 2109-05-15 15:35:11
	WorkerBits:   5,  // 最多32个节点
	SequenceBits: 6,  // 64-{0-9}
}

// NewIdBuilder 返回IdBuilder
// 自定义示例：
// 	builder:=NewIdBuilder("order")
// 	builder.Layout = SnowflakeIdLayout {
// 		TimeBits:     41,
// 		WorkerBits:   10,
// 		SequenceBits: 12,
// }
// 	builder.TimeGen = func() uint64 {
// 		epoch := UTCMillisecond()
// 		twepoch := uint64(1577808000000) // 2020-01-01 00:00:00
// 		return epoch - twepoch           // 2089-09-07 23:47:35 最大
// 	}
// 	gen,_:=builder.Build()
// 	id:=gen()
func NewIdBuilder(name string, layout SnowflakeIdLayout) SnowflakeIdBuilder {
	return SnowflakeIdBuilder{
		Layout:      layout,
		WorkerIDGen: WorkerWithPodIDOrHostname,
		TimeGen:     UTCMillisecond,
	}
}

func ParseSnowflakeId(layout SnowflakeIdLayout, id ID) (epoch, workerId, sequence uint64, err error) {
	var seqMax uint64
	_, _, seqMax, _, err = layout.Validate()
	if err != nil {
		return
	}
	i := uint64(id)
	shift := layout.WorkerBits + layout.SequenceBits
	epoch = i >> shift
	workerId = (i & (uint64(1)<<shift - 1)) >> layout.SequenceBits
	sequence = i & seqMax
	return
}
