package logs

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/junhwong/goost/pkg/field"
	"github.com/junhwong/goost/pkg/field/common"
)

type Tracer struct {
}

func Trace(point, traceId string) {

}

type ISpanContext interface {
	ILogger
	GetTraceId() string
}

type SpanContext struct {
	TranceID     string
	SpanID       string
	SpanParentID string
}

func (span *SpanContext) GetTraceId() string { return span.TranceID }

func (span *Span) newEntry() *Entry {
	entry := span.logger.newEntry()
	traceId := common.Message("bbq")
	entry.Data[traceId.Key] = traceId
	return entry
}
func (span *Span) Finish() {

}

type Span struct {
	ctx    context.Context
	logger *Logger
}

func (span *Span) Debug(a ...interface{}) { span.newEntry().Debug(a...) }
func (span *Span) Info(a ...interface{})  { span.newEntry().Info(a...) }
func (span *Span) Warn(a ...interface{})  { span.newEntry().Warn(a...) }
func (span *Span) Error(a ...interface{}) { span.newEntry().Error(a...) }
func (span *Span) Fatal(a ...interface{}) { span.newEntry().Fatal(a...) }
func (span *Span) Trace(a ...interface{}) { span.newEntry().Trace(a...) }

var (
	_TraceID      = field.String("trace-id")
	_SpanID       = field.String("span-id")
	_SpanParentID = field.String("span-parentid")
)

type ctxKey struct{}

var (
	spanInContextKey = ctxKey{}
)

func (logger *Logger) WithSpan(ctx context.Context, spanName string) *Span {
	if ctx == nil {
		ctx = context.TODO()
	}
	var sctx SpanContext
	if spanCtx, ok := ctx.Value(spanInContextKey).(SpanContext); !ok {
		ctx = context.WithValue(ctx, spanInContextKey, spanCtx)
	} else {
		ctx = context.WithValue(ctx, spanInContextKey, SpanContext{})
		sctx = SpanContext{}
	}
	sctx.SpanID = STraceID{}.String()
	return &Span{ctx, logger}
}

type STraceID struct {
	High uint64
	Low  uint64
}

func (t STraceID) String() string {
	if t.High == 0 {
		return fmt.Sprintf("%016x", t.Low)
	}
	return fmt.Sprintf("%016x%016x", t.High, t.Low)
}

var seededIDGen = rand.New(rand.NewSource(time.Now().UnixNano()))

// randomTimestamped can generate 128 bit time sortable traceid's compatible
// with AWS X-Ray and 64 bit spanid's.
func newTraceId() string {
	id := STraceID{
		High: uint64(time.Now().Unix()<<32) + uint64(seededIDGen.Int31()),
		Low:  uint64(seededIDGen.Int63()),
	}
	return id.String()
}

func (logger *Logger) SpanFromContext(ctx context.Context) *Span {
	span, _ := ctx.Value("goost.logs.span").(*Span)
	if span == nil {
		span = logger.WithSpan(ctx, "")
	}
	return span
}
