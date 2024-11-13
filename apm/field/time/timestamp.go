package time

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// 一个可用于乐观锁的时间戳类型.
//
// 注: 目前的rdb一般不支持纳秒时间戳,该类型序列化为int64的纳秒时间戳
type Timestamp time.Time

func (k Timestamp) MarshalJSON() ([]byte, error) {
	v := time.Time(k).Format(time.RFC3339Nano)
	return []byte(strconv.Quote(v)), nil
}

func (k *Timestamp) UnmarshalJSON(b []byte) (err error) {
	return k.Scan(b)
}

// 实现 sql.Scanner 接口
func (k *Timestamp) Scan(value interface{}) error {
	switch value := value.(type) {
	case []byte:
		var obj any
		if err := json.Unmarshal(value, &obj); err != nil {
			return err
		}
		return k.Scan(obj)
	case string:
		if v, err := strconv.Unquote(value); err == nil {
			value = v
		}
		if n, err := strconv.ParseInt(value, 10, 64); err == nil {
			return k.Scan(n)
		}
		t, err := time.Parse(time.RFC3339Nano, value)
		if err != nil {
			return err
		}
		*k = Timestamp(t)
		return nil
	case int64:
		*k = Timestamp(time.Unix(0, value).Local())
	default:
		return k.Scan(fmt.Sprint(value))
	}
	return nil
}

// 实现 driver.Valuer 接口
func (k Timestamp) Value() (driver.Value, error) {
	return k.Unix(), nil
}

func (k Timestamp) Time() time.Time {
	return time.Time(k)
}

func (k Timestamp) Unix() int64 {
	return time.Time(k).UnixNano()
}

func NowTime() time.Time {
	return time.Now()
}

func Now() Timestamp {
	return Timestamp(time.Now())
}
