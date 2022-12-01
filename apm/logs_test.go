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
	std.Log(context.TODO(), 0, level.Debug, []interface{}{"hello"})

	std.Close()

}

func xerr() error {
	return errors.WithTraceback(fmt.Errorf("mock err"))
}

func TestSpan(t *testing.T) {
	t.Cleanup(Done)
	_, span := std.NewSpan(context.TODO(), 0)
	defer span.End()
	err := xerr()
	span.Error(err)
	span.Debug("hello span")
	Default().Debugf("hhh%v", 2)
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

func BenchmarkAccumulatedContext(b *testing.B) {
	// b.Logf("Logging with some accumulated context.")
	b.Run("goost/apm", func(b *testing.B) {
		// logger := newZapLogger(zap.DebugLevel).With(fakeFields()...)
		logger, std := newTestLog()
		b.Cleanup(func() {
			// time.Sleep(time.Second)
			std.Close()
		})
		b.ResetTimer()
		if logger != nil {
		}
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// logger.Info("getMessage(0)", field.Dynamic("")(""))

				// std.Write(level.Debug, time.Now(), "", runtime.Caller(0))
			}
		})
	})
}

func newTestLog() (Logger, *DefaultLogger) {
	ctx, cancel := context.WithCancel(context.Background())
	f := &TextFormatter{
		timeLayout: "20060102 15:04:05.000",
	}
	std := &DefaultLogger{
		queue:    make(chan Entry, 1024),
		handlers: []Handler{&ConsoleHandler{Formatter: f}},
		cancel:   cancel,
	}
	go std.Run(ctx.Done())
	r := &stdImpl{entryLog: entryLog{ctx: ctx, logger: std}}
	r.spi = std
	r.calldepth = 1
	return r, std
}

type tf func(interface{})
type thand func(interface{}, tf)

func TestHanc(t *testing.T) {

	hds := []thand{
		func(i interface{}, t tf) {
			fmt.Println("1:", i)
			if t != nil {
				t(i)
			}
		},
		func(i interface{}, t tf) {
			fmt.Println("2:", i)
			if t != nil {
				t(i)
			}
		},
		func(i interface{}, t tf) {
			fmt.Println("3:", i)
			if t != nil {
				t(i)
			}
		},
	}

	var r func(interface{})
	var ep func(interface{})
	for i, t2 := range hds {
		prev := r
		var next func(interface{}) = func(ent interface{}) {
			if prev != nil {
				prev(ent)
			}
		}
		if i+1 < len(hds) {
			// fmt.Println("", i)
			nh := hds[i+1]
			next = func(ent interface{}) {
				nh(ent, prev)
			}
		}

		r = func(ent interface{}) {
			t2(ent, next)
		}
		if ep == nil {
			// r("bb")
			ep = r
		}
	}
	hds = nil
	ep("abc")
}

func TestHanc2(t *testing.T) {

	hds := []thand{
		func(i interface{}, t tf) {
			fmt.Println("1:", i)
			if t != nil {
				t(i)
			}
		},
		func(i interface{}, t tf) {
			fmt.Println("2:", i)
			if t != nil {
				t(i)
			}
		},
		func(i interface{}, t tf) {
			fmt.Println("3:", i)
			if t != nil {
				t(i)
			}
		},
		func(i interface{}, t tf) {
			fmt.Println("4:", i)
			if t != nil {
				t(i)
			}
		},
	}

	var r func(interface{})
	// var ep func(interface{})

	for _, t2 := range hds {
		r = cb(r, t2)
	}
	hds = nil
	r("abc")
}

func cb(n tf, c thand) tf {
	return func(ent interface{}) {
		c(ent, n)
	}

}
