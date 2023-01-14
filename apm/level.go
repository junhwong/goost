package apm

import "strings"

// Level type.
type Level int

const (
	LevelUnset Level = iota // 未明确设置(预留给解析日志使用)
	LevelDebug              // 调试(用于开发时的辅助和调试)
	LevelInfo               // 普通(正常需要输出给最终用户的信息)
	LevelWarn               // 警告(不符合预期, 但不妨碍系统运行)
	LevelError              // 错误(可恢复性错误，不确定系统后续是否正常工作)
	LevelFatal              // 故障(严重的错误系统无法继续，程序应该挂掉)
	LevelTrace              // 跟踪(用于跟踪系统运行状态,如：sql执行时间)
)

var levelMap = map[Level]string{
	LevelUnset: "DEBUG",
	LevelDebug: "DEBUG",
	LevelInfo:  "INFO",
	LevelWarn:  "WARNING",
	LevelError: "ERROR",
	LevelFatal: "FATAL",
	LevelTrace: "TRACE",
}

var levelShortMap = map[Level]string{
	LevelUnset: "D",
	LevelDebug: "D",
	LevelInfo:  "I",
	LevelWarn:  "W",
	LevelError: "E",
	LevelFatal: "F",
	LevelTrace: "T",
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
	case "D", "DEBUG":
		return LevelDebug
	case "I", "INFO":
		return LevelInfo
	case "W", "WARNING", "WARN", "ALERT", "NOTICE":
		return LevelWarn
	case "E", "ERROR", "PANIC", "ERR":
		return LevelError
	case "F", "FATAL", "CRITICAL", "EMERG", "CRIT":
		return LevelFatal
	case "T", "TRACE":
		return LevelTrace
	}
	return LevelUnset
}

func LevelFromInt(v int) Level {
	if v <= int(LevelUnset) || v > int(LevelTrace) {
		return LevelUnset
	}
	return Level(v)
}
