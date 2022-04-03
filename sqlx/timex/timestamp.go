package timex

import (
	"database/sql/driver"
)

// Timestamp 表示时间戳。精度: 纳秒(nanoseconds )。
type Timestamp uint64

// Scan implements the Scanner interface.
func (t *Timestamp) Scan(value interface{}) (err error) {

	return nil
}

// Value implements the driver Valuer interface.
func (t Timestamp) Value() (driver.Value, error) {
	return int64(t), nil
}

var ZeroTimestamp Timestamp = 0

// func init() {
// 	t, _ := time.Parse("2006-01-02 15:04:05", "1980-01-01 00:00:00")
// 	ZeroTimestamp = Timestamp(t)
// }
