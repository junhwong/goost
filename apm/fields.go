package apm

import "github.com/junhwong/goost/apm/field"

var (
	LevelKey, LevelField                            = field.Int("level")
	MessageKey, Message                             = field.String("message")
	TimeKey, Time                                   = field.Time("time")
	DurationKey, Duration                           = field.Duration("duration") // 执行的持续时间。微秒
	TraceIDKey, TraceID                             = field.String("trace.id")
	TraceErrorKey, TraceError                       = field.Bool("trace.error")
	TraceServiceNameKey, TraceServiceName           = field.String("trace.service_name")
	SpanIDKey, SpanID                               = field.String("span.id")
	SpanNameKey, SpanName                           = field.String("span.name")
	SpanParentIDKey, SpanParentID                   = field.String("span.parent_id")
	SpanStatusCodeKey, SpanStatusCode               = field.String("span.status_code")
	SpanStatusDescriptionKey, SpanStatusDescription = field.String("span.status_description")
	ErrorMethodKey, ErrorMethod                     = field.String("error.method")
	ErrorStackTraceKey, ErrorStackTrace             = field.String("error.stack_trace")
	TracebackCallerKey, TracebackCaller             = field.String("traceback.caller")
	TracebackPathKey, TracebackPath                 = field.String("traceback.path")
	TracebackLineNoKey, TracebackLineNo             = field.Int("traceback.lineno")
)
