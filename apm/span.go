package apm

import (
	"context"
	"strings"
	"time"

	"github.com/junhwong/goost/apm/field"
)

type SpanFactory interface {
	NewSpan(ctx context.Context, options ...SpanOption) (context.Context, Span)
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

type SpanContext interface {
	IsFirst() bool
	GetTranceID() string
	GetSpanID() string
	GetSpanParentID() string
}

type SpanStatus string

const (
	SpanStatusUnset SpanStatus = ""
	SpanStatusError SpanStatus = "Error"
	SpanStatusOk    SpanStatus = "Ok"
)

type spanContext struct {
	TranceID     string
	SpanID       string
	SpanParentID string

	name  string
	first bool
}

func (ctx *spanContext) IsFirst() bool           { return ctx.first }
func (ctx *spanContext) GetTranceID() string     { return ctx.TranceID }
func (ctx *spanContext) GetSpanID() string       { return ctx.SpanID }
func (ctx *spanContext) GetSpanParentID() string { return ctx.SpanParentID }

var _ Span = (*spanImpl)(nil)

type spanImpl struct {
	*logImpl
	spanContext
	failed    bool
	startTime time.Time
	// option    traceOption

	trimFieldPrefix []string
	name            string
	attrs           Fields
	// delegate        func(*traceOption)
	getName  func() string
	endCalls []func(Span)
}

func (log *logImpl) NewSpan(ctx context.Context, options ...SpanOption) (context.Context, Span) {

	if ctx == nil {
		ctx = context.Background()
	}

	span := &spanImpl{
		// option:    option,
		logImpl: log.clone(),
		attrs:   make(Fields, 0),

		startTime: time.Now(),
		spanContext: spanContext{
			SpanID: NewSpanId(),
		},
	}
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.Apply(span)
	}

	if prent, ok := ctx.Value(spanInContextKey).(*spanImpl); ok && prent != nil {
		span.TranceID = prent.TranceID
		span.SpanParentID = prent.SpanID
	} else {
		span.first = true
		span.TranceID, _ = GetTraceID(ctx) // TODO 上级ID
		if len(span.TranceID) == 0 {
			span.TranceID = NewTraceId()
		}
	}

	// 适配 gin.Context
	if setter, ok := ctx.(interface {
		Set(key string, value interface{})
	}); ok {
		setter.Set(spanInContextKey, span)
	} else {
		ctx = context.WithValue(ctx, spanInContextKey, span)
	}
	span.ctx = ctx
	// span.calldepth = option.calldepth // entrylog

	return ctx, span
}

func (span *spanImpl) End(options ...EndSpanOption) {
	for _, option := range options {
		if option != nil {
			option.Apply(span)
		}
	}
	name := span.name
	if span.getName != nil {
		name = span.getName()
	}
	if len(name) == 0 {
		name = span.name
		if len(name) == 0 {
			s := Caller(span.calldepth).Method
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

	var fs Fields
	fs = append(fs, span.attrs...)
	fs = append(fs, SpanName(name))
	fs = append(fs, Time(span.startTime))
	fs = append(fs, SpanID(span.SpanID))
	fs = append(fs, SpanParentID(span.SpanParentID))
	fs = append(fs, Duration(time.Since(span.startTime))) // Latency
	fs = append(fs, TraceID(span.TranceID))
	for _, fn := range span.endCalls {
		fn(span)
	}
	if span.failed {
		fs = append(fs, TraceError(span.failed))
	}
	fs = append(fs, LevelField(int(LevelTrace)))
	span.LogFS(nil, fs...) //span.ctx , Trace, TODO: calldepth 不能获取到 defer 位置
	span.dispatcher = nil  // 移除关联,
}

func (span *spanImpl) Context() SpanContext { return span }

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

func (span *spanImpl) SetStatus(code SpanStatus, description string, failure ...bool) {
	span.SetAttributes(SpanStatusCode(string(code)), SpanStatusDescription(description))
	for _, v := range failure {
		if v {
			span.failed = true
		}
	}
}
func (s *spanImpl) SetAttributes(a ...field.Field) { s.attrs = append(s.attrs, a...) }
func (s *spanImpl) SetNameGetter(a func() string)  { s.getName = a }
func (s *spanImpl) SetEndCalls(a []func(Span))     { s.endCalls = a }
