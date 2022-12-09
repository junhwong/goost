package apm

import (
	"github.com/junhwong/goost/apm/field"
	"github.com/spf13/cast"
)

type Entry interface {
	GetLevel() (lvl LogLevel)
	GetFields() Fields
}

type (
	Field       = field.Field
	Fields      = field.Fields
	FieldsEntry field.Fields
)

func (e FieldsEntry) GetLevel() (lvl LogLevel) {
	if e == nil {
		return
	}

	if val := e.GetFields().Get(LevelKey); val != nil {
		switch v := val.(type) {
		case int:
			lvl = LevelFromInt(v)
		default:
			lvl = LevelFromInt(cast.ToInt(val))
		}
	}
	return
}

func (e FieldsEntry) GetFields() Fields {
	return Fields(e)
}
