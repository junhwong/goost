package apm

import (
	"context"
	"os"
	"time"

	"github.com/junhwong/goost/runtime"
)

// var (
// 	Root RootLogger = &DefaultLogger{Prefix: "default"}
// )

// // ############### logger interface methods ###############

// func Debug(args ...interface{}) { Root.Debug(args...) }
// func Info(args ...interface{})  { Root.Info(args...) }
// func Warn(args ...interface{})  { Root.Warn(args...) }
// func Error(args ...interface{}) { Root.Error(args...) }
// func Fatal(args ...interface{}) { Root.Fatal(args...) }

// // ############### help methods ###############

// // Crash calls fmt.Fprintln and debug.PrintStack() to print to the stderr.
// // followed by a call to os.Exit(1).
// //
// // Note: this method not logging message to logger.
// func Crash(v ...interface{}) {
// 	fmt.Fprintln(os.Stderr, v...)
// 	debug.PrintStack()
// 	os.Exit(1)
// }
var std = Logger{
	Out:       os.Stdout,
	queue:     make(chan *Entry, 1000),
	Formatter: new(JsonFormatter),
}

func init() {
	ctx, cancel := context.WithCancel(context.TODO())
	std = Logger{
		Out:       os.Stdout,
		queue:     make(chan *Entry, 1000),
		Formatter: new(JsonFormatter),
		cancel:    cancel,
	}
	go std.Run(ctx.Done())
}
func Run(stopCh runtime.StopCh) {
	std.Run(stopCh)
}
func Done() {
	std.cancel()
	time.Sleep(time.Second) // TODO 不延迟处理

	std.flush()
}
