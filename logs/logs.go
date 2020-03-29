package logs

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
