package field

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
)

// Key 表示一个统一规范的 key/value 结构的键名称。
// 如：
//
//	entry={message=logmessage, os.name=linux, os.type=amd64, trace.error=true}
//
// see:
// - https://opentelemetry.io/docs/reference/specification/common/attribute-naming/
// - https://www.w3.org/TR/trace-context/#key
type Key interface {
	Name() string
	Kind() Type
}

type key struct {
	name string
	kind Type
}

func (k key) Name() string {
	return k.name
}

func (k key) Kind() Type {
	return k.kind
}
func (k key) String() string {
	r := Type_name[int32(k.kind)]
	if len(r) == 0 {
		r = "<invalid>"
	}
	return fmt.Sprintf("field.Key(%s: %s)", k.name, r)
}

var (
	keys = sync.Map{}
	// IsValidKey 判断给定的名称是否是合法的。
	//
	// `Key name` 主要参考主流的存储设备来定义，如：ES
	IsValidKeyPatt = regexp.MustCompile(`^[a-zA-Z]([\.\-]?[a-zA-Z]\w*)*$`)
)

func IsValidKey(s string) bool {
	if l := len(s); l > 4 && strings.HasPrefix(s, "__") && strings.HasSuffix(s, "__") {
		return IsValidKeyPatt.MatchString(s[2 : l-2])
	}
	if l := len(s); l > 2 && strings.HasPrefix(s, "@") {
		return IsValidKeyPatt.MatchString(s[1:])
	}
	return IsValidKeyPatt.MatchString(s)
}

func GetKey(name string) Key {
	val, ok := keys.Load(name)
	if !ok || val == nil {
		return nil
	}
	key, _ := val.(*key)
	return key
}

func makeOrGetKey(name string, kind Type, sec ...Type) Key {
	if !IsValidKey(name) {
		panic(fmt.Sprintf("Invalid key name: %s", name))
	}
	obj, _ := keys.LoadOrStore(name, &key{name: name, kind: kind})
	if key, _ := obj.(*key); key != nil && key.kind == kind {
		return key
	}
	panic(fmt.Sprintf("Key already exists, but is not a %s: %s", kind, obj))
}

type Keys []Key

func (x Keys) Len() int           { return len(x) }
func (x Keys) Less(i, j int) bool { return x[i].Name() < x[j].Name() }
func (x Keys) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x Keys) Sort()              { sort.Sort(x) }
