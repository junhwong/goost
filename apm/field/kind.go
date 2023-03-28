package field

import (
	"strings"
)

// Kind 数据类型。
type Kind = Type

const (
	InvalidKind  = Type_UNKNOWN
	StringKind   = Type_STRING
	IntKind      = Type_INT
	UintKind     = Type_UINT
	FloatKind    = Type_FLOAT
	BoolKind     = Type_BOOL
	TimeKind     = Type_TIMESTAMP
	DurationKind = Type_DURATION
	BytesKind    = Type_BYTES
	IPKind       = Type_IP
	LevelKind    = Type_LOGLEVEL

	// SliceKind = pb.Field_UNKNOWN // 数组
	// MapKind                     // 嵌套对象
	// DynamicKind = pb.Field_UNKNOWN // 动态字段。警告：该类型的key是不被检查的。
)

var kindNames = map[Kind]string{
	InvalidKind:  "<invalid>",
	StringKind:   "string",
	IntKind:      "int",
	UintKind:     "uint",
	FloatKind:    "float",
	BoolKind:     "bool",
	TimeKind:     "time",
	DurationKind: "duration",
	BytesKind:    "bytes",
	IPKind:       "ip",
	LevelKind:    "level",

	// SliceKind:   "slice",
	// MapKind:     "map",
	// DynamicKind: "dynamic",
}

//	func (k KeyKind) String() string {
//		s, ok := kindNames[k]
//		if !ok {
//			return strconv.Itoa(int(k))
//		}
//		return s
//	}
func ParseType(v any) Kind {
	var n int32
	switch v := v.(type) {
	case string:
		v = strings.ToLower(v)
		switch v {
		case "timestamp", "time", "datetime":
			v = "timestamp"
		case "long", "short":
			v = "int"
		case "number":
			v = "float"
		}
		n = int32(ParseType(Type_value[strings.ToUpper(v)]))
	case int:
		n = int32(v)
	case int32:
		n = v
	}
	if _, ok := Type_name[n]; ok {
		return Kind(n)
	}

	return InvalidKind
}
