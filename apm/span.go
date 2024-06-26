package apm

import (
	"context"
	"reflect"
	"strings"
	"time"

	"github.com/junhwong/goost/apm/field"
	"github.com/junhwong/goost/apm/field/loglevel"
)

// 用于构造一个新的 Span
type SpanFactory interface {
	NewSpan(ctx context.Context, options ...SpanOption) (context.Context, Span)
}

type Span interface {
	context.Context
	Logger
	End(options ...EndSpanOption)                  // 结束该Span。
	FailIf(err error, description ...string) error // 如果`err`不为`nil`, 则标记失败并返回err。
	PanicIf(err error, description ...string)      // 如果`err`不为`nil`, 则标记失败并`panic`
	SetAttributes(attrs ...*field.Field)           //
	SpanContext() SpanContext                      // 返回与该span关联的上下文
	Name() string                                  // 返回 SpanName. 注意: 只有在 End 后才能最后决定.
	Duration() time.Duration                       // 返回执行时间. 注意: 只有在 End 后才能最后决定.
}

// 与 Span 关联的上下文
type SpanContext interface {
	IsFirst() bool           // 是否是第一个 Span
	GetTranceID() string     // 获取当前 TranceID
	GetSpanID() string       // 获取当前 SpanID
	GetSpanParentID() string // 获取当前 SpanParentID
}

type SpanStatus string

const (
	SpanStatusUnset SpanStatus = ""
	SpanStatusError SpanStatus = "Error"
	SpanStatusOk    SpanStatus = "Ok"
)

var (
	_ Span        = (*spanImpl)(nil)
	_ SpanContext = (*spanImpl)(nil)
)

type spanImpl struct {
	context.Context
	*FieldsEntry
	failed       bool
	failedDesc   string
	TranceID     string
	SpanID       string
	SpanParentID string
	name         string
	duration     time.Duration
	first        bool
	getName      func() string
	endCalls     []func(Span)
	cancel       context.CancelFunc
	warnnings    []error
}

func (e *FieldsEntry) NewSpan(ctx context.Context, options ...SpanOption) (context.Context, Span) {
	if ctx == nil {
		ctx = context.Background()
	}

	span := &spanImpl{
		FieldsEntry: e.new(),
		SpanID:      NewHexID().Low().String(),
	}
	// span.calldepth++

	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.applySpanOption(span)
	}

	span.CallerInfo = CallerFrom(ctx)
	if span.CallerInfo == nil {
		span.CallerInfo = &CallerInfo{}
		doCaller(span.calldepth, span.CallerInfo)
	}

	if len(span.TranceID) != 0 {
	} else if prent := SpanContextFrom(ctx); prent != nil {
		span.TranceID = prent.GetTranceID()
		span.SpanParentID = prent.GetSpanID()
	} else {
		span.first = true
		span.TranceID = NewHexID().String()
		span.SpanParentID = make(HexID, 16).Low().String()
	}

	// if setter, ok := ctx.(interface {
	// 	Set(key string, value interface{})
	// }); ok {
	// 	setter.Set(string(spanInContextKey), span)
	// } else {
	// 	ctx = context.WithValue(ctx, spanInContextKey, span)
	// }

	span.Context, span.cancel = context.WithCancel(ctx)

	// 适配 gin.Context 这类可变 Context, 以贯穿其生命周期
	if setter, ok := ctx.(interface {
		Set(key string, value interface{})
	}); ok {
		setter.Set(string(spanInContextKey), span)
	} else if setter, ok := ctx.(interface {
		SetAttribute(key string, value interface{})
	}); ok {
		setter.SetAttribute(string(spanInContextKey), span)
	}

	// 打印跟踪
	if len(span.warnnings) != 0 {
		for _, it := range span.warnnings {
			span.Warn(it)
		}
	}
	span.warnnings = nil

	return ctx, span
}

// 从写 Context.Value 方法
func (span *spanImpl) Value(key any) any {
	if reflect.DeepEqual(key, spanInContextKey) {
		return span
	}
	return span.Context.Value(key)
}

func (span *spanImpl) End(options ...EndSpanOption) {
	span.mu.Lock()
	if span.FieldsEntry == nil {
		span.mu.Unlock()
		return
	}
	defer span.mu.Unlock()
	defer span.cancel()

	span.duration = time.Since(span.GetTime())
	for _, fn := range span.endCalls {
		if fn != nil {
			fn(span)
		}
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
		span.Set(SpanStatusCode(string(SpanStatusError)))

		if len(span.failedDesc) > 0 {
			span.Set(SpanStatusDescription(string(span.failedDesc)))
		}
	}

	span.do([]any{span.ctx}, func() {
		span.Set(LevelField(loglevel.Trace2))
		span.Set(SpanName(name))
		span.Set(SpanID(span.SpanID))
		if len(span.SpanParentID) > 0 {
			span.Set(SpanParentID(span.SpanParentID))
		}

		span.Set(Duration(span.duration)) // Latency
		span.Set(TraceIDField(span.TranceID))
	})

	span.FieldsEntry = nil
}

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

func (s *spanImpl) SetNameGetter(a func() string) { s.getName = a }
func (s *spanImpl) SetEndCalls(a []func(Span))    { s.endCalls = a }
func (s *spanImpl) SpanContext() SpanContext      { return s }
func (s *spanImpl) Name() string                  { return s.name }
func (s *spanImpl) Duration() time.Duration       { return s.duration }
func (s *spanImpl) IsFirst() bool                 { return s.first }
func (s *spanImpl) GetTranceID() string           { return s.TranceID }
func (s *spanImpl) GetSpanID() string             { return s.SpanID }
func (s *spanImpl) GetSpanParentID() string       { return s.SpanParentID }
