package field

import (
	"fmt"
)

type FieldValueBuild struct {
	Value    any    `json:"value"`
	FieldRef string `json:"fieldRef"`
	Type     string `json:"type"`
	Nullable bool   `json:"nullable"`
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

func (r *FieldsBuild) Build(src FieldSet) (FieldSet, error) {
	fs := FieldSet{}
	for k, fb := range r.Fields {
		if fb.Value != nil {
			fs.Set(Any(k, fb.Value))
			continue
		}
		if ff := src.Get(fb.FieldRef); ff != nil {
			fs.Set(Clone(ff)) // todo 类型
			continue
		}
		if fb.Nullable {
			continue
		}
		return nil, fmt.Errorf("key 未找到: %v", k)
	}
	return fs, nil
}
