package field

import (
	"testing"

	"github.com/junhwong/goost/buffering"
	"github.com/stretchr/testify/assert"
)

func TestMarshalJson(t *testing.T) {
	f := Make("hello").SetString("world")
	m := JsonMarshaler{}
	b := buffering.GetBuffer()
	_, err := m.Marshal(f, b)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, `"world"`, b.String())
	f.SetArray([]*Field{Make("hello2").SetString("world1"), Make("hello2").SetString("world2")})
	b.Reset()
	if _, err := m.Marshal(f, b); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, `["world1","world2"]`, b.String())
	f.SetGroup([]*Field{Make("hello1").SetString("world1"), Make("hello2").SetString("world2")})
	b.Reset()
	if _, err := m.Marshal(f, b); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, `{"hello1":"world1","hello2":"world2"}`, b.String())
}
