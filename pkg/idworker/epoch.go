package idworker

import (
	"math/rand"
	"sync"
	"time"

	"github.com/junhwong/goost/pkg/unixtime"
)

// TWEPOCH_TIME 表示时间戳起始点。
//
// 注意：可通过编译时修改该值。
// 中途修改只能调大，否则造成生成的ID大于历史ID。
var TWEPOCH_TIME string = "2018-01-01T00:00:00.0+08:00"

var (
	lastSequence int64                                      // 序列号占 12 位,十进制范围是 [ 0, 4095 ]
	lastEpoch    int64 = unixtime.NowWithoutLock().Millis() // 上次的时间戳(毫秒级), 1秒=1000毫秒, 1毫秒=1000微秒,1微秒=1000纳秒
	seedSequence int64 = 100                                // 新序列起始随机范围
	maxSequence  int64 = 4095                               // 序列号占 12 位,十进制范围是 [ 0, 4095 ]
	epochMu      sync.Mutex
)

func init() {
	rand.Seed(time.Now().UnixNano())
	lastSequence = rand.Int63n(seedSequence)
	t, _ := time.Parse(time.RFC3339Nano, TWEPOCH_TIME)
	twepoch = unixtime.From(t).Millis()
}

// Next 返回当前毫秒时间戳和一个毫秒内的序列号。
func Next() (epoch int64, sn int64) {

	// FIXME: atomic.LoadInt64和atomic.CompareAndSwapInt64 不能保证线程安全
	epochMu.Lock()
	defer epochMu.Unlock()

	epoch = unixtime.Now().Millis()

	if epoch == lastEpoch {
		sn = lastSequence + 1
	} else {
		sn = rand.Int63n(seedSequence) // 防止时间间隔太长全是0
	}

	if sn > maxSequence {
		time.Sleep(time.Millisecond) // FIXME: 锁占用
		return Next()
	}

	lastSequence = sn
	lastEpoch = epoch

	return
}
