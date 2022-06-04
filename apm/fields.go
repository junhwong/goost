package apm

import "github.com/junhwong/goost/pkg/field"

var (
	LevelKey, Level                                 = field.Int("level")
	MessageKey, Message                             = field.String("message")
	TimeKey, Time                                   = field.Time("time")
	TraceIDKey, TraceID                             = field.String("trace.id")
	TraceErrorKey, TraceError                       = field.Bool("trace.error")
	SpanIDKey, SpanID                               = field.String("span.id")
	SpanNameKey, SpanName                           = field.String("span.name")
	SpanParentIDKey, SpanParentID                   = field.String("span.parent_id")
	SpanStatusCodeKey, spanStatusCode               = field.String("span.status_code")
	SpanStatusDescriptionKey, spanStatusDescription = field.String("span.status_description")
	DurationKey, Duration                           = field.Duration("duration") // 执行的持续时间。微秒
	ErrorMethodKey, ErrorMethod                     = field.String("error.method")
	TracebackCallerKey, TracebackCaller             = field.String("traceback.caller")
	TracebackPathKey, TracebackPath                 = field.String("traceback.path")
	TracebackLineNoKey, TracebackLineNo             = field.Int("traceback.lineno")
	TraceServiceNameKey, TraceServiceName           = field.String("trace.service_name")
)

// 	ErrorStack  = field.String("apm.error.stack_trace")
// 	ErrorCode   = field.String("apm.error.code")
// 	ErrorMethod = field.String("apm.error.method")
