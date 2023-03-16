package apm

import (
	"time"

	"github.com/junhwong/goost/apm/field"
)

type (
	Field  = field.Field
	Fields = []*Field
)
type Entry interface {
	GetLevel() (v Level)
	GetTime() (v time.Time)
	GetMessage() (v string)
	GetLabels() []*Field
	GetCallerInfo() *CallerInfo
	Lookup(key string) (found []*Field)
}

type FieldsEntry struct {
	Level      Level
	Time       time.Time
	Labels     Fields
	CallerInfo CallerInfo
}

func (e FieldsEntry) GetLevel() (v Level) {
	return e.Level
}

func (e FieldsEntry) GetTime() (v time.Time) {
	return e.Time
}
func (e FieldsEntry) GetLabels() Fields {
	return e.Labels
}
func (e FieldsEntry) GetCallerInfo() *CallerInfo {
	if e.CallerInfo.Ok {
		return &e.CallerInfo
	}
	return nil
}

func (e FieldsEntry) GetMessage() (v string) {
	for _, l := range e.GetLabels() {
		if l != nil && l.GetKey() == MessageKey.Name() {
			return l.GetStringValue()
		}
	}
	return
}
func (e FieldsEntry) Lookup(key string) (found []*Field) {
	for _, l := range e.GetLabels() {
		if l != nil && l.GetKey() == key && l.GetType() != field.InvalidKind {
			found = append(found, l)
		}
	}
	return
}
