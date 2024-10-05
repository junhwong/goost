package field

import (
	"fmt"
	"strconv"
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

type logger interface {
	Error(...any)
}

var log logger = &plog{}

type plog struct {
}

func (p *plog) Error(args ...any) {
	fmt.Println(args...)
}

// 转换类型. 转换失败将不会改变
func As(f *Field, t Type, layouts []string, loc *time.Location) error {
	// if f.Parent != nil && f.IsColumn() && f.Parent.Type != t {
	// 	return fmt.Errorf("必须与父级类型一致")
	// }
	if f.IsColumn() {
		f.Type = t
		for _, it := range f.Items {
			if err := As(it, t, layouts, loc); err != nil {
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
			return err
		}
		f.SetInt(v)
		return nil
	case UintKind:
		obj := GetPrimitiveValue(f)
		v, err := cast.ToUint64E(obj)
		if err != nil {
			return err
		}
		f.SetUint(v)
		return nil
	case FloatKind:
		obj := GetPrimitiveValue(f)
		v, err := cast.ToFloat64E(obj)
		if err != nil {
			return err
		}
		f.SetFloat(v)
		return nil
	case StringKind:
		switch f.Type {
		case StringKind:
		default:
			obj := GetPrimitiveValue(f)
			v, err := cast.ToStringE(obj)
			if err != nil {
				return err
			}
			f.SetString(v)
		}
		return nil
	case BoolKind:
		obj := GetPrimitiveValue(f)
		v, err := cast.ToBoolE(obj)
		if err != nil {
			return err
		}
		f.SetBool(v)
		return nil
	case BytesKind:
		switch f.Type {
		case StringKind:
			f.SetBytes([]byte(f.GetStringValue()))
			return nil
		case BytesKind:
			return nil
		default:
			panic("todo convert")
		}

	case TimeKind:
		switch f.GetType() {
		case IntKind:
			panic("todo convert")
		case UintKind:
			panic("todo convert")
		case FloatKind:
			panic("todo convert")
		case StringKind: // 字符串转日期
			v, err := ParseTime(f.GetString(), layouts, loc)
			if err != nil {
				return err
			}
			f.SetTime(v)
			return err
		case TimeKind:
			if loc != nil {
				f.SetTime(f.GetTime().In(loc))
			}
			return nil
		default:
			panic("todo convert")
		}
	case DurationKind:
		panic("todo convert:DurationKind")
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
