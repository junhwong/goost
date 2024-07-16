package apm

import (
	"sync"
	"time"

	"github.com/junhwong/goost/apm/field"
	"github.com/junhwong/goost/apm/field/loglevel"
)

type factoryEntry struct {
	*field.Field
	// CallerInfo *CallerInfo
	calldepth int // 1
	mu        sync.Mutex
}

// func (e *Entry) GetLevel() (v loglevel.Level) {
// 	if f := field.GetLast(e.Items, LevelKey.Name()); f != nil {
// 		return f.GetLevel()
// 	}
// 	return loglevel.Unset
// }

// func (e *Entry) GetTime() (v time.Time) {
// 	if f := field.GetLast(e.Items, TimeKey.Name()); f != nil {
// 		return f.GetTime()
// 	}
// 	return time.Time{}
// }

// func (e *Entry) GetFields() []*field.Field {
// 	return e.Items
// }

// func (e *Entry) GetCallerInfo() *CallerInfo {
// 	if e.CallerInfo.Ok {
// 		return e.CallerInfo
// 	}
// 	return nil
// }

// func (e *Entry) GetMessage() string {
// 	if f := field.GetLast(e.Items, MessageKey.Name()); f != nil {
// 		return f.GetString()
// 	}
// 	return ""
// }

func getTime(e *field.Field) (v time.Time) {
	if f := field.GetLast(e.Items, TimeKey.Name()); f != nil {
		return f.GetTime()
	}
	return time.Time{}
}
func GetLevel(e *field.Field) (v loglevel.Level) {
	if f := field.GetLast(e.Items, LevelKey.Name()); f != nil {
		return f.GetLevel()
	}
	return loglevel.Unset
}
func getMessage(e *field.Field) string {
	if f := field.GetLast(e.Items, MessageKey.Name()); f != nil {
		return f.GetString()
	}
	return ""
}
