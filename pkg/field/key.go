package field

import (
	"fmt"
	"regexp"
	"sync"
)

// KeyKind 表示 key 的数据类型。
type KeyKind uint

const (
	// InvalidKind 表示无效的字段，将被忽略
	InvalidKind KeyKind = iota
	// StringKind 表示一个字符串
	StringKind
	// SliceKind 表示一个数组
	SliceKind
	// MapKind 表示一个嵌套对象
	MapKind
	// IntKind 表示一个整数(int,int8,int16,int32,int64)=int64
	IntKind
	// UintKind 表示一个无符号整数(uint,uint8,uint16,uint32,uint64,byte)=uint64
	UintKind
	// FloatKind 表示一个浮点数(float32,float64)=float64
	FloatKind
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
	SliceKind:   "slice",
	IntKind:     "int",
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
	obj, _ := keys.LoadOrStore(name, &key{name: name, kind: kind})
	if key := obj.(*key); key != nil && key.kind == kind {
		return key
	}
	panic(fmt.Errorf("Key already exists, but is not a %s: %s", kindNames[kind], obj))
}
