package apm

import (
	"time"

	"github.com/junhwong/goost/apm/field"
)

// Level type.
type Level = field.Level

type (
	Field  = field.Field
	Fields = field.FieldSet
)

type Entry interface {
	GetLevel() (v Level)
	GetTime() (v time.Time)
	GetMessage() (v string)
	GetFields() Fields
	GetCallerInfo() *CallerInfo
}

type FieldsEntry struct {
	Time       time.Time
	Level      Level
	Fields     Fields
	CallerInfo CallerInfo
}

func (e FieldsEntry) GetLevel() (v Level) {
	return e.Level
}

func (e FieldsEntry) GetTime() (v time.Time) {
	return e.Time
}

func (e FieldsEntry) GetFields() Fields {
	return e.Fields
}

func (e FieldsEntry) GetCallerInfo() *CallerInfo {
	if e.CallerInfo.Ok {
		return &e.CallerInfo
	}
	return nil
}

func (e FieldsEntry) GetMessage() (v string) {
	if f := e.Fields.Get(MessageKey.Name()); f != nil {
		return f.GetStringValue()
	}
	return
}
