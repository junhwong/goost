package logs

import "strings"

// Level type.
type Level uint16

// log levels.
const (
	UNSET Level = iota // 未明确设置(预留给解析日志使用)
	DEBUG              // 调试(用于开发时的辅助和调试)
	INFO               // 普通(正常需要输出给最终用户的信息)
	WARN               // 警告(不符合预期, 但不妨碍系统运行)
	ERROR              // 错误(可恢复性错误，不确定系统后续是否正常工作)
	FATAL              // 故障(严重的错误系统无法继续，程序应该挂掉)
	TRACE              // 跟踪(用于跟踪系统运行状态,如：sql执行)
)

var levelMap = map[Level]string{
	UNSET: "UNSET",
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	FATAL: "FATAL",
	TRACE: "TRACE",
}

var levelShortMap = map[Level]string{
	UNSET: "U",
	DEBUG: "D",
	INFO:  "I",
	WARN:  "W",
	ERROR: "E",
	FATAL: "F",
	TRACE: "T",
}

func (lvl Level) String() string {
	return levelMap[lvl]
}
func (lvl Level) Short() string {
	return levelShortMap[lvl]
}

// ParseLevel return a Level from given string.
func ParseLevel(s string) Level {
	switch strings.TrimSpace(strings.ToUpper(s)) {
	case "DEBUG", "D":
		return DEBUG
	case "INFO", "I":
		return INFO
	case "WARN", "W", "WARNING", "ALERT", "NOTICE":
		return WARN
	case "ERROR", "E", "PANIC", "ERR":
		return ERROR
	case "FATAL", "F", "CRITICAL", "EMERG", "CRIT":
		return FATAL
	case "TRACE", "T":
		return TRACE
	}
	return UNSET
}
