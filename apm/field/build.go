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
		if k == "" || k == fb.FieldRef {
			return fmt.Errorf("键非法: %v", k)
		}
		if fb.Value == nil {
			if fb.FieldRef == "" {
				return fmt.Errorf("键 %v 的值或引用不能为空", k)
			}
		}
	}

	if r.Fields == nil {
		r.Fields = map[string]FieldValueBuild{}
	}
	return nil
}

func (r *FieldsBuild) Build(src []*Field) ([]*Field, error) {
	fs := []*Field{}
	for k, fb := range r.Fields {
		if fb.Value != nil {
			fs = append(fs, Any(k, fb.Value))
			continue
		}
		if ff := GetLast(src, fb.FieldRef); ff != nil {
			rf := Clone(ff)
			rf.Name = k
			fs = append(fs, rf)
			continue
		}
		if fb.Nullable {
			continue
		}
		return nil, fmt.Errorf("key %s 未找到引用 %v", k, fb.FieldRef)
	}
	return fs, nil
}
