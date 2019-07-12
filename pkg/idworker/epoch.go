package idworker

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// EPOCH 表示时间戳起始点。
//
// 警告：中途修改只能调大，否则造成生成的ID大于历史ID。
const EPOCH int64 = 123

var (
	sequence      int64             // 序列号占 12 位,十进制范围是 [ 0, 4095 ]
	lastTimeStamp int64 = timeGen() // 上次的时间戳(毫秒级), 1秒=1000毫秒, 1毫秒=1000微秒,1微秒=1000纳秒
	seedSequence  int64 = 100       // 新序列起始随机范围
	maxSequence   int64 = 4095      // 序列号占 12 位,十进制范围是 [ 0, 4095 ]
	epochMu       sync.Mutex
)

func init() {
	rand.Seed(time.Now().UnixNano())
	sequence = rand.Int63n(seedSequence)
	t, _ := time.Parse(time.RFC3339Nano, "2018-01-01T00:00:00.0+08:00")
	twepoch = t.UnixNano() / 1e6
}

// timeGen 返回当前毫秒时间戳
func timeGen() int64 {
	return time.Now().UnixNano() / 1e6
}

// Next 返回当前毫秒时间戳和一个毫秒内的序列号。
func Next() (epoch int64, sn int64) {
	// FIXME: atomic.LoadInt64和atomic.CompareAndSwapInt64不能保证线程安全
	epochMu.Lock()
	defer epochMu.Unlock()
LOOP:
	for {
		epoch = timeGen()
		last := atomic.LoadInt64(&lastTimeStamp)
		n := atomic.LoadInt64(&sequence)
		switch {
		case epoch < last:
			if offset := epoch - last; offset <= 5 {
				time.Sleep(time.Millisecond * time.Duration(offset<<1)) // 时间偏差大小小于5ms，则等待两倍时间
				// 时间回归后继续
				if timeGen() > last {
					continue LOOP
				}
			}
			panic(fmt.Errorf("clock is moving backwards. rejecting requests until %v.", last)) // 机器时钟发生回拨
		case epoch > last:
			if !atomic.CompareAndSwapInt64(&lastTimeStamp, last, epoch) {
				continue LOOP
			}
			sn = rand.Int63n(seedSequence) // 防止时间间隔太长全是0
		default:
			sn = n + 1
		}

		if sn > maxSequence {
			time.Sleep(time.Millisecond)
			continue
		}
		if atomic.CompareAndSwapInt64(&sequence, n, sn) {
			break
		}
	}

	return
}
