package timestamp

import (
	"time"
)

// Timestamp is a Unix time, the number of nanoseconds elapsed
// since January 1, 1970 UTC.
type Timestamp time.Duration

var CST, _ = time.LoadLocation("Asia/Chongqing")

const nanoseconds = int64(time.Second)

// 注意：这是标准UTC时间，GO在format时会按时区调整显示。
// 如：当前local为 CST 时会差2个 UTC+8 的偏移。
func (ts Timestamp) UTC() time.Time {
	nano := int64(ts)
	sec := nano / nanoseconds
	nsec := nano - (sec * nanoseconds)
	return time.Unix(sec, nsec).In(time.UTC)
}

// CST 返回时间戳的中国标准时间。
func (ts Timestamp) CST() time.Time {
	nt := ts.UTC()
	return nt.In(CST).Add(time.Hour * 8)
}

// Local 返回当前时区的时间。
func (ts Timestamp) Local() time.Time {
	nt := ts.UTC().In(time.Local)
	if _, offs := nt.Zone(); offs != 0 {
		return nt.Add(time.Second * time.Duration(offs))
	}
	return nt
}

func (ts Timestamp) Format(layout string) string {
	nt := ts.Local()
	return nt.Format(layout)
}

func (ts Timestamp) String() string {
	nt := ts.Local()
	return nt.String()
}

func Now() Timestamp {
	return From(time.Now())
}

func From(nt time.Time) Timestamp {
	if _, offs := nt.Zone(); offs != 0 {
		nt = nt.Add(time.Second * time.Duration(-offs))
	}
	return Timestamp(nt.UnixNano())
}
