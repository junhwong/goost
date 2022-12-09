package deflog

// // type SpanOption interface {
// // 	applySpanOption(*traceOption)
// // }
// // type StartSpanOption interface {
// // 	applyStartSpanOption(*traceOption)
// // }
// // type EndSpanOption interface {
// // 	applyEndOption(*traceOption)
// // }

// //	type SpanOptionSetter interface {
// //		SetNameGetter(a func() string)
// //		SetAttributes(a []field.Field)
// //		SetCalldepth(a int)
// //	}
// //
// //	type EndSpanOptionSetter interface {
// //		SetNameGetter(a func() string)
// //		SetAttributes(a []field.Field)
// //		SetEndCalls(a []func(apm.Span))
// //	}
// type traceOption struct {
// 	trimFieldPrefix []string
// 	name            string
// 	attrs           []field.Field
// 	delegate        func(*traceOption)
// 	getName         func() string
// 	calldepth       int
// 	endCalls        []func(apm.Span)
// }

// func (opt *traceOption) SetNameGetter(a func() string)  { opt.getName = a }
// func (opt *traceOption) SetAttributes(a []field.Field)  { opt.attrs = a }
// func (opt *traceOption) SetCalldepth(a int)             { opt.calldepth = a }
// func (opt *traceOption) SetEndCalls(a []func(apm.Span)) { opt.endCalls = a }

// // func (opt *traceOption) applySpanOption(target *traceOption) {
// // 	opt.delegate(target)
// // }
// // func (opt *traceOption) applyEndOption(target *traceOption) {
// // 	opt.applySpanOption(target)
// // }

// // func (opt *traceOption) applyStartSpanOption(target *traceOption) {
// // 	opt.applyStartSpanOption(target)
// // }

// // func WithName(name string) SpanOption {
// // 	return &traceOption{delegate: func(target *traceOption) {
// // 		target.name = name
// // 	}}
// // }

// // // 调整日志堆栈记录深度
// // func WithCallDepth(depth int) *traceOption {
// // 	return &traceOption{delegate: func(target *traceOption) {
// // 		target.calldepth = depth
// // 	}}
// // }

// // // 替换SpanName
// // func WithReplaceSpanName(getName func() string) EndSpanOption {
// // 	if getName == nil {
// // 		panic("apm: getName cannot be nil")
// // 	}
// // 	return &traceOption{delegate: func(target *traceOption) {
// // 		target.getName = getName
// // 	}}
// // }

// // type appendFields func() []field.Field

// // //	func (fn appendFields) applyInterface(impl *stdImpl) {
// // //		fs := fn()
// // //		impl.fields = append(impl.fields, fs...)
// // //	}
// // func (fn appendFields) applySpanOption(opt *traceOption) {
// // 	opt.attrs = append(opt.attrs, fn()...)
// // }
// // func (fn appendFields) applyEndOption(opt *traceOption) {
// // 	opt.attrs = append(opt.attrs, fn()...)
// // }

// // // Deprecated 已经废弃
// // func WithTrimFieldPrefix(prefix ...string) SpanOption {
// // 	return &traceOption{delegate: func(target *traceOption) {
// // 		target.trimFieldPrefix = prefix
// // 	}}
// // }

// // // func Start(ctx context.Context, options ...SpanOption) (context.Context, Span) {
// // // 	return std.NewSpan(ctx, 0, options...)
// // // }

// // // 调整日志堆栈记录深度
// // func WithEndCall(fn func(apm.Span)) *traceOption {
// // 	return &traceOption{delegate: func(target *traceOption) {
// // 		// if target.endCalls == nil {
// // 		// 	target.endCalls = []func(){fn}
// // 		// 	return
// // 		// }
// // 		target.endCalls = append(target.endCalls, fn)
// // 	}}
// // }

// // // 调整日志堆栈记录深度
// // func WithClearup(closer interface{}) *traceOption {
// // 	return WithEndCall(func(s apm.Span) {
// // 		// Close(closer, s)
// // 	})
// // }
