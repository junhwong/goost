package logs

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type FieldKind uint

const (
	FKInvalid   FieldKind = iota // 无效的字段，将被忽略
	FKAny                        // 任意
	FKBool                       // 布尔值
	FKString                     // 简单字符串，符合规则： `a-zA-Z0-9\.\-\:`
	FKInteger                    // 整数。
	FKFloat                      // 浮点数
	FKTime                       // 时间
	FKDuration                   // 时长
	FKTriceback                  // 堆栈
)

var fkmap = map[FieldKind]string{
	FKInvalid:   "invalid",
	FKAny:       "any",
	FKBool:      "bool",
	FKString:    "string",
	FKInteger:   "int",
	FKFloat:     "float",
	FKTime:      "time",
	FKDuration:  "duration",
	FKTriceback: "triceback",
}

func (k FieldKind) String() string {
	return fkmap[k]
}

func CheckSafetyString(v string) bool {
	return false
}

// Field 表示一个日志标签。
//
// 注意：不建议直接使用该结构初始化，应该使用 `logs.String()`, `logs.Int()`等这类辅助方法创建。
type Field struct {
	// 命名规则 `[a-z][a-z0-9_]*`根据木桶原理，由当前市面上比较流行的几种存储和分析软件调整而来。
	//
	// - MySQL and Systemd-Journal.
	//
	// - [Elasticsearch](https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-create-index.html)
	//
	// - [InfluxDB](https://v2.docs.influxdata.com/v2.0/reference/line-protocol/#naming-restrictions)
	//
	// - [Prometheus](https://prometheus.io/docs/concepts/data_model/)
	Name    string
	Value   interface{}
	Index   bool
	NoIndex bool
	Type    FieldKind
}

// MarkIndex 标记为索引
func (f *Field) MarkIndex(force ...bool) *Field {
	switch {
	case f.Index && f.Type == FKInvalid || f.Type == FKTriceback:
	case f.Type != FKString && len(force) > 0 && force[0]:
		f.Value = fmt.Sprint(f.Value)
		f.Type = FKString
		f.Index = true
	case f.Type == FKString:
		f.Index = true
	}
	return f
}

var nameRegex = regexp.MustCompile(`[a-z][a-z0-9_]*`)

func newField(n string, v interface{}, t FieldKind) *Field {
	if v == nil || !nameRegex.MatchString(n) {
		t = FKInvalid
	}
	f := &Field{
		Name:  n,
		Value: v,
		Type:  t,
	}
	return f
}

func Any(name string, value interface{}) *Field {
	switch v := value.(type) {
	case int:
		return newField(name, int64(v), FKInteger)
	case uint:
		return newField(name, int64(v), FKInteger)
	case int16:
		return newField(name, int64(v), FKInteger)
	case uint16:
		return newField(name, int64(v), FKInteger)
	case int32:
		return newField(name, int64(v), FKInteger)
	case uint32:
		return newField(name, int64(v), FKInteger)
	case int64:
		return newField(name, v, FKInteger)
	case uint64:
		return newField(name, int64(v), FKInteger)
	case uint8:
		return newField(name, int64(v), FKInteger)
	case uintptr:
		return newField(name, int64(v), FKInteger)
	case float32:
		f := float64(v)
		if !strings.Contains(strconv.FormatFloat(f, 'f', -1, 64), ".") {
			return newField(name, int64(v), FKInteger)
		} else {
			return newField(name, f, FKFloat)
		}
	case float64:
		if !strings.Contains(strconv.FormatFloat(v, 'f', -1, 64), ".") {
			return newField(name, int64(v), FKInteger)
		} else {
			return newField(name, v, FKFloat)
		}
	case bool:
		return newField(name, v, FKBool)
	case string:
		t := FKString
		if v == "" {
			t = FKInvalid
		}
		return newField(name, v, t)
	case time.Time:
		return newField(name, v, FKTime)
	case time.Duration:
		return newField(name, v, FKDuration)
	default:
		return newField(name, v, FKAny)
	}
}

func String(name, value string) *Field {
	t := FKString
	if value == "" {
		t = FKInvalid
	}
	return newField(name, value, t)
}

func Int(name string, value int64) *Field {
	return newField(name, value, FKInteger)
}

func Float(name string, value float64) *Field {
	return newField(name, value, FKFloat)
}
