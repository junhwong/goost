package apm

import (
	"context"
	"strings"
	"time"

	"github.com/junhwong/goost/apm/level"
	"github.com/junhwong/goost/pkg/field"
	"github.com/junhwong/goost/runtime"
)

type SpanInterface interface {
	NewSpan(ctx context.Context, calldepth int, options ...SpanOption) (context.Context, Span)
}
type Span interface {
	Logger
	End(options ...EndSpanOption)                                   // 结束该Span。
	Fail()                                                          // Fail 标记该Span为失败。
	FailIf(err error) bool                                          // 如果`err`不为`nil`, 则标记失败并返回`true`，否则`false`
	PanicIf(err error)                                              // 如果`err`不为`nil`, 则标记失败并`panic`
	SetStatus(code SpanStatus, description string, failure ...bool) // 设置状态
	SetAttributes(attrs ...field.Field)                             //
	Context() SpanContext                                           // 返回与该span关联的上下文
}
type SpanContext struct {
	TranceID     string
	SpanID       string
	SpanParentID string

	name  string
	first bool
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

func newSpan(ctx context.Context, logger *DefaultLogger, calldepth int, options []SpanOption) (context.Context, *spanImpl) {
	if logger == nil {
		panic("apm: logger cannot be nil")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	option := traceOption{
		calldepth: calldepth,
		attrs:     make([]field.Field, 0),
	}
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.applySpanOption(&option)
	}

	span := &spanImpl{
		option:    option,
		startTime: time.Now(),
		SpanContext: SpanContext{
			SpanID: newSpanId(),
			name:   option.name,
		},
	}

	if prent, ok := ctx.Value(spanInContextKey).(*spanImpl); ok && prent != nil {
		span.TranceID = prent.TranceID
		span.SpanParentID = prent.SpanID
	} else {
		span.first = true
		span.TranceID, _ = getTraceID(ctx) // TODO 上级ID
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
	span.calldepth = option.calldepth // entrylog
	return ctx, span
}

func (span *spanImpl) End(options ...EndSpanOption) {
	for _, option := range options {
		if option != nil {
			option.applyEndOption(&span.option)
		}
	}
	name := span.name
	if span.option.getName != nil {
		name = span.option.getName()
	}
	if len(name) == 0 {
		name = span.name
		if len(name) == 0 {
			s := runtime.Caller(span.calldepth).Method
			i := strings.LastIndex(s, ".")
			if i > 0 {
				name = strings.Trim(s[i+1:], ".")
				s = s[:i]
			}
			s = strings.Trim(strings.SplitN(s, "(", 2)[0], ".")
			if len(name) > 0 {
				name = s + "." + name
			} else {
				name = s
			}
		}
	}

	fs := []interface{}{}
	for _, it := range span.option.attrs {
		fs = append(fs, it)
	}
	fs = append(fs, SpanName(name))
	fs = append(fs, Time(span.startTime))
	fs = append(fs, SpanID(span.SpanID))
	fs = append(fs, SpanParentID(span.SpanParentID))
	fs = append(fs, Duration(time.Since(span.startTime))) // Latency

	for _, fn := range span.option.endCalls {
		fn(span)
	}

	if span.failed {
		fs = append(fs, TraceError(span.failed))
	}
	span.calldepth = span.option.calldepth
	span.logger.Log(span.ctx, span.calldepth, level.Trace, fs) // TODO: calldepth 不能获取到 defer 位置
	span.logger = nil                                          // 移除关联,
}

func (span *spanImpl) Context() SpanContext { return span.SpanContext }

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

type SpanStatus string

const (
	SpanStatusUnset SpanStatus = "Unset"
	SpanStatusError SpanStatus = "Error"
	SpanStatusOk    SpanStatus = "Ok"
)

func (span *spanImpl) SetStatus(code SpanStatus, description string, failure ...bool) {
	span.SetAttributes(spanStatusCode(string(code)), spanStatusDescription(description))
	for _, v := range failure {
		if v {
			span.failed = true
		}
	}
}

func (span *spanImpl) SetAttributes(attrs ...field.Field) {
	span.option.attrs = append(span.option.attrs, attrs...)
}

//failure
