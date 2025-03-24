package field

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	stdtime "time"

	"github.com/junhwong/goost/apm/field/times"
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
func As(f *Field, target Type, layouts []string, loc *stdtime.Location, baseTime stdtime.Time, failToDefault bool) error {
	switch target {
	case IntKind:
		obj := GetPrimitiveValue(f)
		v, err := cast.ToFloat64E(obj)
		if err != nil {
			if failToDefault {
				f.SetInt(0)
			}
			return err
		}
		i, err := strconv.ParseInt(strconv.FormatFloat(v, 'f', 0, 64), 10, 64)
		if err != nil {
			panic(err)
		}
		f.SetInt(i)
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
		case TimeKind:
			if loc != nil {
				f.SetTime(f.GetTime().In(loc))
			}
			return nil
		case StringKind: // 字符串转日期
			s := f.GetString()
			v, err := times.ParseTime(s, layouts, loc)
			if err == nil && s == "" {
				err = fmt.Errorf("空字符串不能转换为日期")
			}
			if err != nil {
				if failToDefault {
					f.SetTime(stdtime.Time{})
				}
				return err
			}
			f.SetTime(v)
			return err
		case IntKind:
			d := f.GetInt()
			for _, l := range layouts {
				switch strings.ToUpper(l) {
				case "UNIX_MS":
					f.SetTime(stdtime.UnixMilli(d))
					return nil
				case "UNIX_US":
					f.SetTime(stdtime.UnixMicro(d))
					return nil
				case "UNIX_NS":
					f.SetTime(stdtime.Unix(0, d))
					return nil
				case "UNIX":
					f.SetTime(stdtime.Unix(d, 0))
					return nil
				}
			}
			f.SetTime(stdtime.Unix(0, d))
			return nil
		default:
			err := fmt.Errorf("todo convert to %v from %#v", target, f)
			if err != nil {
				if failToDefault {
					f.SetTime(stdtime.Time{})
				}
				return err
			}
			f.SetNull(true)
			return nil
		}
	case DurationKind:
		switch f.GetType() {
		case StringKind:
			s := f.GetString()
			d, err := times.ParseDuration(s)
			if err != nil {
				d, err = times.ParseMomentDuration(s)
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
					f.SetDuration(stdtime.Duration(d) * stdtime.Millisecond)
					return nil
				case "MICSECONDS", "US":
					f.SetDuration(stdtime.Duration(d) * stdtime.Microsecond)
					return nil
				case "NANOSECONDS", "NS":
					f.SetDuration(stdtime.Duration(d))
					return nil
				case "SECONDS", "S":
					f.SetDuration(stdtime.Duration(d) * stdtime.Second)
					return nil
				case "MINUTES", "M":
					f.SetDuration(stdtime.Duration(d) * stdtime.Minute)
					return nil
				case "HOURS", "H":
					f.SetDuration(stdtime.Duration(d) * stdtime.Hour)
					return nil
				}
			}
			f.SetDuration(stdtime.Duration(d) * stdtime.Microsecond)
			return nil
		case FloatKind:
			v := f.GetFloat()
			for _, l := range layouts {
				switch strings.ToUpper(l) {
				case "MILLISECONDS", "MS":
					v *= float64(stdtime.Millisecond)
					goto LOOP
				case "MICSECONDS", "US":
					v *= float64(stdtime.Microsecond)
					goto LOOP
				case "NANOSECONDS", "NS":
					v *= float64(stdtime.Nanosecond)
					goto LOOP
				case "SECONDS", "S":
					v *= float64(stdtime.Second)
					goto LOOP
				case "MINUTES", "M":
					v *= float64(stdtime.Minute)
					goto LOOP
				case "HOURS", "H":
					v *= float64(stdtime.Hour)
					goto LOOP
				}
			}
			v *= float64(stdtime.Second)
		LOOP:
			f.SetDuration(stdtime.Duration(v))
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
		if f.IsGroup() {
			var l string
			if len(layouts) != 0 { // todo 多好个layout
				l = layouts[0]
			}
			switch l {
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
			case "":
				f.Type = target
				// todo 设置名称
				return nil
			default:
				return fmt.Errorf("field:转换为RowTable失败: %v", "不支持的field类型")
			}
		}
		if f.IsArray() {
			// todo 处理 layouts
			return nil
		}
		// todo 处理 layouts
		if f.IsNull() || (f.Type == StringKind && f.GetString() == "") {
			f.SetArray(nil)
			return nil
		}
		c := Clone(f)
		f.SetArray([]*Field{c})
		return nil
	}

	panic(fmt.Sprintf("todo convert to %v from %#v", target, f))
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
