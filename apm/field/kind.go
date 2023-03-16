package field

import (
	"github.com/junhwong/goost/apm/field/pb"
)

// KeyKind 表示 key 的数据类型。
type KeyKind = pb.Field_Type

const (
	InvalidKind  = pb.Field_UNKNOWN   // 无效的字段，将被忽略
	StringKind   = pb.Field_STRING    // 字符串 string
	IntKind      = pb.Field_INT       // 整数 int64
	UintKind     = pb.Field_UINT      // 整数 uint64
	FloatKind    = pb.Field_FLOAT     // 浮点数 float64
	BoolKind     = pb.Field_BOOL      // 布尔值 bool
	TimeKind     = pb.Field_TIMESTAMP // 时间 time.Time
	DurationKind = pb.Field_DURATION  // 时间 time.Duration
	BytesKind    = pb.Field_BYTES     // bytes

	// SliceKind = pb.Field_UNKNOWN // 数组
	// MapKind                     // 嵌套对象
	// DynamicKind = pb.Field_UNKNOWN // 动态字段。警告：该类型的key是不被检查的。
)

var kindNames = map[KeyKind]string{
	InvalidKind: "<invalid>",
	StringKind:  "string",
	IntKind:     "int",
	UintKind:    "uint",
	FloatKind:   "float",
	BoolKind:    "bool",
	TimeKind:    "time",

	// SliceKind:   "slice",
	// MapKind:     "map",
	// DynamicKind: "dynamic",
}

// func (k KeyKind) String() string {
// 	s, ok := kindNames[k]
// 	if !ok {
// 		return strconv.Itoa(int(k))
// 	}
// 	return s
// }
