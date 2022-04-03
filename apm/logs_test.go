package apm

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/junhwong/goost/apm/level"
	"github.com/junhwong/goost/errors"
)

// func TestLog(t *testing.T) {
// 	root.New().Debug("hello")
// 	root.New().Info("hello")
// }

// func TestBuildTemplete(t *testing.T) {
// 	testCases := []struct {
// 		desc string
// 	}{
// 		{
// 			desc: "${level|-5s} [${time|yyyy-MM-ddTHH:mm:ssS}] ${message}",
// 		},
// 		{
// 			desc: "${level|-5s} [${time|2006-01-02T15:04:05.999Z}] ${message}",
// 		},
// 		{
// 			desc: "${data}",
// 		},
// 	}
// 	for _, tC := range testCases {
// 		t.Run(tC.desc, func(t *testing.T) {
// 			buildTemplete(tC.desc)
// 		})
// 	}
// }

func TestLog(t *testing.T) {
	// std := Logger{
	// 	queue: make(chan *LogEntry, 1000),
	// }
	std.Log(context.TODO(), 2, level.Debug, []interface{}{"here %s", "world", _entryMessage("bbq")})

	std.Close()
}

func xerr() error {
	return errors.WithTraceback(fmt.Errorf("mock err"))
}

func TestSpan(t *testing.T) {
	if std == nil {
		t.Fatal("empty")
	}
	t.Cleanup(Done)
	_, span := std.NewSpan(context.TODO(), WithName("GET /test/abc123"))
	defer span.End()
	err := xerr()
	span.Error(err)
	span.Debug("hello span")
}

func TestX(t *testing.T) {
	var v interface{} = 12
	z, _ := v.(string)
	t.Log(z)
	log.SetPrefix("[abc]")
	log.Println("hhhh")
}

func TestIDgen(t *testing.T) {
	t.Log(newTraceId())
	t.Log(newTraceId())
	t.Log(newTraceId())
	t.Log(newTraceId())
	t.Log(newTraceId())
	t.Log(newTraceId())
}

func TestY(t *testing.T) {
}
