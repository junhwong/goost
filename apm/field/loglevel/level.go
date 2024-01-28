package loglevel

import (
	"strconv"
	"strings"
)

// Level type.
// see https://opentelemetry.io/docs/reference/specification/logs/data-model/#severity-fields
type Level uint

// 重要程度. 由于希望默认行为, 所以未严格按照 OTEL 的级别显示名称, 但值相同.
const (
	Unset Level = 0 // 未明确设置或未知(预留给解析日志使用)

	Trace  Level = 1 // 1-4 建议日志中不包含该区间的级别, 处理其它系统定义时也转义为 Debug1
	Trace2 Level = 2 // 特殊: 跟踪
	Metric Level = 3 // 特殊: 指标
	Trace4 Level = 4

	Debug  Level = 5 // 5-8 调试(用于开发时的辅助和调试)
	Debug2 Level = 6 //
	Debug3 Level = 7 // 配置
	Debug4 Level = 8 //

	Info  Level = 9 // 9-12 普通(正常需要输出给最终用户的信息)
	Info2 Level = 10
	Info3 Level = 11
	Info4 Level = 12

	Warn  Level = 13 // 13-16 警告(不符合预期, 但不妨碍系统运行)
	Warn2 Level = 14
	Warn3 Level = 15
	Warn4 Level = 16

	Error  Level = 17 // 17-20 错误(可恢复性错误，不确定系统后续是否正常工作)
	Error2 Level = 18
	Error3 Level = 19
	Error4 Level = 20

	Fatal  Level = 21 // 21-24 故障(严重的错误系统无法继续，程序应该挂掉)
	Fatal2 Level = 22
	Fatal3 Level = 23
	Fatal4 Level = 24
)

var levelMap = map[Level]string{

	Trace:  "trace1",
	Trace2: "trace",
	Metric: "metric",
	Trace4: "trace4",

	Debug2: "debug1",
	Debug:  "debug",
	Debug3: "debug3",
	Debug4: "debug4",

	Info:  "info",
	Info2: "info2",
	Info3: "info3",
	Info4: "info4",

	Warn:  "warn",
	Warn2: "warn1",
	Warn3: "warn2",
	Warn4: "warn3",

	Error:  "error",
	Error2: "error2",
	Error3: "error3",
	Error4: "error4",

	Fatal:  "fatal",
	Fatal2: "fatal2",
	Fatal3: "fatal3",
	Fatal4: "fatal4",
}

var levelShortMap = map[Level]string{
	Trace:  "T1",
	Trace2: "T",
	Metric: "M",
	Trace4: "T4",

	Debug2: "D1",
	Debug:  "D",
	Debug3: "D3",
	Debug4: "D4",

	Info:  "I",
	Info2: "I2",
	Info3: "I3",
	Info4: "I4",

	Warn:  "W",
	Warn2: "W1",
	Warn3: "W2",
	Warn4: "W3",

	Error:  "E",
	Error2: "E2",
	Error3: "E3",
	Error4: "E4",

	Fatal:  "F",
	Fatal2: "F2",
	Fatal3: "F3",
	Fatal4: "F4",
}

func (lvl Level) String() string {
	return levelMap[lvl]
}

func (lvl Level) Short() string {
	return levelShortMap[lvl]
}

func FromInt(v int) Level {
	if v <= 0 || v > 24 {
		return Unset
	}
	return Level(v)
}

// Parse return a Level from given string.
// see https://opentelemetry.io/docs/reference/specification/logs/data-model-appendix/#appendix-b-severitynumber-example-mappings
func Parse(s string) Level {
	s = strings.ReplaceAll(strings.TrimSpace(s), "LogLevel.", "")
	switch strings.ToUpper(s) {
	case "T1", "FINEST":
		return Trace
	case "T2", "T", "TRACE":
		return Trace2
	case "T3", "M", "METRIC":
		return Metric
	case "T4":
		return Trace4
	case "D1", "D", "DEBUG", "FINER", "FINE":
		return Debug
	case "D2":
		return Debug2
	case "D3", "CONFIG":
		return Debug3
	case "D4":
		return Debug4
	case "I1", "I", "INFO", "INFORMATION":
		return Info
	case "I2":
		return Info2
	case "I3":
		return Info3
	case "I4":
		return Info4
	case "W1", "W", "WARNING", "WARN", "NOTICE":
		return Warn
	case "W2":
		return Warn2
	case "W3":
		return Warn3
	case "W4":
		return Warn4
	case "E1", "E", "ERROR", "SEVERE", "ERR":
		return Error
	case "E2", "CRITICAL", "CRIT":
		return Error2
	case "E3", "ALERT", "PANIC":
		return Error3
	case "E4":
		return Error4
	case "F1", "F", "FATAL", "EMERG":
		return Fatal
	case "F2":
		return Fatal2
	case "F3":
		return Fatal3
	case "F4":
		return Fatal4
	}
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return Unset
	}
	return FromInt(int(i))
}
