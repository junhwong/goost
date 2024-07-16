package apm

import (
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

var ZeroHexID = make(HexID, 16)

// 符合 W3C 规范的 TraceID 或 SpanID.
//
// see: https://www.w3.org/TR/trace-context/#trace-id
type HexID []byte

func (id HexID) Bytes() []byte { return id }
func (id HexID) High() HexID {
	if len(id) != 16 {
		return nil
	}
	return id[:8]
}
func (id HexID) Low() HexID {
	if len(id) != 16 {
		return nil
	}
	return id[8:]
}
func (id HexID) String() string {
	if l := len(id); !(l == 0 || l == 8 || l == 16) {
		return "<invalid>"
	}
	return fmt.Sprintf("%x", id.Bytes())
}
func (id HexID) Equal(b HexID) bool {
	if len(id) != len(b) {
		return false
	}
	for i := range b {
		if b[i] != id[i] {
			return false
		}
	}
	return true
}

var seededIDGen = rand.New(rand.NewSource(time.Now().UnixNano()))
var mu sync.Mutex

// randomTimestamped can generate 128 bit time sortable traceid's compatible
// with AWS X-Ray and 64 bit spanid's.
func NewHexID() HexID {
	mu.Lock()
	var b []byte
	if i := uint64(time.Now().Unix()<<32) + uint64(seededIDGen.Int31()); i > 0 {
		b = binary.BigEndian.AppendUint64(b, i)
	}
	if i := uint64(seededIDGen.Int63()); i > 0 {
		b = binary.BigEndian.AppendUint64(b, i)
	}
	mu.Unlock()
	return b
}

var (
	errInvalidHexID = errors.New("hex-id can only contain hex characters, len (16 or 32)")
)

// ParseHexID returns a HexID from a hex string.
func ParseHexID(h string) (HexID, error) {
	decoded, err := hex.DecodeString(h)
	if err != nil {
		return nil, errInvalidHexID
	}
	switch len(decoded) {
	case 16:
		return ZeroHexID, nil
	case 8:
		decoded = append(make([]byte, 8), decoded...)
	default:
		return nil, errInvalidHexID
	}
	for i := range decoded {
		if decoded[i] != ZeroHexID[i] {
			return decoded, nil
		}
	}

	return ZeroHexID, nil
}

// 解析 W3C trace.
//
// 示例: `00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01`.
//
// see: https://www.w3.org/TR/trace-context/#traceparent-header
func ParseW3Traceparent(traceparent string) (version byte, traceID, parentSpanID string, flags byte, err error) {
	if traceparent == "" {
		return
	}
	arr := strings.Split(traceparent, "-")
	if len(arr) != 4 {
		err = fmt.Errorf("invalid format")
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

	traceID = arr[1]
	parentSpanID = arr[2]

	// traceID, err = ParseHexID(arr[1])
	// if err != nil {
	// 	return
	// }

	// parentSpanID, err = ParseHexID(arr[2])
	// if err != nil {
	// 	return
	// }

	return
}

// 解析 W3C tracestate.
//
// 示例: `rojo=00f067aa0ba902b7,congo=t61rcWkgMzE`.
//
// see: https://www.w3.org/TR/trace-context/#tracestate-header
func ParseW3Tracestate(tracestate string) (fs []*field.Field, err error) {
	arr := strings.Split(tracestate, ",")
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
		f := field.Make(kv[0]).SetString(kv[1]) // TODO 推断值类型?
		if f.GetType() == field.StringKind {
			fs = append(fs, f)
		}
	}
	return
}
