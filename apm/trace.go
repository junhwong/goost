package apm

import (
	"context"

	"github.com/junhwong/goost/pkg/field"
)

type SpanOption interface {
	applySpanOption(*traceOption)
}
type StartSpanOption interface {
	applyStartSpanOption(*traceOption)
}
type EndSpanOption interface {
	applyEndOption(*traceOption)
}
type traceOption struct {
	trimFieldPrefix []string
	name            string
	attrs           []field.Field
	delegate        func(*traceOption)
	getName         func() string
	calldepth       int
	endCalls        []func(Span)
}

func (opt *traceOption) applySpanOption(target *traceOption) {
	opt.delegate(target)
}
func (opt *traceOption) applyEndOption(target *traceOption) {
	opt.applySpanOption(target)
}

// func (opt *traceOption) applyStartSpanOption(target *traceOption) {
// 	opt.applyStartSpanOption(target)
// }

func WithName(name string) SpanOption {
	return &traceOption{delegate: func(target *traceOption) {
		target.name = name
	}}
}

// 调整日志堆栈记录深度
func WithCallDepth(depth int) *traceOption {
	return &traceOption{delegate: func(target *traceOption) {
		target.calldepth = depth
	}}
}

// 替换SpanName
func WithReplaceSpanName(getName func() string) EndSpanOption {
	if getName == nil {
		panic("apm: getName cannot be nil")
	}
	return &traceOption{delegate: func(target *traceOption) {
		target.getName = getName
	}}
}

//
func (fn appendFields) applySpanOption(opt *traceOption) {
	opt.attrs = append(opt.attrs, fn()...)
}
func (fn appendFields) applyEndOption(opt *traceOption) {
	opt.attrs = append(opt.attrs, fn()...)
}

// Deprecated 已经废弃
func WithTrimFieldPrefix(prefix ...string) SpanOption {
	return &traceOption{delegate: func(target *traceOption) {
		target.trimFieldPrefix = prefix
	}}
}

func Start(ctx context.Context, options ...SpanOption) (context.Context, Span) {
	return defi.NewSpan(ctx, options...)
}

// 调整日志堆栈记录深度
func WithEndCall(fn func(Span)) *traceOption {
	return &traceOption{delegate: func(target *traceOption) {
		// if target.endCalls == nil {
		// 	target.endCalls = []func(){fn}
		// 	return
		// }
		target.endCalls = append(target.endCalls, fn)
	}}
}

// 调整日志堆栈记录深度
func WithClearup(closer interface{}) *traceOption {
	return WithEndCall(func(s Span) {
		Close(closer, s)
	})
}
