package field

import (
	"fmt"
	"regexp"
	"sync"
)

// KeyKind 表示 key 的数据类型。
type KeyKind uint

const (
	InvalidKind KeyKind = iota // 表示无效的字段，将被忽略
	StringKind                 // 表示一个字符串
	SliceKind                  // 表示一个数组
	MapKind                    // 表示一个嵌套对象
	IntKind                    // 表示一个整数 int64
	FloatKind                  // 表示一个浮点数(float32,float64)=float64
	// // UintKind 表示一个无符号整数(uint,uint8,uint16,uint32,uint64,byte)=uint64
	// UintKind
	StringsKind // 表示一个数组
	IntsKind    // 表示一个数组
	FloatsKind  // 表示一个数组
)

var (
	keys = sync.Map{}
	// IsValidKeyName 判断给定的名称是否是合法的。
	//
	// `Key name` 主要参考主流的存储设备来定义，如：ES
	IsValidKeyName = regexp.MustCompilePOSIX(`^[a-zA-Z][a-zA-Z0-9_\-]{0,}(\.[a-zA-Z][a-zA-Z0-9_\-]{0,}){0,}$`).MatchString
)
var kindNames = map[KeyKind]string{
	InvalidKind: "",
	StringKind:  "string",
	StringsKind: "strings",
	SliceKind:   "slice",
	MapKind:     "map",
	IntKind:     "int",
	FloatKind:   "float",
}

// Key 表示一个统一规范的 key/value 结构的键名称。
// 如：
//    log entry={message=logmessage, os.name=linux, os.type=amd64, error=true}
type Key interface {
	Name() string
	Kind() KeyKind
}

type key struct {
	name string
	kind KeyKind
}

func (k *key) Name() string {
	return k.name
}
func (k *key) Kind() KeyKind {
	return k.kind
}
func (k *key) String() string {
	return fmt.Sprintf("field.Key(%s: %s)", k.name, kindNames[k.kind])
}

func makeOrGetKey(name string, kind KeyKind) Key {
	if !IsValidKeyName(name) {
		panic(fmt.Errorf("Invalid key name: %s", name))
	}
	if kind <= InvalidKind || kind > FloatsKind {
		panic(fmt.Errorf("Out of range KeyKind: %v", kind))
	}
	obj, _ := keys.LoadOrStore(name, &key{name: name, kind: kind})
	if key := obj.(*key); key != nil && key.kind == kind {
		return key
	}
	panic(fmt.Errorf("Key already exists, but is not a %s: %s", kindNames[kind], obj))
}
