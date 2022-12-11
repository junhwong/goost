package field

import (
	"fmt"
	"regexp"
	"sort"
	"sync"
)

// Key 表示一个统一规范的 key/value 结构的键名称。
// 如：
//
//	entry={message=logmessage, os.name=linux, os.type=amd64, trace.error=true}
type Key interface {
	Name() string
	Kind() KeyKind
}

type key struct {
	name   string
	kind   KeyKind
	sec    KeyKind // 子类型, 仅 kind=slice 有效
	reason string
}

func (k key) Name() string {
	return k.name
}
func (k key) Reason() string {
	return k.reason
}
func (k key) Kind() KeyKind {
	return k.kind
}
func (k key) String() string {
	return fmt.Sprintf("field.Key(%s: %s)", k.name, kindNames[k.kind])
}

var (
	keys = sync.Map{}
	// IsValidKeyName 判断给定的名称是否是合法的。
	//
	// `Key name` 主要参考主流的存储设备来定义，如：ES
	IsValidKeyName = regexp.MustCompilePOSIX(`^[a-zA-Z_][a-zA-Z0-9_\-]{0,}(\.[a-zA-Z][a-zA-Z0-9_\-]{0,}){0,}$`).MatchString
)

func GetKey(name string) Key {
	val, ok := keys.Load(name)
	if !ok || val == nil {
		return nil
	}
	key, _ := val.(*key)
	return key
}

func makeOrGetKey(name string, kind KeyKind, sec ...KeyKind) Key {
	if !IsValidKeyName(name) {
		k := &key{name: name, kind: InvalidKind, reason: fmt.Sprintf("Invalid key name: %s", name)}
		return k
	}
	if kind <= InvalidKind || kind > DynamicKind {
		k := &key{name: name, kind: InvalidKind, reason: fmt.Sprintf("Out of range KeyKind: %v", kind)}
		return k
	}
	obj, _ := keys.LoadOrStore(name, &key{name: name, kind: kind})
	if key, _ := obj.(*key); key != nil && key.kind == kind {
		return key
	}
	k := &key{name: name, kind: InvalidKind, reason: fmt.Sprintf("Key already exists, but is not a %s: %s", kindNames[kind], obj)}
	return k
}

type Keys []Key

func (x Keys) Len() int           { return len(x) }
func (x Keys) Less(i, j int) bool { return x[i].Name() < x[j].Name() }
func (x Keys) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x Keys) Sort()              { sort.Sort(x) }