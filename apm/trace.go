package apm

import (
	"context"
	"time"

	"github.com/junhwong/goost/pkg/field"
	"github.com/junhwong/goost/pkg/field/common"
)

type SpanContext struct {
	TranceID     string
	SpanID       string
	SpanParentID string
	Name         string
}

var SpanName = field.String("trace.span.name")

func (span *Span) newEntry() *Entry {
	entry := span.logger.newEntry()
	for _, it := range span.attrs {
		entry.Data.Set(it)
	}
	entry.Data.Set(common.TraceID(span.TranceID))
	return entry
}

type FinishOption interface {
	apply(*finishOption)
}
type finishOption struct {
	getName  func() string
	delegate func(target *finishOption)
}

func (opt *finishOption) apply(target *finishOption) {
	opt.delegate(target)
}
func WithReplaceSpanName(getName func() string) FinishOption {
	if getName == nil {
		panic("getName can't be nil")
	}
	return &finishOption{delegate: func(target *finishOption) {
		target.getName = getName
	}}
}
func (span *Span) Finish(options ...FinishOption) {
	ops := &finishOption{}
	for _, option := range options {
		if option != nil {
			option.apply(ops)
		}
	}
	entry := span.newEntry()
	entry.Time = span.startTime
	entry.Data.Set(common.SpanID(span.SpanID))
	entry.Data.Set(common.SpanParentID(span.SpanParentID))
	entry.Data.Set(common.Duration(time.Since(entry.Time)))
	name := span.Name
	if ops.getName != nil {
		name = ops.getName()
	}
	if name == "" {
		name = span.Name // TODO: 处理未定义名称的情况
	}
	entry.Data.Set(SpanName(name))
	entry.Trace()

	span.logger = nil // 移除关联
}
func (span *Span) Fail() {
	span.failed = true
}

type Span struct {
	ctx context.Context
	SpanContext
	logger    *Logger
	attrs     Fields
	failed    bool
	startTime time.Time
}

func (span *Span) Context() context.Context  { return span.ctx }
func (span *SpanContext) GetTraceId() string { return span.TranceID }
func (span *Span) Debug(a ...interface{})    { span.newEntry().Debug(a...) }
func (span *Span) Info(a ...interface{})     { span.newEntry().Info(a...) }
func (span *Span) Warn(a ...interface{})     { span.newEntry().Warn(a...) }
func (span *Span) Error(a ...interface{})    { span.newEntry().Error(a...) }
func (span *Span) Fatal(a ...interface{})    { span.newEntry().Fatal(a...) }
func (span *Span) Trace(a ...interface{})    { span.newEntry().Trace(a...) }

const (
	spanInContextKey = "$apm.spanInContextKey"
)

func (logger *Logger) WithSpan(ctx context.Context, options ...Option) *Span {
	if ctx == nil {
		ctx = context.TODO()
	}
	option := &traceOption{
		attrs: Fields{},
	}
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(option)
	}
	span := &Span{
		logger:    logger,
		startTime: time.Now(),
		attrs:     option.attrs,
		SpanContext: SpanContext{
			SpanID: newSpanId(),
			Name:   option.name,
		},
	}

	if prent, ok := ctx.Value(spanInContextKey).(*Span); ok && prent != nil {
		span.TranceID = prent.TranceID
		span.SpanParentID = prent.SpanID
		span.logger = prent.logger
	} else {
		span.TranceID = TraceIDFromContext(ctx)
	}

	if span.logger == nil {
		span.logger = logger
	}

	//适配 gin.Context
	if setter, ok := ctx.(interface {
		Set(key string, value interface{})
	}); ok {
		setter.Set(spanInContextKey, span)
		span.ctx = ctx
	} else {
		span.ctx = context.WithValue(ctx, spanInContextKey, span) // nolint
	}

	return span
}

type SpanInterface interface {
	LoggerInterface
	// Finish 结束该Span。
	Finish(options ...FinishOption)
	// Fail 标记该Span为失败。
	Fail()
	Context() context.Context
}

type traceOption struct {
	name     string
	attrs    Fields
	delegate func(*traceOption)
}

func (opt *traceOption) apply(target *traceOption) {
	opt.delegate(target)
}

type Option interface {
	apply(*traceOption)
}

func WithName(name string) Option {
	return &traceOption{delegate: func(target *traceOption) {
		target.name = name
	}}
}

func WithFields(fs ...*Field) Option {
	return &traceOption{delegate: func(target *traceOption) {
		for _, f := range fs {
			target.attrs.Set(f)
		}
	}}
}

func FromContext(ctx context.Context) LoggerInterface {
	span, _ := ctx.Value(spanInContextKey).(*Span)
	if span != nil && span.logger != nil {
		return span
	}
	panic("apm.FromContext: span not found in context")
}

func Start(ctx context.Context, options ...Option) SpanInterface {
	return std.WithSpan(ctx, options...)
}
