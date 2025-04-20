package field

import (
	"fmt"
	"testing"
	"time"

	"github.com/junhwong/goost/buffering"
)

func TestWrie(t *testing.T) {
	f := Make("a").SetString("b")
	buf := buffering.GetBuffer()
	err := Marshal(f, buf, true)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("data: %v\n", buf.Bytes())
	f2 := Make("ttt")
	err = Unmarshal(buf, f2)
	if err != nil {
		t.Fatal(err)
	}
	s1 := fmt.Sprintf("%#v", f)
	s2 := fmt.Sprintf("%#v", f2)
	if s1 != s2 {
		t.Fatal(s1, s2)
	}
	fmt.Printf("s1: %v\n", s1)
	fmt.Printf("s2: %v\n", s2)
}

func TestMarshalUnmarshal(t *testing.T) {
	// 初始化缓冲区
	buf := buffering.GetBuffer()

	// 测试用例：包含所有支持类型的Field
	field := &Field{
		kind:    GroupKind,
		typFlag: 0,
		Items: []*Field{
			{kind: StringKind, name: "str", strVal: "hello"},
			{kind: IntKind, name: "int", intVal: -42},
			{kind: UintKind, name: "uint", uintVal: 42},
			{kind: FloatKind, name: "float", floatVal: 3.14},
			{kind: BoolKind, name: "bool", intVal: 1},
			{kind: TimeKind, name: "time", intVal: time.Now().UnixNano()},
			{kind: DurationKind, name: "duration", intVal: int64(time.Hour)},
			{kind: BytesKind, name: "bytes", bytesVal: []byte{0xDE, 0xAD}},
			{kind: IPKind, name: "ip", bytesVal: []byte{127, 0, 0, 1}},
			{kind: LevelKind, name: "level", intVal: 2},
			{
				kind: ArrayKind,
				Items: []*Field{
					{kind: IntKind, intVal: 1},
					{kind: IntKind, intVal: 2},
				},
			},
		},
	}

	// 序列化
	if err := Marshal(field, buf, true); err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	newField := &Field{}
	if err := Unmarshal(buf, newField); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证反序列化后的数据
	// 这里应添加详细的字段比对逻辑，此处仅示例
	if newField.Items[0].GetString() != "hello" {
		t.Error("String field mismatch")

	}

	s1 := fmt.Sprintf("%#v\n", field)
	s2 := fmt.Sprintf("%#v\n", newField)
	if s1 != s2 {
		t.Logf("s1: %s\n", s1)
		t.Logf("s2: %s\n", s2)
		t.Error("String field mismatch")
	}
}

func TestInvalidType(t *testing.T) {
	buf := buffering.GetBuffer()
	field := &Field{kind: Type(255)} // 无效类型

	if err := Marshal(field, buf, true); err == nil {
		t.Error("Expected error for invalid type")
	}
}
