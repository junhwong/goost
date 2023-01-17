package unixtime

import (
	"fmt"
	"sync"
	"time"
)

var CST, _ = time.LoadLocation("Asia/Chongqing")

// Unixtime is a Unix time, the number of nanoseconds elapsed
// since January 1, 1970 UTC.
type Unixtime int64

func (t Unixtime) Millis() int64 {
	return int64(t) / 1e6
}

func (t Unixtime) Seconds() int64 {
	return int64(t) / 1e9
}

// From an time transform to `Unixtime` with time zone.
func From(nt time.Time) Unixtime {
	if _, offs := nt.Zone(); offs != 0 {
		nt = nt.Add(time.Second * time.Duration(-offs))
	}
	return Unixtime(nt.UnixNano())
}

var (
	mu   sync.Mutex
	last Unixtime = NowWithoutLock()
)

func Now() Unixtime {
	mu.Lock()
	defer mu.Unlock()

	ut := NowWithoutLock()
	if ut < last {
		offs := ut - last
		if offs <= 5000 {
			time.Sleep(time.Nanosecond * time.Duration(offs<<1)) // 时间偏差大小小于5ms，则等待两倍时间
			ut = NowWithoutLock()
		}
		fmt.Println("", offs, ut, last)
		// 时间回归后继续验证
		if ut <= last {
			panic(fmt.Errorf("clock is moving backwards. rejecting requests until %v.", last)) // 机器时钟发生回拨
		}
	}
	last = ut
	return ut
}

func NowWithoutLock() Unixtime {
	ut := From(time.Now())
	return ut
}

// type UnixTimeMillis Unixtime

// type UnixTimeNano Unixtime
