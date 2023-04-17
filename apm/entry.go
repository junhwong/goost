package apm

import (
	"context"
	"sync"
	"time"

	"github.com/junhwong/goost/apm/field"
)

// Level type.
type Level = field.Level

type Entry interface {
	GetLevel() (v field.Level)
	GetTime() (v time.Time)
	GetMessage() (v string)
	GetFields() field.FieldSet
	GetCallerInfo() *CallerInfo
}

type FieldsEntry struct {
	Time       time.Time
	Level      field.Level
	Fields     field.FieldSet
	CallerInfo CallerInfo
	calldepth  int // 1
	ctx        context.Context
	mu         sync.Mutex
}

func (e *FieldsEntry) GetLevel() (v Level) {
	return e.Level
}

func (e *FieldsEntry) GetTime() (v time.Time) {
	return e.Time
}

func (e *FieldsEntry) GetFields() field.FieldSet {
	return e.Fields
}

func (e *FieldsEntry) GetCallerInfo() *CallerInfo {
	if e.CallerInfo.Ok {
		return &e.CallerInfo
	}
	return nil
}

func (e *FieldsEntry) GetMessage() (v string) {
	if f := e.Fields.Get(MessageKey.Name()); f != nil {
		return f.GetStringValue()
	}
	return
}
