package field

import (
	"testing"

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
			name: "foo.bar0",
			pass: true,
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
			if IsValidKeyName(tC.name) != tC.pass {
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
	k := makeOrGetKey(name, IntKind)
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
func BenchmarkStruct2(b *testing.B) {
	_, f := Int("BenchmarkStruct2")
	fs := make(Fields)
	for i := 0; i < b.N; i++ {
		fs.Set(f(i))
	}
}
func BenchmarkFunc2(b *testing.B) {
	_, f := makeInt("BenchmarkFunc2")
	fs := make(Fields2)
	for i := 0; i < b.N; i++ {
		fs.Set(f(i))
	}
}
