package logs

import (
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
	span := std.SpanWithContext(nil)
	// https://opentracing.io/docs/getting-started/
	// .StartSpan("hello")
	// .Finish()
	span.Debug()
}

func TestX(t *testing.T) {
	var v interface{} = 12
	z, _ := v.(string)
	t.Log(z)
}
