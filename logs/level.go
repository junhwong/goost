package logs

import "strings"

// Level type.
type Level uint16

// log levels.
const (
	UnsetLevel   Level = iota // 未明确设置(预留给解析日志使用)
	DebugLevel                // 调试(用于开发时的辅助和调试)
	InfoLevel                 // 普通(正常需要输出给最终用户的信息)
	WarningLevel              // 警告(不符合预期, 但不妨碍系统运行)
	ErrorLevel                // 错误(可恢复性错误，不确定系统后续是否正常工作)
	FatalLevel                // 故障(严重的错误系统无法继续，程序应该挂掉)
	TraceLevel                // 跟踪(用于跟踪系统运行状态,如：sql执行时间)
)

var levelMap = map[Level]string{
	UnsetLevel:   "UNSET",
	DebugLevel:   "DEBUG",
	InfoLevel:    "INFO",
	WarningLevel: "WARNING",
	ErrorLevel:   "ERROR",
	FatalLevel:   "FATAL",
	TraceLevel:   "TRACE",
}

var levelShortMap = map[Level]string{
	UnsetLevel:   "-",
	DebugLevel:   "D",
	InfoLevel:    "I",
	WarningLevel: "W",
	ErrorLevel:   "E",
	FatalLevel:   "F",
	TraceLevel:   "T",
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
		return DebugLevel
	case "I", "INFO":
		return InfoLevel
	case "W", "WARNING", "WARN", "ALERT", "NOTICE":
		return WarningLevel
	case "E", "ERROR", "PANIC", "ERR":
		return ErrorLevel
	case "F", "FATAL", "CRITICAL", "EMERG", "CRIT":
		return FatalLevel
	case "T", "TRACE":
		return TraceLevel
	}
	return UnsetLevel
}
