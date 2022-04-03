package common

import "github.com/junhwong/goost/pkg/field"

var (
	MessageKey, Message           = field.String("message")
	TraceIDKey, TraceID           = field.String("trace.id")
	SpanIDKey, SpanID             = field.String("trace.span.id")
	SpanParentIDKey, SpanParentID = field.String("trace.span.parentid")
	DurationKey, Duration         = field.Duration("duration") // 纳秒
)
