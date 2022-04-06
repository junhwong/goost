package apm

import (
	"context"

	"github.com/junhwong/goost/pkg/field"
)

type Option interface {
	apply(*traceOption)
}
type StartOption interface {
	applyStart(*traceOption)
}
type EndOption interface {
	applyEnd(*traceOption)
}
type traceOption struct {
	trimFieldPrefix []string
	name            string
	attrs           []field.Field
	delegate        func(*traceOption)
	getName         func() string
	calldepth       int
}

func (opt *traceOption) apply(target *traceOption) {
	opt.delegate(target)
}
func (opt *traceOption) applyStart(target *traceOption) {
	opt.apply(target)
}
func (opt *traceOption) applyEnd(target *traceOption) {
	opt.apply(target)
}
func WithName(name string) Option {
	return &traceOption{delegate: func(target *traceOption) {
		target.name = name
	}}
}
func WithCalldepth(depth int) *traceOption {
	return &traceOption{delegate: func(target *traceOption) {
		target.calldepth = depth
	}}
}

// 替换SpanName
func WithReplaceSpanName(getName func() string) EndOption {
	if getName == nil {
		panic("apm: getName cannot be nil")
	}
	return &traceOption{delegate: func(target *traceOption) {
		target.getName = getName
	}}
}
func WithFields(fs ...field.Field) *traceOption {
	return &traceOption{delegate: func(target *traceOption) {
		target.attrs = append(target.attrs, fs...)
	}}
}

func WithTrimFieldPrefix(prefix ...string) Option {
	return &traceOption{delegate: func(target *traceOption) {
		target.trimFieldPrefix = prefix
	}}
}

// func FromContext(ctx context.Context) LoggerInterface {
// 	span, _ := ctx.Value(spanInContextKey).(*Span)
// 	if span != nil && span.logger != nil {
// 		return span
// 	}
// 	panic("apm.FromContext: span not found in context")
// }

func Start(ctx context.Context, options ...Option) (context.Context, Span) {
	return std.NewSpan(ctx, options...)
}

// // SpanFromConext 返回与ctx关联的SpanContext，如果未找到则创建一个新的对象。
// func SpanFromConext(ctx context.Context) (span SpanContext) {
// 	if ctx == nil {
// 		ctx = context.TODO()
// 	}
// 	if prent, ok := ctx.Value(spanInContextKey).(*Span); ok && prent != nil {
// 		span.TranceID = prent.TranceID
// 		span.SpanParentID = prent.SpanParentID
// 		span.Name = prent.Name
// 		span.SpanID = prent.SpanID
// 	} else {
// 		span.SpanID = newSpanId()
// 		span.TranceID = TraceIDFromContext(ctx)
// 		span.first = true
// 	}
// 	return
// }
