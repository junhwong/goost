package apm

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/junhwong/goost/apm/field"
)

type SpanFactory interface {
	NewSpan(ctx context.Context, options ...SpanOption) (context.Context, Span)
}

type Span interface {
	Logger
	End(options ...EndSpanOption)                                   // 结束该Span。
	Fail(error) error                                               // Fail 标记该Span为失败。
	FailIf(err error) bool                                          // 如果`err`不为`nil`, 则标记失败并返回`true`，否则`false`
	PanicIf(err error)                                              // 如果`err`不为`nil`, 则标记失败并`panic`
	SetStatus(code SpanStatus, description string, failure ...bool) // 设置状态
	SetAttributes(attrs ...*Field)                                  //
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
	*FieldsEntry
	mu sync.Mutex
	*logImpl
	spanContext
	failed   bool
	name     string
	getName  func() string
	endCalls []func(Span)

	// trimFieldPrefix []string
	// attrs           Fields
	// delegate        func(*traceOption)
}

func (log *logImpl) NewSpan(ctx context.Context, options ...SpanOption) (context.Context, Span) {
	if ctx == nil {
		ctx = context.Background()
	}
	log = log.clone()
	id := NewHexID()
	id.High = 0
	span := &spanImpl{
		FieldsEntry: &FieldsEntry{
			Level:  field.LevelTrace,
			Time:   time.Now(),
			Labels: log.fields,
		},
		logImpl: log,
		spanContext: spanContext{
			SpanID: id.String(),
		},
	}
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.applySpanOption(span)
	}
	var ok bool
	span.CallerInfo, ok = CallerFromContext(ctx)
	if !ok {
		doCaller(span.calldepth, &span.CallerInfo)
		ctx = context.WithValue(ctx, callerContextKey, span.CallerInfo)
	}
	if prent, ok := ctx.Value(spanInContextKey).(*spanImpl); ok && prent != nil {
		span.TranceID = prent.TranceID
		span.SpanParentID = prent.SpanID
	} else {
		span.first = true
		span.TranceID, span.SpanParentID = GetTraceID(ctx)
		if len(span.TranceID) == 0 {
			span.TranceID = NewHexID().String()
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
	if span.logImpl == nil {
		return
	}

	span.mu.Lock()
	defer span.mu.Unlock()
	if span.logImpl == nil {
		return
	}
	for _, option := range options {
		if option != nil {
			option.applyEndSpanOption(span)
		}
	}
	name := span.name
	if span.getName != nil {
		name = span.getName()
	}
	if len(name) == 0 {
		s := span.CallerInfo.Method
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
		if i = strings.LastIndex(name, "/"); i > 0 {
			name = name[i+1:]
		}
	}

	span.Labels = append(span.Labels, SpanName(name))
	span.Labels = append(span.Labels, SpanID(span.SpanID))
	if len(span.SpanParentID) > 0 {
		span.Labels = append(span.Labels, SpanParentID(span.SpanParentID))
	}
	span.Labels = append(span.Labels, Duration(time.Since(span.Time))) // Latency
	span.Labels = append(span.Labels, TraceIDField(span.TranceID))
	for _, fn := range span.endCalls {
		fn(span)
	}
	if span.failed {
		span.Labels = append(span.Labels, SpanStatusCode(string(SpanStatusError)))
	}

	span.LogFS(span.FieldsEntry, []any{span.ctx}) //span.ctx , Trace, TODO: calldepth 不能获取到 defer 位置
	span.logImpl = nil
	span.FieldsEntry = nil
}

func (span *spanImpl) Context() SpanContext { return span }

// 标记失败
func (span *spanImpl) Fail(err error) error {
	if err == nil {
		return nil
	}
	span.failed = true
	return err
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
func (s *spanImpl) SetAttributes(a ...*Field)     { s.Labels = append(s.Labels, a...) }
func (s *spanImpl) SetNameGetter(a func() string) { s.getName = a }
func (s *spanImpl) SetEndCalls(a []func(Span))    { s.endCalls = a }
