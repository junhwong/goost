package apm

import "strings"

// LogLevel type.
type LogLevel int

const (
	Unset LogLevel = iota // 未明确设置(预留给解析日志使用)
	Debug                 // 调试(用于开发时的辅助和调试)
	Info                  // 普通(正常需要输出给最终用户的信息)
	Warn                  // 警告(不符合预期, 但不妨碍系统运行)
	Error                 // 错误(可恢复性错误，不确定系统后续是否正常工作)
	Fatal                 // 故障(严重的错误系统无法继续，程序应该挂掉)
	Trace                 // 跟踪(用于跟踪系统运行状态,如：sql执行时间)
)

var levelMap = map[LogLevel]string{
	Unset: "DEBUG",
	Debug: "DEBUG",
	Info:  "INFO",
	Warn:  "WARNING",
	Error: "ERROR",
	Fatal: "FATAL",
	Trace: "TRACE",
}

var levelShortMap = map[LogLevel]string{
	Unset: "D",
	Debug: "D",
	Info:  "I",
	Warn:  "W",
	Error: "E",
	Fatal: "F",
	Trace: "T",
}

func (lvl LogLevel) String() string {
	return levelMap[lvl]
}

func (lvl LogLevel) Short() string {
	return levelShortMap[lvl]
}

// ParseLevel return a Level from given string.
func ParseLevel(s string) LogLevel {
	switch strings.TrimSpace(strings.ToUpper(s)) {
	case "D", "DEBUG":
		return Debug
	case "I", "INFO":
		return Info
	case "W", "WARNING", "WARN", "ALERT", "NOTICE":
		return Warn
	case "E", "ERROR", "PANIC", "ERR":
		return Error
	case "F", "FATAL", "CRITICAL", "EMERG", "CRIT":
		return Fatal
	case "T", "TRACE":
		return Trace
	}
	return Unset
}

func LevelFromInt(v int) LogLevel {
	if v <= int(Unset) || v > int(Trace) {
		return Unset
	}
	return LogLevel(v)
}
