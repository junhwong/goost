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
	End(options ...EndSpanOption)                  // 结束该Span。
	FailIf(err error, description ...string) error // 如果`err`不为`nil`, 则标记失败并返回err。
	PanicIf(err error, description ...string)      // 如果`err`不为`nil`, 则标记失败并`panic`
	SetAttributes(attrs ...*field.Field)           //
	SpanContext() SpanContext                      // 返回与该span关联的上下文
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
	spanContext
	failed     bool
	failedDesc string
	name       string
	getName    func() string
	endCalls   []func(Span)
}

func (log *FieldsEntry) NewSpan(ctx context.Context, options ...SpanOption) (context.Context, Span) {
	if ctx == nil {
		ctx = context.Background()
	}

	span := &spanImpl{
		FieldsEntry: log.clone(),
		spanContext: spanContext{
			SpanID: NewHexID().Low().String(),
		},
	}
	span.Level = field.LevelTrace
	span.Time = time.Now()

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

	// 适配 gin.Context 这类可变 Context
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
	if span.FieldsEntry == nil {
		return
	}

	span.mu.Lock()
	defer span.mu.Unlock()

	if span.FieldsEntry == nil {
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

	if span.failed {
		span.Fields.Set(SpanStatusCode(string(SpanStatusError)))

		if len(span.failedDesc) > 0 {
			span.Fields.Set(SpanStatusDescription(string(span.failedDesc)))
		}
	}
	for _, fn := range span.endCalls {
		fn(span)
	}

	span.do([]any{span.ctx}, func() {
		span.Fields.Set(SpanName(name))
		span.Fields.Set(SpanID(span.SpanID))
		if len(span.SpanParentID) > 0 {
			span.Fields.Set(SpanParentID(span.SpanParentID))
		}
		span.Fields.Set(Duration(time.Since(span.Time))) // Latency
		span.Fields.Set(TraceIDField(span.TranceID))
	})

	span.FieldsEntry = nil
}

func (span *spanImpl) SpanContext() SpanContext { return span }

// 标记失败
func (span *spanImpl) FailIf(err error, description ...string) error {
	if err == nil {
		return nil
	}
	span.failed = true
	if len(description) > 0 {
		span.failedDesc = description[len(description)-1]
	} else {
		span.failedDesc = err.Error()
	}
	return err
}

func (span *spanImpl) PanicIf(err error, description ...string) {
	if err := span.FailIf(err, description...); err != nil {
		panic(err) // TODO 错误包装
	}
}

//	func (span *spanImpl) SetStatus(code SpanStatus, description string, failure ...bool) {
//		span.SetAttributes(SpanStatusCode(string(code)), SpanStatusDescription(description))
//		for _, v := range failure {
//			if v {
//				span.failed = true
//			}
//		}
//	}
func (s *spanImpl) SetAttributes(a ...*field.Field) { s.Fields = append(s.Fields, a...) }
func (s *spanImpl) SetNameGetter(a func() string)   { s.getName = a }
func (s *spanImpl) SetEndCalls(a []func(Span))      { s.endCalls = a }
