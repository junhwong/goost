package apm

import (
	"time"

	"github.com/junhwong/goost/apm/field"
	"github.com/spf13/cast"
)

type Entry interface {
	GetLevel() (v Level)
	GetTime() (v time.Time)
	GetMessage() (v string)
	GetFields() Fields
}

type (
	Field       = field.Field
	Fields      = []Field
	FieldsEntry field.Fields
)

func (e FieldsEntry) GetLevel() (v Level) {
	if len(e) == 0 {
		return
	}
	switch a := field.Fields(e).Get(LevelKey).(type) {
	case int:
		v = LevelFromInt(a)
	default:
		v = LevelFromInt(cast.ToInt(a))
	}
	return
}

func (e FieldsEntry) GetTime() (v time.Time) {
	if len(e) == 0 {
		return
	}
	switch a := field.Fields(e).Get(TimeKey).(type) {
	case time.Time:
		v = a
	default:
	}
	return
}

func (e FieldsEntry) GetMessage() (v string) {
	if len(e) == 0 {
		return
	}
	switch a := field.Fields(e).Get(MessageKey).(type) {
	case string:
		v = a
	default:
	}
	return
}

func (e FieldsEntry) GetFields() Fields {
	if len(e) == 0 {
		return nil
	}
	return field.Fields(e).List()
}
