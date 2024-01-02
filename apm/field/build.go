package field

import "fmt"

type FieldValueBuild struct {
	Value    any    `json:"value"`
	FieldRef string `json:"fieldRef"`
	Type     string `json:"type"`
}

type FieldsBuild struct {
	Fields map[string]FieldValueBuild `json:"fields"`
}

func (r *FieldsBuild) Init() error {
	for k, fb := range r.Fields {
		if len(k) == 0 || k == fb.FieldRef {
			return fmt.Errorf("键非法: %v", k)
		}
		if fb.Value == nil || fb.Value == "" {
			if len(fb.FieldRef) == 0 {
				return fmt.Errorf("键 %v 的值或引用不能为空", k)
			}
			fb.Value = nil
		}
	}
	if r.Fields == nil {
		r.Fields = map[string]FieldValueBuild{}
	}
	return nil
}

func (r *FieldsBuild) Build(f *Field, src FieldSet) (FieldSet, error) {
	fs := FieldSet{}
	for k, fb := range r.Fields {
		if fb.Value == nil {
			if f != nil && fb.FieldRef == f.GetKey() {
				fb.Value = GetObject(f)
			} else if ff := src.Get(fb.FieldRef); ff != nil {
				fb.Value = GetObject(ff)
			}
		}
		if fb.Value == nil {
			continue
		}
		fs.Set(Any(k, fb.Value))
	}
	return fs, nil
}
