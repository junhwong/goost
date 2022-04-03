package level

import "strings"

// Level type.
type Level = int

const (
	Unset Level = iota // 未明确设置(预留给解析日志使用)
	Debug              // 调试(用于开发时的辅助和调试)
	Info               // 普通(正常需要输出给最终用户的信息)
	Warn               // 警告(不符合预期, 但不妨碍系统运行)
	Error              // 错误(可恢复性错误，不确定系统后续是否正常工作)
	Fatal              // 故障(严重的错误系统无法继续，程序应该挂掉)
	Trace              // 跟踪(用于跟踪系统运行状态,如：sql执行时间)

	_maxLevel
)

var levelMap = map[Level]string{
	Unset: "",
	Debug: "DEBUG",
	Info:  "INFO",
	Warn:  "WARNING",
	Error: "ERROR",
	Fatal: "FATAL",
	Trace: "TRACE",
}

var levelShortMap = map[Level]string{
	Unset: "",
	Debug: "D",
	Info:  "I",
	Warn:  "W",
	Error: "E",
	Fatal: "F",
	Trace: "T",
}

func String(lvl Level) string {
	return levelMap[lvl]
}
func Short(lvl Level) string {
	return levelShortMap[lvl]
}

// ParseLevel return a Level from given string.
func Parse(s string) Level {
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

func FromInt(v int) Level {
	if v <= Unset || v >= _maxLevel {
		return Unset
	}
	return v
}
