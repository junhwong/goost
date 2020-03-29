package idworker

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync/atomic"
	"time"
)

// SnowflakeIdBuilder 用于构造一个snowflake算法(Twitter)的 ID 生成器。
//
// 注意：所有参数需服务第一次上线时间确定, 设置后不允许改变，否则ID会错乱。
type SnowflakeIdBuilder struct {
	// TimeBits 表示时间的位数
	TimeBits uint64
	// WorkerBits 表示工作节点的位数
	WorkerBits uint64
	// SequenceBits 表示序号的位数
	SequenceBits uint64
	// WorkerIDGen 用于获取 workerID
	WorkerIDGen func(max uint64) (workerID uint64, err error)
	// TimeGen 用于获取当前时间戳。
	//
	// 注意：该方法必须返回UTC时间的时间戳。
	TimeGen func() uint64
}

func (config SnowflakeIdBuilder) Build() (IDGenerator, error) {
	timeBits, workerBist, sequenceBits := config.TimeBits, config.WorkerBits, config.SequenceBits
	bits := timeBits + workerBist + sequenceBits
	idMax, workerMax, sequenceMax, epochMax := (uint64(1)<<bits)-1, (uint64(1)<<workerBist)-1, (uint64(1)<<sequenceBits)-1, (uint64(1)<<timeBits)-1
	if timeBits == 0 || sequenceBits == 0 || bits > 63 {
		return nil, fmt.Errorf("bad bits")
	}
	var workerID, lastSequence, lastEpoch uint64
	if workerBist > 0 {
		if config.WorkerIDGen == nil {
			return nil, fmt.Errorf("`WorkerIDGen` cannot be nil")
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
	}

	timeGen := config.TimeGen
	if timeGen == nil {
		return nil, fmt.Errorf("`TimeGen` cannot be nil")
	}
	if n := timeGen(); n > epochMax {
		return nil, fmt.Errorf("time is overflow max %d/%d", n, epochMax)
	}

	idGen := func(epoch, seq uint64) ID {
		id := epoch & idMax
		id <<= sequenceBits + workerBist
		id |= workerID << sequenceBits
		id |= seq
		if id >= idMax {
			// id耗尽
			panic(fmt.Errorf("idworker: ID are exhausted. max id is %d", idMax))
		}

		return ID(strconv.FormatUint(id, 10))
	}
	seed := 10
	if n := int(sequenceBits); n < seed {
		seed = n
	}
	return func() ID {

		for {
			epoch := timeGen()
			// 如果时钟无变化
			if epoch == lastEpoch {
				seq := lastSequence + 1
				if seq > sequenceMax {
					// 等待时钟变化
					time.Sleep(time.Nanosecond)
					continue
				}
				if !atomic.CompareAndSwapUint64(&lastSequence, lastSequence, seq) {
					// 如果更新序列号不成功，等待并继续
					time.Sleep(time.Nanosecond)
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

				if !atomic.CompareAndSwapUint64(&lastEpoch, lastEpoch, epoch) || !atomic.CompareAndSwapUint64(&lastSequence, lastSequence, seq) {
					// 如果更新序列号不成功，等待并继续
					time.Sleep(time.Nanosecond)
					continue
				}
				return idGen(epoch, seq)
			}
		}
	}, nil
}

// 9007199256903680
// 9007199254740990
// http://www.ecma-international.org/ecma-262/6.0/index.html#sec-ecmascript-language-types-number-type
// 0x1FF FFFF FFFF

func NewShortSnowflakeIdBuilder() SnowflakeIdBuilder {
	return SnowflakeIdBuilder{
		TimeBits:     32,
		WorkerBits:   5,
		SequenceBits: 16,
		WorkerIDGen:  WorkerWithPodIDOrHostname,
		TimeGen: func() uint64 {
			return UTCMillisecond() / 1e3 // 秒，最大 2106-02-07 14:28:15
		},
	}
}

func NewLongSnowflakeIdBuilder() SnowflakeIdBuilder {
	return SnowflakeIdBuilder{
		TimeBits:     41,
		WorkerBits:   10,
		SequenceBits: 12,
		WorkerIDGen:  WorkerWithPodIDOrHostname,
		TimeGen: func() uint64 {
			epoch := UTCMillisecond()
			twepoch := uint64(1577808000000) // 2020-01-01 00:00:00
			return epoch - twepoch           // 2089-09-07 23:47:35 最大
		},
	}
}
