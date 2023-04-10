package field

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/spf13/cast"
)

func TestIsValidKeyName(t *testing.T) {
	testCases := []struct {
		name string
		pass bool
	}{
		{
			name: "foo",
			pass: true,
		},
		{
			name: "rootDir",
			pass: true,
		},
		{
			name: "foo_bar",
			pass: true,
		},
		{
			name: "__foo__",
			pass: true,
		},
		{
			name: "____",
			pass: false,
		},
		{
			name: "__foo",
			pass: false,
		},
		{
			name: "_foo",
			pass: false,
		},
		{
			name: "@foo",
			pass: true,
		},
		{
			name: "foo@",
			pass: false,
		},
		{
			name: "foo.bar0",
			pass: true,
		},
		{
			name: "foo.bar.",
			pass: false,
		},
		{
			name: "foo..bar0",
			pass: false,
		},
		{
			name: "fo_o.b-ar",
			pass: true,
		},
		{
			name: "fo_o.b-ar.oth",
			pass: true,
		},
		{
			name: "foo.123",
			pass: false,
		},
		{
			name: "123.bar",
			pass: false,
		},
		{
			name: ".",
			pass: false,
		},
		{
			name: "-._",
			pass: false,
		},
		{
			name: "*dfd5#",
			pass: false,
		},
		{
			name: "中文",
			pass: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			if IsValidKey(tC.name) != tC.pass {
				t.Fail()
			}
		})
	}
}

type FieldInvoker func() (Key, interface{})
type Fields2 map[Key]interface{}

func (fs Fields2) Set(f FieldInvoker) {
	if f == nil {
		return
	}
	k, v := f()
	if k == nil || v == nil {
		return
	}
	fs[k] = v
}

func makeInt(name string) (Key, func(interface{}) FieldInvoker) {
	k := makeOrGetKey(name, Type_INT)
	return k, func(v interface{}) FieldInvoker {
		v, e := cast.ToInt64E(v)
		return makeResult(k, v, e)
	}
}

func makeResult(k Key, v interface{}, err error) FieldInvoker {
	if err != nil {
		v = nil
	}
	return func() (Key, interface{}) { return k, v }
}

func t() {

	_, v := makeInt("")

	v(123)

}

func BenchmarkStruct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, f := Int("i")
		f(i)
	}
}
func BenchmarkFunc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, f := makeInt("i")
		f(i)
	}
}

//	func BenchmarkStruct2(b *testing.B) {
//		_, f := Int("BenchmarkStruct2")
//		fs := make(Fields)
//		for i := 0; i < b.N; i++ {
//			fs.Set(f(i))
//		}
//	}
func BenchmarkFunc2(b *testing.B) {
	_, f := makeInt("BenchmarkFunc2")
	fs := make(Fields2)
	for i := 0; i < b.N; i++ {
		fs.Set(f(i))
	}
}

type ts string

func TestRV(t *testing.T) {
	testCases := []struct {
		o any
		k Type
	}{
		{
			o: "r",
			k: Type_STRING,
		},
		{
			o: ts("s"),
			k: Type_STRING,
		},
		{
			o: time.Now(),
			k: Type_TIMESTAMP,
		},
		{
			o: time.Second,
			k: Type_DURATION,
		},
	}
	for _, tC := range testCases {
		t.Run(fmt.Sprint(tC.k, ":", tC.o), func(t *testing.T) {
			v, k := InferPrimitiveValueByReflect(reflect.ValueOf(tC.o))
			if tC.k != k {
				t.Fatal(k, v)
			}
		})
	}
}
