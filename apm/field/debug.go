package field

import "fmt"

func toString(f *Field, ident int) string {
	var v string
	var num string
	is := func() string {
		s := ""
		for i := 0; i < ident; i++ {
			s += "\t"
		}
		return s
	}()

	switch {
	case f.IsArray():
		num = fmt.Sprintf(" len: %v", len(f.Items))
		v = "[\n"
		for i, f2 := range f.Items {
			if i != 0 {
				v += ",\n"
			}
			v += toString(f2, ident+1)
		}
		v += "\n"
		v += is + "]"
	case f.IsGroup():
		v = "{\n"
		for i, f2 := range f.Items {
			if i != 0 {
				v += ",\n"
			}
			v += toString(f2, ident+1)
		}
		v += "\n"
		v += is + "}"
	case f.Type == BytesKind:
		v = is + "<bytes>"
	default:
		v = fmt.Sprintf("%v", GetValue(f))
	}

	return fmt.Sprintf("%vField(Name:%v type: %v %v value: %v)", is, f.Name, f.Type, num, v)
}
