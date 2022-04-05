package apm

import (
	"context"
	"time"

	"github.com/junhwong/goost/apm/level"
	"github.com/junhwong/goost/pkg/field"
)

type SpanInterface interface {
	NewSpan(ctx context.Context, options ...Option) (context.Context, Span)
}
type Span interface {
	Logger
	// Trace(...interface{})     //
	End(options ...EndOption) // 结束该Span。
	Fail()                    // Fail 标记该Span为失败。
	FailIf(err error) bool    // 如果`err`不为`nil`, 则标记失败并返回`true`，否则`false`
	PanicIf(err error)        // 如果`err`不为`nil`, 则标记失败并`panic`
	Context() context.Context // 返回与该span关联的上下文
}
type SpanContext struct {
	TranceID     string
	SpanID       string
	SpanParentID string
	Name         string
	first        bool
}

func (ctx *SpanContext) IsFirst() bool {
	return ctx.first
}

const (
	spanInContextKey = "$apm.spanInContextKey"
)

var _ Span = (*spanImpl)(nil)

type spanImpl struct {
	entryLog
	SpanContext
	failed    bool
	startTime time.Time
	option    traceOption
}

func newSpan(ctx context.Context, logger *DefaultLogger, options []Option) (context.Context, *spanImpl) {
	if logger == nil {
		panic("apm: logger cannot be nil")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	option := traceOption{
		attrs: make([]field.Field, 0),
	}
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(&option)
	}
	span := &spanImpl{
		option:    option,
		startTime: time.Now(),
		SpanContext: SpanContext{
			SpanID: newSpanId(),
			Name:   option.name,
		},
	}

	if prent, ok := ctx.Value(spanInContextKey).(*spanImpl); ok && prent != nil {
		span.TranceID = prent.TranceID
		span.SpanParentID = prent.SpanID
	} else {
		span.first = true
		span.TranceID = getTraceID(ctx)
		if span.TranceID == "" {
			span.TranceID = newTraceId()
		}
	}

	// 适配 gin.Context
	if setter, ok := ctx.(interface {
		Set(key string, value interface{})
	}); ok {
		setter.Set(spanInContextKey, span)
	} else {
		ctx = context.WithValue(ctx, spanInContextKey, span) // nolint
	}
	span.option = option
	span.logger = logger
	span.ctx = ctx
	return ctx, span
}

func (span *spanImpl) End(options ...EndOption) {
	for _, option := range options {
		if option != nil {
			option.applyEnd(&span.option)
		}
	}
	name := span.Name
	if span.option.getName != nil {
		name = span.option.getName()
	}
	if name == "" {
		name = span.Name // TODO: 处理未定义名称的情况
	}

	fs := []interface{}{}
	for _, it := range span.option.attrs {
		fs = append(fs, it)
	}
	fs = append(fs, _entryTime(span.startTime))
	fs = append(fs, _entrySpanID(span.SpanID))
	fs = append(fs, _entrySpanParentID(span.SpanParentID))
	fs = append(fs, _entryDuration(time.Since(span.startTime))) // Latency
	fs = append(fs, _entrySpanName(name))
	if span.failed {
		fs = append(fs, _entryTraceError(span.failed))
	}
	span.logger.Log(span.ctx, 1, level.Trace, fs)
	span.logger = nil // 移除关联,
}

func (span *spanImpl) Context() context.Context { return span.ctx }

// func (span *SpanContext) GetTraceId() string         { return span.TranceID }

// 标记失败
func (span *spanImpl) Fail() {
	span.failed = true
}

func (span *spanImpl) FailIf(err error) bool {
	if err == nil {
		return false
	}
	span.failed = true
	return true
}

func (span *spanImpl) PanicIf(err error) {
	if err == nil {
		return
	}
	span.failed = true
	panic(err) // TODO 错误包装
}
