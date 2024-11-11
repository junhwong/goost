package time

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cast"
)

type Duration time.Duration

func (t *Duration) Scan(value interface{}) error {
	switch value := value.(type) {
	case float64:
		return t.Scan(int64(value))
	case int64:
		if value >= 0 {
			*t = Duration(time.Duration(value))
			// *t = Duration{time.Duration(value)}
			return nil
		}
	case []byte:
		return t.Scan(string(value))
	case string:
		if s, err := strconv.Unquote(value); err == nil {
			value = s
		}
		v, err := ParseDuration(value)
		if err != nil {
			return err
		}
		*t = Duration(v)
		return nil
	default:
		return t.Scan(fmt.Sprint(value))
	}
	return fmt.Errorf("不支持的时间戳格式, %T: %v", value, value)
}

func (t Duration) Value() (driver.Value, error) {
	return t, nil
}
func (t Duration) String() string {
	return time.Duration(t).String()
}
func (t Duration) GoString() string {
	return t.String()
}
func (t Duration) MarshalJSON() ([]byte, error) {
	v := t.String()
	return []byte(v), nil
}
func (t *Duration) UnmarshalJSON(b []byte) (err error) {
	return t.Scan(b)
}

// 解析相对于基础时间，返回时间间隔
func ParseMomentDuration(s string) (time.Duration, error) {
	base := time.Now()
	s = strings.ReplaceAll(s, "&nbsp;", " ")
	s = strings.TrimSpace(s)
	switch s {
	case "":
		return 0, nil
	case "刚刚", "1秒前", "一秒前", "just", "just now", " last second", "1 second ago":
		return base.Sub(base.Add(-1 * time.Second)), nil
	case "1分钟前", "1 min ago", "1 minute ago", "last minute", "一分钟前":
		return base.Sub(base.Add(-1 * time.Minute)), nil
	case "1小时前", "1 hour ago", "last hour", "一小时前", "last hour ago":
		return base.Sub(base.Add(-1 * time.Hour)), nil
	case "昨天", "yesterday", "1 day ago", "一天前", "1天前", "last day":
		return base.Sub(base.AddDate(0, 0, -1)), nil
	case "前天", "the day before yesterday":
		return base.Sub(base.AddDate(0, 0, -2)), nil
	case "一周前", "a week ago", "1周前", "1 week ago", "last week":
		return base.Sub(base.AddDate(0, 0, -7)), nil
	case "一个月前", "a month ago", "1个月前", "1 month ago", "last month":
		return base.Sub(base.AddDate(0, -1, 0)), nil
	case "一年前", "a year ago", "1年前", "1 year ago", "last year":
		return base.Sub(base.AddDate(-1, 0, 0)), nil
	}

	if s, ok := strings.CutSuffix(s, "years ago"); ok {
		n, err := cast.ToIntE(s)
		if err != nil {
			return 0, err
		}
		return base.Sub(base.AddDate(-1*n, 0, 0)), nil
	}
	if s, ok := strings.CutSuffix(s, "年前"); ok {
		n, err := cast.ToIntE(s)
		if err != nil {
			return 0, err
		}
		return base.Sub(base.AddDate(-1*n, 0, 0)), nil
	}
	if s, ok := strings.CutSuffix(s, "months ago"); ok {
		n, err := cast.ToIntE(s)
		if err != nil {
			return 0, err
		}
		return base.Sub(base.AddDate(0, -1*n, 0)), nil
	}
	if s, ok := strings.CutSuffix(s, "个月前"); ok {
		n, err := cast.ToIntE(s)
		if err != nil {
			return 0, err
		}
		return base.Sub(base.AddDate(0, -1*n, 0)), nil
	}
	if s, ok := strings.CutSuffix(s, "days ago"); ok {
		n, err := cast.ToIntE(s)
		if err != nil {
			return 0, err
		}
		return base.Sub(base.AddDate(0, 0, -1*n)), nil
	}
	if s, ok := strings.CutSuffix(s, "天前"); ok {
		n, err := cast.ToIntE(s)
		if err != nil {
			return 0, err
		}
		return base.Sub(base.AddDate(0, 0, -1*n)), nil
	}

	if s, ok := strings.CutSuffix(s, "ago"); ok {
		return ParseDuration(s)
	}

	if s, ok := strings.CutSuffix(s, "前"); ok {
		return ParseDuration(s)
	}

	return 0, fmt.Errorf("unsupported format %s", s)
}

func ParseDuration(os string) (time.Duration, error) {
	s := strings.ReplaceAll(os, "&nbsp;", " ")
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}
	if strings.Contains(s, ":") { // 00:00:00
		a := strings.Split(s, ":")
		l := len(a)
		if l < 3 {
			return 0, fmt.Errorf("unsupported format %s", os)
		}
		panic("TODO")
		return 0, nil
	}

	s = strings.ReplaceAll(s, "小时", "h")
	s = strings.ReplaceAll(s, "分钟", "m")
	s = strings.ReplaceAll(s, "秒钟", "s")
	s = strings.ReplaceAll(s, "毫秒", "ms")
	s = strings.ReplaceAll(s, "微妙", "us")
	s = strings.ReplaceAll(s, "纳秒", "ns")
	s = strings.ReplaceAll(s, "μ", "ms")
	s = strings.ReplaceAll(s, "μs", "ms")
	s = strings.ReplaceAll(s, "天", "d")
	s = strings.ReplaceAll(s, "时", "h")
	s = strings.ReplaceAll(s, "分", "m")
	s = strings.ReplaceAll(s, "秒", "s")
	s = strings.ReplaceAll(s, "年", "y")
	s = strings.ReplaceAll(s, "月", "M")
	s = strings.ReplaceAll(s, "周", "w")
	s = strings.TrimSpace(s)

	h := 0
	if arr := strings.SplitN(s, "y", 2); len(arr) == 2 {
		n, err := cast.ToIntE(arr[0])
		if err != nil {
			return 0, err
		}
		h += n * 360 * 24
		s = arr[1]
	}
	if arr := strings.SplitN(s, "M", 2); len(arr) == 2 {
		n, err := cast.ToIntE(arr[0])
		if err != nil {
			return 0, err
		}
		h += n * 30 * 24
		s = arr[1]
	}
	if arr := strings.SplitN(s, "w", 2); len(arr) == 2 {
		n, err := cast.ToIntE(arr[0])
		if err != nil {
			return 0, err
		}
		h += n * 7 * 24
		s = arr[1]
	}
	if arr := strings.SplitN(s, "d", 2); len(arr) == 2 {
		n, err := cast.ToIntE(arr[0])
		if err != nil {
			return 0, err
		}
		h += n * 24
		s = arr[1]
	}
	if arr := strings.SplitN(s, "h", 2); len(arr) == 2 {
		n, err := cast.ToIntE(arr[0])
		if err != nil {
			return 0, err
		}
		h += n
		s = arr[1]
	}
	s = strconv.Itoa(h) + "h" + s
	return time.ParseDuration(s)
}
