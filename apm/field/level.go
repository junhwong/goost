package field

import "strings"

// Level type.
// see https://opentelemetry.io/docs/reference/specification/logs/data-model/#severity-fields
type Level uint

// 重要程度. 由于希望默认行为, 所以未严格按照 OTEL 的级别显示名称, 但值相同.
const (
	LevelUnset Level = 0 // 未明确设置或未知(预留给解析日志使用)

	levelTrace1 Level = 1 // 1-4 建议日志中不包含该区间的级别, 处理其它系统定义时也转义为 Debug1
	LevelTrace  Level = 2 // 特殊: 跟踪
	LevelMetric Level = 3 // 特殊: 指标
	levelTrace4 Level = 4

	levelDebug1 Level = 5 //
	LevelDebug  Level = 6 // 5-8 调试(用于开发时的辅助和调试)
	levelDebug3 Level = 7 // 配置
	levelDebug4 Level = 8 //

	LevelInfo  Level = 9 // 9-12 普通(正常需要输出给最终用户的信息)
	levelInfo2 Level = 10
	levelInfo3 Level = 11
	levelInfo4 Level = 12

	LevelWarn  Level = 13 // 13-16 警告(不符合预期, 但不妨碍系统运行)
	levelWarn2 Level = 14
	levelWarn3 Level = 15
	levelWarn4 Level = 16

	LevelError  Level = 17 // 17-20 错误(可恢复性错误，不确定系统后续是否正常工作)
	levelError2 Level = 18
	levelError3 Level = 19
	levelError4 Level = 20

	LevelFatal  Level = 21 // 21-24 故障(严重的错误系统无法继续，程序应该挂掉)
	levelFatal2 Level = 22
	levelFatal3 Level = 23
	levelFatal4 Level = 24
)

var levelMap = map[Level]string{

	levelTrace1: "trace1",
	LevelTrace:  "trace",
	LevelMetric: "metric",
	levelTrace4: "trace4",

	levelDebug1: "debug1",
	LevelDebug:  "debug",
	levelDebug3: "debug3",
	levelDebug4: "debug4",

	LevelInfo:  "info",
	levelInfo2: "info2",
	levelInfo3: "info3",
	levelInfo4: "info4",

	LevelWarn:  "warn",
	levelWarn2: "warn1",
	levelWarn3: "warn2",
	levelWarn4: "warn3",

	LevelError:  "error",
	levelError2: "error2",
	levelError3: "error3",
	levelError4: "error4",

	LevelFatal:  "fatal",
	levelFatal2: "fatal2",
	levelFatal3: "fatal3",
	levelFatal4: "fatal4",
}

var levelShortMap = map[Level]string{
	levelTrace1: "T1",
	LevelTrace:  "T",
	LevelMetric: "M",
	levelTrace4: "T4",

	levelDebug1: "D1",
	LevelDebug:  "D",
	levelDebug3: "D3",
	levelDebug4: "D4",

	LevelInfo:  "I",
	levelInfo2: "I2",
	levelInfo3: "I3",
	levelInfo4: "I4",

	LevelWarn:  "W",
	levelWarn2: "W1",
	levelWarn3: "W2",
	levelWarn4: "W3",

	LevelError:  "E",
	levelError2: "E2",
	levelError3: "E3",
	levelError4: "E4",

	LevelFatal:  "F",
	levelFatal2: "F2",
	levelFatal3: "F3",
	levelFatal4: "F4",
}

func (lvl Level) String() string {
	return levelMap[lvl]
}

func (lvl Level) Short() string {
	return levelShortMap[lvl]
}

// ParseLevel return a Level from given string.
// see https://opentelemetry.io/docs/reference/specification/logs/data-model-appendix/#appendix-b-severitynumber-example-mappings
func ParseLevel(s string) Level {
	s = strings.ReplaceAll(strings.TrimSpace(s), "LogLevel.", "")
	switch strings.ToUpper(s) {
	case "T1", "FINEST":
		return levelTrace1
	case "T2", "T", "TRACE":
		return LevelTrace
	case "T3", "M", "METRIC":
		return LevelMetric
	case "T4":
		return levelTrace4
	case "D1":
		return levelDebug1
	case "D2", "D", "DEBUG", "FINER", "FINE":
		return LevelDebug
	case "D3", "CONFIG":
		return levelDebug3
	case "D4":
		return levelDebug4
	case "I1", "I", "INFO", "INFORMATION":
		return LevelInfo
	case "I2":
		return levelInfo2
	case "I3":
		return levelInfo3
	case "I4":
		return levelInfo4
	case "W1", "W", "WARNING", "WARN", "NOTICE":
		return LevelWarn
	case "W2":
		return levelWarn2
	case "W3":
		return levelWarn3
	case "W4":
		return levelWarn4
	case "E1", "E", "ERROR", "SEVERE", "ERR":
		return LevelError
	case "E2", "CRITICAL", "CRIT":
		return levelError2
	case "E3", "ALERT", "PANIC":
		return levelError3
	case "E4":
		return levelError4
	case "F1", "F", "FATAL", "EMERG":
		return LevelFatal
	case "F2":
		return levelFatal2
	case "F3":
		return levelFatal3
	case "F4":
		return levelFatal4
	}
	return LevelUnset
}

func LevelFromInt(v int) Level {
	if v <= 0 || v > 24 {
		return LevelUnset
	}
	return Level(v)
}
