package apm

import (
	"context"
	"fmt"

	"github.com/junhwong/goost/apm/field"
)

type SpanOption interface {
	applySpanOption(target *spanImpl)
}

type EndSpanOption interface {
	applyEndSpanOption(target *spanImpl)
}

type funcSpanOption func(target *spanImpl)

func (f funcSpanOption) applySpanOption(target *spanImpl) {
	if f == nil {
		return
	}
	f(target)
}

type callDepthProperty interface {
	SetCalldepth(v int)
	GetCalldepth() int
}
type callDepthOption func(target callDepthProperty)

// 实现 SpanOption 接口
func (f callDepthOption) applySpanOption(target *spanImpl) {
	if f == nil {
		return
	}
	f(target)
}
func (f callDepthOption) applyWithOption(target *factoryEntry) {
	if f == nil {
		return
	}
	f(target)
}

// 调整日志堆栈记录深度
func WithCallDepth(depth int) callDepthOption {
	return callDepthOption(func(target callDepthProperty) {
		target.SetCalldepth(depth)
	})
}

// 在当前日志堆栈记录深度上增加指定值
func WithCallDepthAdd(depth int) callDepthOption {
	return callDepthOption(func(target callDepthProperty) {
		target.SetCalldepth(depth + target.GetCalldepth())
	})
}

type fieldsSetter interface {
	SetAttributes(a ...*field.Field)
}
type withFieldsOption func(target fieldsSetter)

// 实现 SpanOption 接口
func (f withFieldsOption) applySpanOption(target *spanImpl) {
	if f == nil {
		return
	}
	f(target)
}

// impl WithOption
func (f withFieldsOption) applyWithOption(target *factoryEntry) {
	if f == nil {
		return
	}
	f(target)
}

// 设置字段
func WithFields(fs ...*field.Field) withFieldsOption {
	return withFieldsOption(func(target fieldsSetter) {
		target.SetAttributes(fs...)
	})
}

// 设置外部 trance parent
func WithExternalTrace(traceID, parentSpanID string) funcSpanOption {
	if traceID == "" {
		return nil
	}
	return funcSpanOption(func(target *spanImpl) {
		tid, err := ParseHexID(traceID)
		if err != nil {
			target.warnnings = append(target.warnnings, fmt.Errorf("invalid parent trace-id %q: %v", traceID, err))
		} else if tid.Equal(ZeroHexID) {
			target.warnnings = append(target.warnnings, fmt.Errorf("invalid parent trace-id %q: %v", traceID, err))
			return
		} else {
			target.TranceID = tid.String()
		}
		sid, err := ParseHexID(parentSpanID)
		if err != nil {
			target.warnnings = append(target.warnnings, fmt.Errorf("invalid parent span-id %q: %v", traceID, err))
		} else if !sid.Equal(ZeroHexID) {
			target.SpanParentID = sid.Low().String()
		}
	})
}

type funcEndSpanOption func(target *spanImpl)

//	func (f funcEndSpanOption) applySpanOption(target *spanImpl) {
//		f(target)
//	}
func (f funcEndSpanOption) applyEndSpanOption(target *spanImpl) {
	f(target)
}

// 用于动态替换 SpanName
func WithName(getName func() string) nameGetter {
	if getName == nil {
		panic("apm: getName cannot be nil")
	}
	return func(target *spanImpl) {
		target.SetNameGetter(getName)
	}
}

type nameGetter func(target *spanImpl)

func (f nameGetter) applySpanOption(target *spanImpl) {
	f(target)
}
func (f nameGetter) applyEndSpanOption(target *spanImpl) {
	f(target)
}

// 结束时调用
func WithEndCall(fn func(Span)) funcEndSpanOption {
	if fn == nil {
		return nil
	}
	return funcEndSpanOption(func(target *spanImpl) {
		fn(target)
		// target.SetEndCalls([]func(Span){fn})
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
	return defaultEntry.NewSpan(ContextWithCaller(ctx, 3), options...)
}

// // 调整日志堆栈记录深度
//
//	func WithClearup(closer interface{}) *traceOption {
//		return WithEndCall(func(s Span) {
//			Close(closer, s)
//		})
//	}
type ctxkey string

const (
	nameInContextKey ctxkey = "github.com/junhwong/goost$apm.nameInContextKey"
	spanInContextKey ctxkey = "github.com/junhwong/goost$apm.spanInContextKey"
	callerContextKey ctxkey = "github.com/junhwong/goost$apm.callerContextKey"
)

func ContextWithName(ctx context.Context, s string) context.Context {
	if setter, _ := ctx.(interface {
		Set(key string, value interface{})
	}); setter != nil {
		setter.Set(string(nameInContextKey), s)
		return ctx
	}
	return context.WithValue(ctx, nameInContextKey, s)
}
func NameFromContext(ctx context.Context) (s string) {
	s, _ = ctx.Value(nameInContextKey).(string)
	return
}

func SpanContextFrom(ctx context.Context) SpanContext {
	span, _ := ctx.Value(spanInContextKey).(*spanImpl)
	if span != nil {
		return span
	}
	return nil
}

func SpanFrom(ctx context.Context, cotr ...func() (canCreateNew bool, opts []SpanOption)) (context.Context, Span) {
	span, _ := ctx.Value(spanInContextKey).(*spanImpl)
	if span != nil {
		select {
		case <-span.Done():
		default:
			return ctx, span
		}
	}

	if len(cotr) == 0 {
		cotr = append(cotr, func() (canCreateNew bool, opts []SpanOption) { return true, nil })
	}
	fn := cotr[len(cotr)-1]
	if fn == nil {
		return ctx, nil
	}
	canCreateNew, opts := fn()
	if !canCreateNew {
		return ctx, nil
	}
	return defaultEntry.NewSpan(ContextWithCaller(ctx, 3), append([]SpanOption{
		funcSpanOption(func(target *spanImpl) {
			if span == nil {
				return
			}
			target.TranceID = span.GetTranceID()
			target.SpanParentID = span.GetSpanID()
			target.first = false
		}),
	}, opts...)...)
}

var NoCreateNewSpan func() (canCreateNew bool, opts []SpanOption) = nil
