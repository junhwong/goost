package apm

import (
	"context"
	"sync"
	"time"

	"github.com/junhwong/goost/apm/field"
	"github.com/junhwong/goost/apm/field/loglevel"
)

// Level type.

type Entry interface {
	GetLevel() (v loglevel.Level)
	GetTime() (v time.Time)
	GetMessage() (v string)
	GetFields() []*field.Field
	GetCallerInfo() *CallerInfo
}

type FieldsEntry struct {
	field.Field
	CallerInfo *CallerInfo
	calldepth  int // 1
	ctx        context.Context
	mu         sync.Mutex
}

func (e *FieldsEntry) GetLevel() (v loglevel.Level) {
	if f := field.GetLast(e.Items, LevelKey.Name()); f != nil {
		return f.GetLevel()
	}
	return loglevel.Unset
}

func (e *FieldsEntry) GetTime() (v time.Time) {
	if f := field.GetLast(e.Items, TimeKey.Name()); f != nil {
		return f.GetTime()
	}
	return time.Time{}
}

func (e *FieldsEntry) GetFields() []*field.Field {
	return e.Items
}

func (e *FieldsEntry) GetCallerInfo() *CallerInfo {
	if e.CallerInfo.Ok {
		return e.CallerInfo
	}
	return nil
}

func (e *FieldsEntry) GetMessage() string {
	if f := field.GetLast(e.Items, MessageKey.Name()); f != nil {
		return f.GetString()
	}
	return ""
}
