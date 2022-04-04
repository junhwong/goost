package field

// KeyKind 表示 key 的数据类型。
type KeyKind uint

const (
	InvalidKind KeyKind = iota // 无效的字段，将被忽略
	StringKind                 // 字符串 string
	IntKind                    // 整数 int64
	UintKind                   // 整数 uint64
	FloatKind                  // 浮点数 float64
	BoolKind                   // 布尔值 bool
	TimeKind                   // 时间 time.Time
	SliceKind                  // 数组
	MapKind                    // 嵌套对象
	DynamicKind                // 动态字段。警告：该类型的key是不被检查的。
)

var kindNames = map[KeyKind]string{
	InvalidKind: "<invalid>",
	StringKind:  "string",
	IntKind:     "int",
	UintKind:    "uint",
	FloatKind:   "float",
	TimeKind:    "time",
	SliceKind:   "slice",
	MapKind:     "map",
}
