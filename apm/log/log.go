package log

import "github.com/junhwong/goost/apm"

var log = apm.Default(apm.WithCallDepthAdd(1))

func Debug(a ...interface{}) { log.Debug(a...) }
func Info(a ...interface{})  { log.Info(a...) }
func Warn(a ...interface{})  { log.Warn(a...) }
func Error(a ...interface{}) { log.Error(a...) }
func Fatal(a ...interface{}) { log.Fatal(a...) }
