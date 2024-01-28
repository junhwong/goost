package field

import (
	"github.com/junhwong/goost/jsonpath"
)

// // 字段集合
// type FieldSet []*Field

// func (x FieldSet) Len() int           { return len(x) }
// func (x FieldSet) Less(i, j int) bool { return x[i].GetKey() < x[j].GetKey() } // 字典序
// func (x FieldSet) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
// func (x FieldSet) Sort() {
// 	if len(x) == 0 {
// 		return
// 	}
// 	sort.Sort(x)
// }

// func (fs *FieldSet) Set(f *Field) *Field {
// 	f, _ = fs.Put(f)
// 	return f
// }
// func (fs *FieldSet) Put(f *Field) (crt, old *Field) {
// 	crt = f
// 	i := fs.At(f.GetKey())
// 	if i < 0 {
// 		*fs = append(*fs, f)
// 		return
// 	}
// 	tmp := *fs
// 	old = tmp[i]
// 	tmp[i] = f
// 	return
// }

// func (fs *FieldSet) SetWith(f *Field, prt *Field) *Field {
// 	if prt != nil {
// 		// if prt.IsList() {
// 		// 	f.Parent = prt
// 		// 	prt.ItemsValue = append(prt.ItemsValue, f)
// 		// 	return f
// 		// }
// 		// if prt.Type == MapKind {
// 		// 	var ifs FieldSet = prt.ItemsValue
// 		// 	f.Parent = prt
// 		// 	ifs.Set(f)
// 		// 	prt.ItemsValue = ifs
// 		// 	return f
// 		// }
// 		// if prt.Parent != nil {
// 		// 	return fs.SetWith(f, prt.Parent)
// 		// }
// 	}
// 	f, _ = fs.Put(f)
// 	return f
// }

// func (fs FieldSet) get(k string) (int, *Field) {
// 	for i, v := range fs {
// 		if v.GetKey() == k {
// 			return i, v
// 		}
// 	}
// 	return -1, nil
// }

// func (fs FieldSet) Get(k string) *Field {
// 	_, f := fs.get(k)
// 	return f
// }
// func (fs FieldSet) At(k string) int {
// 	i, _ := fs.get(k)
// 	return i
// }

// func (fs *FieldSet) Remove(k string) (f *Field) {
// 	tmp := *fs
// 	start := 0
// 	t := len(tmp)

// LOOP:
// 	for {
// 		for i := start; i < t; i++ {
// 			tf := tmp[i]
// 			if tf.Key == k {
// 				f = tf
// 				start = i + 1
// 				for j := i; j < t-1; j++ { // 将后面的元素提前
// 					tmp[j] = tmp[j+1]
// 				}
// 				t--
// 				continue LOOP
// 			}
// 		}
// 		break
// 	}
// 	if f == nil {
// 		return
// 	}
// 	*fs = tmp[:t]
// 	return f
// }

// // 清除重复
// // TODO: 优化效率
// func (fs *FieldSet) Unique() FieldSet {
// 	if fs == nil {
// 		return nil
// 	}
// 	tmp := FieldSet{}
// 	for _, f := range *fs {
// 		tmp.Set(f)
// 	}
// 	*fs = tmp
// 	return tmp
// }

// func (fs FieldSet) Find(keyOrPath string) (FieldSet, error) {
// 	if f := fs.Get(keyOrPath); f != nil {
// 		return FieldSet{f}, nil
// 	}
// 	seg, err := jsonpath.Parse(keyOrPath)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return doFind(fs, seg, nil), nil
// }

func doFind(fs []*Field, s jsonpath.Segment, n jsonpath.Segment) []*Field {
	switch s.Type() {
	case jsonpath.PathSegment:
		segs := s.(jsonpath.Path)
		if len(segs) == 0 {
			return nil
		}
		r := fs
		for i := 0; i < len(segs); i++ {
			n = nil
			if i+1 < len(segs) {
				n = segs[i+1]
			}
			r = doFind(r, segs[i], n)
		}
		return r
	case jsonpath.MulSegment:
		var r []*Field
		for _, v := range s.(jsonpath.Multiple) {
			for _, v2 := range doFind(fs, v, n) {
				r = append(r, v2)
			}
		}
		return r
	case jsonpath.IndexSegment:
		p := s.(jsonpath.Index)
		if i := int(p); i > 0 && i < len(fs) {
			if n != nil {
				return fs[i].Items
			}
			return fs[i:i]
		}
		// todo 超出索引
		return nil
	case jsonpath.KeySegment, jsonpath.QuoteSegment:
		if f := Get(fs, s.Key()); f != nil {
			if n != nil {
				return f.Items
			}
			return []*Field{f}
		}
		return nil
	case jsonpath.RangeSegment:
		seg := s.(jsonpath.Range)
		i := seg[0]
		if i < 0 {
			i += len(fs)
		}
		if i < 0 {
			return nil // todo 超出索引
		}

		j := seg[0]
		if j < 0 {
			j += len(fs)
		}
		if j < 0 || j < i {
			return nil // todo 超出索引
		}

		if n != nil {
			var r []*Field
			for _, v := range fs[i:j] {
				r = append(r, v.Items...)
			}
			return r
		}

		return fs[i:j]

	case jsonpath.SymbolSegment:
		panic("todo: find SymbolSegment")
	default:
		panic("todo: find ")
	}

}

func Find(fs []*Field, nameOrPath string) ([]*Field, error) {
	seg, err := jsonpath.Parse(nameOrPath)
	if err != nil {
		return nil, err
	}
	return FindWithPath(fs, seg)
}

func FindWithPath(fs []*Field, p jsonpath.Segment) ([]*Field, error) {
	return doFind(fs, p, nil), nil
}

func Get(fs []*Field, name string) *Field {
	for i := len(fs) - 1; i >= 0; i-- {
		if fs[i].Name == name {
			return fs[i]
		}
	}
	return nil
}
