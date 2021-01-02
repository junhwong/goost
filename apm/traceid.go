package apm

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

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

func newSpanId() string {
	id := STraceID{
		High: 0,
		Low:  uint64(seededIDGen.Int63()),
	}
	return id.String()
}

// TODO
func TraceIDFromContext(ctx context.Context) string {
	if ctx == nil {
		panic("ctx cannot be nil")
	}
	//https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_conn_man/headers#id21
	if s, ok := ctx.Value("trace_id").(string); ok && s != "" {
		return s
	}
	//https://www.w3.org/TR/trace-context/
	if s, ok := ctx.Value("traceparent").(string); ok && s != "" {
		// version
		// trace-id
		// parent-id
		// trace-flags
		return s
	}
	if s, ok := ctx.Value("request_id").(string); ok && s != "" {
		return s
	}
	return newTraceId()
}
