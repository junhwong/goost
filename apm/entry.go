package apm

import (
	"time"

	"github.com/junhwong/goost/apm/field"
)

// Level type.
type Level = field.Level

type (
	Field  = field.Field
	Fields = []*Field
)

type Entry interface {
	GetLevel() (v Level)
	GetTime() (v time.Time)
	GetMessage() (v string)
	GetFields() []*Field
	GetCallerInfo() *CallerInfo
	Lookup(key string) (found []*Field)
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
	for _, l := range e.GetFields() {
		if l != nil && l.GetKey() == MessageKey.Name() {
			return l.GetStringValue()
		}
	}
	return
}
func (e FieldsEntry) Lookup(key string) (found []*Field) {
	for _, l := range e.GetFields() {
		if l != nil && l.GetKey() == key && l.GetType() != field.InvalidKind {
			found = append(found, l)
		}
	}
	return
}
