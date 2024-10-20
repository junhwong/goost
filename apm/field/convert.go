package field

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cast"
)

// 从值倒推类型(只能是基本类型)
func InvertType(f *Field) *Field {
	switch f.Type {
	case IntKind:
		if f.IntValue == nil {
			f.Type = InvalidKind
		}
	case UintKind:
		if f.UintValue == nil {
			f.Type = InvalidKind
		}
	case FloatKind:
		if f.FloatValue == nil {
			f.Type = InvalidKind
		}
	case StringKind:
		if f.StringValue == nil {
			f.Type = InvalidKind
		}
	case BoolKind:
		if f.IntValue == nil {
			f.Type = InvalidKind
		}
	case BytesKind:
		if f.BytesValue == nil {
			f.Type = InvalidKind
		}
	case TimeKind:
		if f.UintValue == nil {
			f.Type = InvalidKind
		}
	case DurationKind:
		if f.IntValue == nil {
			f.Type = InvalidKind
		}
	default:
		f.Type = InvalidKind
	}
	if f.Type != InvalidKind {
		return f
	}
	switch {
	case f.BytesValue != nil:
		f.Type = BytesKind
	case f.FloatValue != nil:
		f.Type = FloatKind
	case f.IntValue != nil:
		f.Type = IntKind
	case f.UintValue != nil:
		f.Type = UintKind
	case f.StringValue != nil:
		f.Type = StringKind
	}
	return f
}

// 转换类型. 转换失败将不会改变
func As(f *Field, t Type, layouts []string, loc *time.Location, baseTime time.Time, failToDefault bool) error {
	// if f.Parent != nil && f.IsColumn() && f.Parent.Type != t {
	// 	return fmt.Errorf("必须与父级类型一致")
	// }
	if f.IsColumn() {
		f.Type = t
		for _, it := range f.Items {
			if err := As(it, t, layouts, loc, baseTime, failToDefault); err != nil {
				return err
			}
		}
		return nil
	}

	switch t {
	case IntKind:
		obj := GetPrimitiveValue(f)
		v, err := cast.ToInt64E(obj)
		if err != nil {
			if failToDefault {
				f.SetInt(0)
			}
			return err
		}
		f.SetInt(v)
		return nil
	case UintKind:
		obj := GetPrimitiveValue(f)
		v, err := cast.ToUint64E(obj)
		if err != nil {
			if failToDefault {
				f.SetUint(0)
			}
			return err
		}
		f.SetUint(v)
		return nil
	case FloatKind:
		obj := GetPrimitiveValue(f)
		v, err := cast.ToFloat64E(obj)
		if err != nil {
			if failToDefault {
				f.SetFloat(0)
			}
			return err
		}
		f.SetFloat(v)
		return nil
	case StringKind:
		switch f.Type {
		case StringKind:
			return nil
		default:
			obj := GetPrimitiveValue(f)
			v, err := cast.ToStringE(obj)
			if err != nil {
				if failToDefault {
					f.SetString("")
				}
				return err
			}
			f.SetString(v)
			return nil
		}
	case BoolKind:
		obj := GetPrimitiveValue(f)
		v, err := cast.ToBoolE(obj)
		if err != nil {
			if failToDefault {
				f.SetBool(false)
			}
			return err
		}
		f.SetBool(v)
		return nil
	case BytesKind:
		switch f.Type {
		case StringKind:
			f.SetBytes([]byte(f.GetString()))
			return nil
		case BytesKind:
			return nil
		}
	case TimeKind:
		switch f.GetType() {
		case StringKind: // 字符串转日期
			v, err := ParseTime(f.GetString(), layouts, loc)
			if err != nil {
				if failToDefault {
					f.SetTime(time.Time{})
				}
				return err
			}
			f.SetTime(v)
			return err
		case TimeKind:
			if loc != nil {
				f.SetTime(f.GetTime().In(loc))
			}
			return nil
		case IntKind:
			d := f.GetInt()
			for _, l := range layouts {
				switch strings.ToUpper(l) {
				case "UNIX_MS":
					f.SetTime(time.UnixMilli(d))
					return nil
				case "UNIX_US":
					f.SetTime(time.UnixMicro(d))
					return nil
				case "UNIX_NS":
					f.SetTime(time.Unix(0, d))
					return nil
				case "UNIX":
					f.SetTime(time.Unix(d, 0))
					return nil
				}
			}
			f.SetTime(time.Unix(0, d))
		}
	case DurationKind:
		switch f.GetType() {
		case StringKind:
			s := f.GetString()
			d, err := ParseDuration(s)
			if err != nil {
				d, err = ParseMomentDuration(s)
			}
			if err != nil {
				if failToDefault {
					f.SetDuration(0)
				}
				return fmt.Errorf("invalid duration: %s", s)
			}
			f.SetDuration(d)
			return nil
		case IntKind:
			d := f.GetInt()
			for _, l := range layouts {
				switch strings.ToUpper(l) {
				case "MILLISECONDS", "MS":
					f.SetDuration(time.Duration(d) * time.Millisecond)
					return nil
				case "MICSECONDS", "US":
					f.SetDuration(time.Duration(d) * time.Microsecond)
					return nil
				case "NANOSECONDS", "NS":
					f.SetDuration(time.Duration(d))
					return nil
				case "SECONDS", "S":
					f.SetDuration(time.Duration(d) * time.Second)
					return nil
				case "MINUTES", "M":
					f.SetDuration(time.Duration(d) * time.Minute)
					return nil
				case "HOURS", "H":
					f.SetDuration(time.Duration(d) * time.Hour)
					return nil
				}
			}
			f.SetDuration(time.Duration(d) * time.Microsecond)
			return nil
		case FloatKind:
			v := f.GetFloat()
			for _, l := range layouts {
				switch strings.ToUpper(l) {
				case "MILLISECONDS", "MS":
					v *= float64(time.Millisecond)
					goto LOOP
				case "MICSECONDS", "US":
					v *= float64(time.Microsecond)
					goto LOOP
				case "NANOSECONDS", "NS":
					v *= float64(time.Nanosecond)
					goto LOOP
				case "SECONDS", "S":
					v *= float64(time.Second)
					goto LOOP
				case "MINUTES", "M":
					v *= float64(time.Minute)
					goto LOOP
				case "HOURS", "H":
					v *= float64(time.Hour)
					goto LOOP
				}
			}
			v *= float64(time.Second)
		LOOP:
			f.SetDuration(time.Duration(v))
			return nil
		}
	case IPKind:
		switch f.GetType() {
		case StringKind:
			s := f.GetString()
			ip := net.ParseIP(s)
			if strings.Contains(s, ":") {
				ip = ip.To16()
			} else {
				ip = ip.To4()
			}
			if len(ip) == 0 {
				if failToDefault {
					f.SetIP(net.IP{})
				}
				return fmt.Errorf("invalid ip: %s", s)
			}
			f.SetIP(ip)
			return nil
		}
	case GroupKind:
		switch {
		case f.IsGroup():
			if len(layouts) == 0 {
				return nil
			}
			switch layouts[0] {
			case "RowTable":
				group := false
				array := false
				for _, it := range f.Items {
					if it.IsArray() {
						array = true
					}
					if it.IsGroup() {
						group = true
					}
				}

				if group && array {
					return fmt.Errorf("field:转换为RowTable失败: %v", "列中不能同时混杂group和array")
				}
				if group {
					return nil
				}
				items := make([]*Field, len(f.Items))
				copy(items, f.Items)

				ToRowTable(f, items)
				return nil
			}
		case f.IsArray():
			// todo 目前只是单纯的转换
			ItemsCopy := make([]*Field, len(f.Items))
			copy(ItemsCopy, f.Items)
			f.SetGroup(nil)
			for i, it := range ItemsCopy {
				if it.Name == "" {
					it.Name = strconv.Itoa(i)
				}
				exists := false
				for _, eit := range f.Items {
					if eit.Name == it.Name {
						exists = true
						break
					}
				}
				if exists {
					it.Name += "_" + strconv.Itoa(i)
				}
				f.Set(it)
			}
			return nil
		}
	case ArrayKind:
		switch f.Type {
		case GroupKind:
			if len(layouts) == 0 {
				panic("todo convert:ArrayKind-GroupKind")
			}
			switch layouts[0] {
			case "RowTable":
				group := false
				array := false
				for _, it := range f.Items {
					if it.IsArray() {
						array = true
					}
					if it.IsGroup() {
						group = true
					}
				}

				if group && array {
					return fmt.Errorf("field:转换为RowTable失败: %v", "列中不能同时混杂group和array")
				}
				if group { // ?
					return nil
				}
				items := make([]*Field, len(f.Items))
				copy(items, f.Items)

				ToRowTable(f, items)
				return nil
			}
		}
	}

	panic(fmt.Sprintf("todo convert %v->%v", f.GetType(), t))
}

func ToRowTable(dest *Field, cols []*Field) {
	if len(cols) == 0 {
		log.Error("ToRowTable: empty cols")
		dest.SetArray(nil, false)
		return
	}
	rows := []*Field{}
	colcnt := len(cols)
	rowcnt := len(cols[0].Items) // 行数
	for i := 0; i < rowcnt; i++ {
		row := Make(dest.Name).SetKind(GroupKind, false, false)
		for j := 0; j < colcnt; j++ {
			f := cols[j].Items[i]
			f.Name = cols[j].Name
			row.Set(f)
		}
		rows = append(rows, row)
	}
	dest.SetArray(rows, false)
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
	return time.ParseDuration(strconv.Itoa(h) + "h" + s)
}
