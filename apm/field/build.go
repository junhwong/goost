package field

import (
	"fmt"
)

type FieldValueBuild struct {
	Value    any               `json:"value"`
	FieldRef string            `json:"ref"`
	Group    FieldsBuildMap    `json:"group"`
	Array    []FieldValueBuild `json:"array"`
	Type     string            `json:"type"`
	Nullable bool              `json:"nullable"`
}

func (fb *FieldValueBuild) Init(k string) error {
	if k == "" {
		return fmt.Errorf("field name is empty")
	}
	if fb.FieldRef != "" {
		return nil
	}
	if fb.Group != nil {
		if err := fb.Group.Init(); err != nil {
			return err
		}
		return nil
	}
	if fb.Array != nil {
		for i := range fb.Array {
			if err := fb.Array[i].Init(k); err != nil {
				return err
			}
		}
		return nil
	}
	if fb.Value != nil {
		return nil
	}
	fb.FieldRef = k
	return nil
}
func (vb *FieldValueBuild) Build(k string, src ...[]*Field) (*Field, error) {
	if vb.Value != nil {
		return Any(k, vb.Value), nil
	}
	if vb.Group != nil {
		gfs, err := vb.Group.Build(src...)
		if err != nil {
			return nil, err
		}
		return Make(k).SetGroup(gfs), nil
	}
	if vb.Array != nil {
		var tmp []*Field
		for _, b := range vb.Array {
			f, err := b.Build(k, src...)
			if err != nil {
				return nil, err
			}
			tmp = append(tmp, f)
		}

		return Make(k).SetArray(tmp), nil
	}
	for _, s := range src {
		if ff := GetLast(s, vb.FieldRef); ff != nil {
			rf := Clone(ff)
			rf.Name = k
			return rf, nil
		}
	}

	if vb.Nullable {
		return nil, nil
	}
	return nil, fmt.Errorf("key %s 未找到引用 %v", k, vb.FieldRef)
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

type FieldsBuildMap map[string]*FieldValueBuild

func (r *FieldsBuildMap) Init() error {
	rr := *r
	for k, fb := range rr {
		if err := fb.Init(k); err != nil {
			return fmt.Errorf("键 %v 初始化失败: %v", k, err)
		}
	}

	if rr == nil {
		rr = FieldsBuildMap{}
	}
	*r = rr
	return nil
}

func (r FieldsBuildMap) Build(src ...[]*Field) ([]*Field, error) {
	fs := []*Field{}
	for k, fb := range r {
		f, err := fb.Build(k, src...)
		if err != nil {
			return nil, err
		}
		if f != nil {
			fs = append(fs, f)
		}
	}
	return fs, nil
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
