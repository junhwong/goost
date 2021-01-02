package apm

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/junhwong/goost/pkg/field/common"
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
	std := Logger{
		Out:       os.Stdout,
		Formatter: new(JsonFormatter),
	}

	std.Debug("here %s", "world", common.Message("bbq"))
}

func TestSpan(t *testing.T) {
	std := Logger{
		Out:       os.Stdout,
		queue:     make(chan *Entry, 1000),
		Formatter: new(JsonFormatter),
	}
	ctx, cancel := context.WithCancel(context.TODO())
	go func() {
		std.Run(ctx.Done())
	}()
	go func() {
		span := std.WithSpan(context.TODO(), WithName("GET /test/abc123"))
		defer span.Finish()
		// https://opentracing.io/docs/getting-started/
		// .StartSpan("hello")
		// .Finish()
		//   5fb53ec94b6d4d4a035451cea68a0e77
		//00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01

		span.Debug("hello span")
	}()

	time.Sleep(time.Second * 2)
	cancel()
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
