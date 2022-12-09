package apm

import (
	"context"

	"github.com/junhwong/goost/apm/field"
)

type SpanOptionSetter interface {
	SetNameGetter(a func() string)
	SetAttributes(a []field.Field)
	SetCalldepth(a int)
}
type EndSpanOptionSetter interface {
	SetNameGetter(a func() string)
	SetAttributes(a []field.Field)
	SetEndCalls(a []func(Span))
}

type SpanOption interface {
	Apply(target SpanOptionSetter)
}

//	type StartSpanOption interface {
//		applyStartSpanOption(*traceOption)
//	}
type EndSpanOption interface {
	Apply(target EndSpanOptionSetter)
}

type funcSpanOption func(SpanOptionSetter)

func (f funcSpanOption) Apply(target SpanOptionSetter) {
	f(target)
}

func WithName(name string) SpanOption {
	return funcSpanOption(func(target SpanOptionSetter) {
		target.SetNameGetter(func() string {
			return name
		})
	})
}

// 调整日志堆栈记录深度
func WithCallDepth(depth int) SpanOption {
	return funcSpanOption(func(target SpanOptionSetter) {
		target.SetCalldepth(depth)
	})
}

type funcEndSpanOption func(EndSpanOptionSetter)

func (f funcEndSpanOption) Apply(target EndSpanOptionSetter) {
	f(target)
}

// 替换SpanName
func WithReplaceSpanName(getName func() string) EndSpanOption {
	if getName == nil {
		panic("apm: getName cannot be nil")
	}
	return funcEndSpanOption(func(target EndSpanOptionSetter) {
		target.SetNameGetter(getName)
	})
}

// 调整日志堆栈记录深度
func WithEndCall(fn func(Span)) EndSpanOption {
	return funcEndSpanOption(func(target EndSpanOptionSetter) {
		target.SetEndCalls([]func(Span){fn})
	})
}

// //
// func (fn appendFields) applySpanOption(opt *traceOption) {
// 	opt.attrs = append(opt.attrs, fn()...)
// }
// func (fn appendFields) applyEndOption(opt *traceOption) {
// 	opt.attrs = append(opt.attrs, fn()...)
// }

// // Deprecated 已经废弃
// func WithTrimFieldPrefix(prefix ...string) SpanOption {
// 	return &traceOption{delegate: func(target *traceOption) {
// 		target.trimFieldPrefix = prefix
// 	}}
// }

func Start(ctx context.Context, options ...SpanOption) (context.Context, Span) {
	return std.NewSpan(ctx, options...)
}

// // 调整日志堆栈记录深度
// func WithClearup(closer interface{}) *traceOption {
// 	return WithEndCall(func(s Span) {
// 		Close(closer, s)
// 	})
// }
