package logs

import (
	"context"
	"log"
	"os"
	"testing"

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
		Formatter: new(JsonFormatter),
	}
	span := std.WithSpan(context.TODO(), "test/abc123")
	defer span.Finish()
	// https://opentracing.io/docs/getting-started/
	// .StartSpan("hello")
	// .Finish()

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
