package common

import "github.com/junhwong/goost/pkg/field"

var (
	Message      = field.String("message")
	TraceID      = field.String("trace.id")
	SpanID       = field.String("trace.span.id")
	SpanParentID = field.String("trace.span.parentid")
	Duration     = field.Duration("duration") // 纳秒
)
