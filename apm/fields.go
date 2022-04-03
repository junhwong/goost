package apm

import "github.com/junhwong/goost/pkg/field"

var (
	LevelKey, _entryLevel               = field.Int("apm.level")
	MessageKey, _entryMessage           = field.String("apm.message")
	TimeKey, _entryTime                 = field.Time("apm.time")
	FailKey, _entryFail                 = field.Bool("apm.trace.error")
	TraceIDKey, _entryTraceID           = field.String("apm.trace.id")
	SpanIDKey, _entrySpanID             = field.String("apm.span.id")
	SpanNameKey, _entrySpanName         = field.String("apm.span.name")
	SpanParentIDKey, _entrySpanParentID = field.String("apm.span.parent_id")
	DurationKey, _entryDuration         = field.Duration("apm.duration") // 执行的持续时间。微秒
	ErrorMethodKey, _entryErrorMethod   = field.String("apm.error.method")
	SourcefileKey, _entrySourcefile     = field.String("apm.sourcefile")
)

// 	ErrorStack  = field.String("apm.error.stack_trace")
// 	ErrorCode   = field.String("apm.error.code")
// 	ErrorMethod = field.String("apm.error.method")
