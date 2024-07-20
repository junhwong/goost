package apm

import "github.com/junhwong/goost/apm/field"

var (
	LevelKey, LevelField                            = field.BuildLevel("level")
	MessageKey, Message                             = field.String("message")
	TimeKey, Time                                   = field.Time("__time__")
	DurationKey, Duration                           = field.Duration("duration")
	TraceIDKey, TraceIDField                        = field.String("trace_id")
	ServiceNameKey, ServiceName                     = field.String("service.name")
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
	LogAdapterKey, LogAdapter                       = field.String("log.adapter")
	LogComponentKey, LogComponent                   = field.String("log.component")
)
