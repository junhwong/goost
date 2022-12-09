package deflog

// const (
// 	spanInContextKey = "$apm.spanInContextKey"
// )

// var _ apm.Span = (*spanImpl)(nil)

// type spanImpl struct {
// 	entryLog
// 	SpanContext
// 	failed    bool
// 	startTime time.Time
// 	option    traceOption
// }

// type SpanContext struct {
// 	TranceID     string
// 	SpanID       string
// 	SpanParentID string

// 	name  string
// 	first bool
// }

// func (ctx *SpanContext) IsFirst() bool {
// 	return ctx.first
// }
// func (ctx *SpanContext) GetTranceID() string     { return ctx.TranceID }
// func (ctx *SpanContext) GetSpanID() string       { return ctx.SpanID }
// func (ctx *SpanContext) GetSpanParentID() string { return ctx.SpanParentID }

// func newSpan(ctx context.Context, logger *DefaultLogger, calldepth int, options []apm.SpanOption) (context.Context, *spanImpl) {
// 	if logger == nil {
// 		panic("apm: logger cannot be nil")
// 	}
// 	if ctx == nil {
// 		ctx = context.Background()
// 	}
// 	option := traceOption{
// 		calldepth: calldepth,
// 		attrs:     make([]field.Field, 0),
// 	}
// 	for _, opt := range options {
// 		if opt == nil {
// 			continue
// 		}
// 		opt.Apply(&option)
// 	}

// 	span := &spanImpl{
// 		option:    option,
// 		startTime: time.Now(),
// 		SpanContext: SpanContext{
// 			SpanID: apm.NewSpanId(),
// 			name:   option.name,
// 		},
// 	}

// 	if prent, ok := ctx.Value(spanInContextKey).(*spanImpl); ok && prent != nil {
// 		span.TranceID = prent.TranceID
// 		span.SpanParentID = prent.SpanID
// 	} else {
// 		span.first = true
// 		span.TranceID, _ = apm.GetTraceID(ctx) // TODO 上级ID
// 		if span.TranceID == "" {
// 			span.TranceID = apm.NewTraceId()
// 		}
// 	}

// 	// 适配 gin.Context
// 	if setter, ok := ctx.(interface {
// 		Set(key string, value interface{})
// 	}); ok {
// 		setter.Set(spanInContextKey, span)
// 	} else {
// 		ctx = context.WithValue(ctx, spanInContextKey, span) // nolint
// 	}
// 	span.option = option
// 	span.logger = logger
// 	span.ctx = ctx
// 	span.calldepth = option.calldepth // entrylog
// 	return ctx, span
// }

// func (span *spanImpl) End(options ...apm.EndSpanOption) {
// 	for _, option := range options {
// 		if option != nil {
// 			option.Apply(&span.option)
// 		}
// 	}
// 	name := span.name
// 	if span.option.getName != nil {
// 		name = span.option.getName()
// 	}
// 	if len(name) == 0 {
// 		name = span.name
// 		if len(name) == 0 {
// 			s := runtime.Caller(span.calldepth).Method
// 			i := strings.LastIndex(s, ".")
// 			if i > 0 {
// 				name = strings.Trim(s[i+1:], ".")
// 				s = s[:i]
// 			}
// 			s = strings.Trim(strings.SplitN(s, "(", 2)[0], ".")
// 			if len(name) > 0 {
// 				name = s + "." + name
// 			} else {
// 				name = s
// 			}
// 		}
// 	}

// 	fs := []interface{}{}
// 	for _, it := range span.option.attrs {
// 		fs = append(fs, it)
// 	}
// 	fs = append(fs, apm.SpanName(name))
// 	fs = append(fs, apm.Time(span.startTime))
// 	fs = append(fs, apm.SpanID(span.SpanID))
// 	fs = append(fs, apm.SpanParentID(span.SpanParentID))
// 	fs = append(fs, apm.Duration(time.Since(span.startTime))) // Latency

// 	for _, fn := range span.option.endCalls {
// 		fn(span)
// 	}

// 	if span.failed {
// 		fs = append(fs, apm.TraceError(span.failed))
// 	}
// 	span.calldepth = span.option.calldepth
// 	span.logger.Log(span.ctx, span.calldepth, level.Trace, fs) // TODO: calldepth 不能获取到 defer 位置
// 	span.logger = nil                                          // 移除关联,
// }

// func (span *spanImpl) Context() apm.SpanContext { return span }

// // 标记失败
// func (span *spanImpl) Fail() {
// 	span.failed = true
// }

// func (span *spanImpl) FailIf(err error) bool {
// 	if err == nil {
// 		return false
// 	}
// 	span.failed = true
// 	return true
// }

// func (span *spanImpl) PanicIf(err error) {
// 	if err == nil {
// 		return
// 	}
// 	span.failed = true
// 	panic(err) // TODO 错误包装
// }

// func (span *spanImpl) SetStatus(code apm.SpanStatus, description string, failure ...bool) {
// 	span.SetAttributes(apm.SpanStatusCode(string(code)), apm.SpanStatusDescription(description))
// 	for _, v := range failure {
// 		if v {
// 			span.failed = true
// 		}
// 	}
// }

// func (span *spanImpl) SetAttributes(attrs ...field.Field) {
// 	span.option.attrs = append(span.option.attrs, attrs...)
// }

// //failure
