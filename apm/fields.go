package apm

import "github.com/junhwong/goost/pkg/field"

var (
	LevelKey, _entryLevel               = field.Int("level")
	MessageKey, _entryMessage           = field.String("message")
	TimeKey, _entryTime                 = field.Time("time")
	TraceIDKey, _entryTraceID           = field.String("trace.id")
	TraceErrorKey, _entryTraceError     = field.Bool("trace.error")
	SpanIDKey, _entrySpanID             = field.String("span.id")
	SpanNameKey, _entrySpanName         = field.String("span.name")
	SpanParentIDKey, _entrySpanParentID = field.String("span.parent_id")
	DurationKey, _entryDuration         = field.Duration("duration") // 执行的持续时间。微秒
	ErrorMethodKey, _entryErrorMethod   = field.String("error.method")
	SourcefileKey, _entrySourcefile     = field.String("sourcefile")
)

// 	ErrorStack  = field.String("apm.error.stack_trace")
// 	ErrorCode   = field.String("apm.error.code")
// 	ErrorMethod = field.String("apm.error.method")
