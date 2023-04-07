package apm

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/junhwong/goost/apm/field"
)

// 符合 W3C 规范的 TraceID 或 SpanID.
// https://www.w3.org/TR/trace-context/#trace-id
type HexID struct {
	High uint64
	Low  uint64
}

func (id HexID) Bytes() []byte {
	var b []byte
	if i := id.High; i > 0 {
		b = binary.BigEndian.AppendUint64(b, i)
	}
	if i := id.Low; i > 0 {
		b = binary.BigEndian.AppendUint64(b, i)
	}
	return b
}

func (id HexID) String() string {
	b := id.Bytes()
	if len(b) == 0 {
		return ""
	}
	return fmt.Sprintf("%x", id.Bytes())
}

var seededIDGen = rand.New(rand.NewSource(time.Now().UnixNano()))
var mu sync.Mutex

// randomTimestamped can generate 128 bit time sortable traceid's compatible
// with AWS X-Ray and 64 bit spanid's.
func NewHexID() HexID {
	mu.Lock()
	id := HexID{
		High: uint64(time.Now().Unix()<<32) + uint64(seededIDGen.Int31()),
		Low:  uint64(seededIDGen.Int63()),
	}
	mu.Unlock()
	return id
}

var (
	errInvalidHexID = errors.New("hex-id can only contain hex characters, len (16 or 32)")
)

// ParseHexID returns a HexID from a hex string.
func ParseHexID(h string) (HexID, error) {
	t := HexID{}
	decoded, err := hex.DecodeString(h)
	if err != nil {
		return t, errInvalidHexID
	}
	switch len(decoded) {
	case 16:
		t.High = binary.BigEndian.Uint64(decoded[:8])
		decoded = decoded[8:]
		fallthrough
	case 8:
		t.Low = binary.BigEndian.Uint64(decoded)
	default:
		return t, errInvalidHexID
	}
	return t, nil
}

const (
	spanInContextKey = "$apm.spanInContextKey"
)

// Deprecated: Drivers
func GetTraceID(ctx context.Context) (traceID, spanID string) {
	if ctx == nil {
		return "", ""
	}
	if p, ok := ctx.Value(spanInContextKey).(SpanContext); ok && p != nil {
		return p.GetTranceID(), p.GetSpanID()
	}
	// https://opentelemetry.io/docs/reference/specification/sdk-environment-variables/
	// https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_conn_man/headers#id21
	if s, ok := ctx.Value("trace_id").(string); ok && s != "" {
		return s, ""
	}
	// todo https://www.w3.org/TR/trace-context/
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

// 解析 W3C trace.
//
// 示例: `00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01`.
//
// see: https://www.w3.org/TR/trace-context/#traceparent-header
func ParseW3Traceparent(traceparent string) (version byte, traceID, parentSpanID HexID, flags byte, err error) {
	arr := strings.Split(traceparent, "-")
	if len(arr) != 4 {
		return
	}
	decoded, ex := hex.DecodeString(arr[0])
	if ex != nil || len(decoded) != 1 {
		err = fmt.Errorf("invalid version")
		return
	}
	version = decoded[0]
	decoded, ex = hex.DecodeString(arr[3])
	if ex != nil || len(decoded) != 1 {
		err = fmt.Errorf("invalid flags")
		return
	}
	flags = decoded[0]

	traceID, err = ParseHexID(arr[1])
	if err != nil {
		return
	}

	parentSpanID, err = ParseHexID(arr[2])
	if err != nil {
		return
	}

	return
}

// 解析 W3C tracestate.
//
// 示例: `rojo=00f067aa0ba902b7,congo=t61rcWkgMzE`.
//
// see: https://www.w3.org/TR/trace-context/#tracestate-header
func ParseW3Tracestate(tracestate string) (fs Fields, err error) {
	arr := strings.Split(tracestate, "-")
	if len(arr) == 0 {
		return nil, nil
	}
	for _, s := range arr {
		s := strings.TrimSpace(s)
		if len(s) == 0 {
			continue
		}
		kv := strings.SplitN(s, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid state item")
		}
		f := field.New(kv[0])
		field.SetString(f, kv[1]) // TODO 推断值类型?
		fs = append(fs, f)
	}
	return
}
