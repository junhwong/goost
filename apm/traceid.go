package apm

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type TraceID2 struct {
	High int64
	Low  int64
}

func (t TraceID2) String() string {
	if t.High == 0 {
		return fmt.Sprintf("%016x", t.Low)
	}
	return fmt.Sprintf("%016x%016x", t.High, t.Low)
}

var seededIDGen = rand.New(rand.NewSource(time.Now().UnixNano()))

func NewTraceID() TraceID2 {
	return TraceID2{
		High: int64(time.Now().Unix()<<32) + int64(seededIDGen.Int31()),
		Low:  int64(seededIDGen.Int63()),
	}
}

// randomTimestamped can generate 128 bit time sortable traceid's compatible
// with AWS X-Ray and 64 bit spanid's.
func NewTraceId() string {
	id := TraceID2{
		High: int64(time.Now().Unix()<<32) + int64(seededIDGen.Int31()),
		Low:  int64(seededIDGen.Int63()),
	}
	return id.String()
}

func NewSpanId() string {
	id := TraceID2{
		High: 0,
		Low:  int64(seededIDGen.Int63()),
	}
	return id.String()
}

const (
	spanInContextKey = "$apm.spanInContextKey"
)

func GetTraceID(ctx context.Context) (traceID, spanID string) {
	if ctx == nil {
		return "", ""
	}
	if prent, ok := ctx.Value(spanInContextKey).(SpanContext); ok && prent != nil {
		return prent.GetTranceID(), prent.GetSpanID()
	}
	// https://opentelemetry.io/docs/reference/specification/sdk-environment-variables/
	// https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_conn_man/headers#id21
	if s, ok := ctx.Value("trace_id").(string); ok && s != "" {
		return s, ""
	}
	// https://www.w3.org/TR/trace-context/
	if s, ok := ctx.Value("traceparent").(string); ok && s != "" {
		// version
		// trace-id
		// parent-id
		// trace-flags
		return s, ""
	}
	if s, ok := ctx.Value("request_id").(string); ok && s != "" {
		return s, ""
	}
	return "", ""
}
