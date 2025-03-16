package jsonpath

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		desc string
		err  bool
		ex   string
	}{
		{
			desc: `@@.@@`,
		},
		{
			desc: `@@.*`,
		},
		{
			desc: `@.call(1,2.3,true,false)`,
		},
		{
			desc: `@.call()`,
		},
		{
			desc: `$.fool.call()`,
		},
		{
			desc: `$.call()`,
			err:  true,
		},
		{
			desc: `$.*`,
		},
		{
			desc: `$.*.[].[]`,
		},
		{
			desc: `@.*`,
		},
		{
			desc: `$.foo`,
		},
		// {
		// 	desc: `$.foo.*`, // todo
		// },
		{
			desc: `@foo`,
			err:  true,
		},
		{
			desc: "$.foo[]",
		},
		{
			desc: "foo[].bar",
			err:  true,
		},
		{
			desc: "foo[][0]",
			err:  true,
		},
		{
			desc: "$.foo[6]",
		},
		{
			desc: "$.foo[-1]",
		},
		{
			desc: "$.[1:20]",
		},
		{
			desc: "$.[5]",
		},
		{
			desc: "$.[:]",
		},
		{
			desc: "$.[:].*",
		},
		{
			desc: "$.[:].foo",
		},
		{
			desc: "$.foo[:]",
		},
		{
			desc: "$.foo[:]bar",
			err:  true,
		},
		{
			desc: "$.foo[:][]",
		},
		{
			desc: "$.foo[:][:][:]",
		},
		{
			desc: "$.[::]",
			err:  true,
		},
		{
			desc: "$.[:2]",
		},
		{
			desc: "$.[1:]",
		},
		{
			desc: "$.[0:2]",
			ex:   "$.[:2]",
		},
		{
			desc: "$.[1:-1]",
			ex:   "$.[1:]",
		},
		// {
		// 	desc: "[1:2:3]",
		// },
		{
			desc: "$.foo[bar]",
		},
		{
			desc: "$.foo.[]",
		},
		{
			desc: "$.中文",
		},
		{
			desc: "$.foo",
		},
		{
			desc: "$.foo.bar",
		},
		{
			desc: "$.foo..bar",
		},
		{
			desc: "$.foo.*",
		},
		{
			desc: "$.foo.*.x",
		},
		{
			desc: `$."foo"`,
		},
		{
			desc: `$."foo".bar`,
		},
		{
			desc: `$.foo"bar"`,
			err:  true,
		},
		{
			desc: `$foo`,
			err:  true,
		},
		{
			desc: `$__foo`,
			err:  true,
		},
		{
			desc: `$."foo"bar`,
			err:  true,
		},
		{
			desc: `$.[?()]`,
			err:  true,
		},
		{
			desc: `$.[?(@.foo<5)]`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ex := tC.desc
			r, _, err := Parse(ex)
			if tC.err {
				assert.NotNil(t, err)
				return
			}
			if !assert.Nil(t, err) {
				return
			}
			buf := bytes.NewBuffer(nil)
			buf.Reset()
			v := &printer{Out: buf}
			v.Visit = func(e Expr) {
				Visit(e, v, v.SetError)
			}
			v.Visit(r)
			assert.Nil(t, v.Error())
			if tC.ex != "" {
				ex = tC.ex
			}
			assert.Equal(t, ex, buf.String())
		})
	}
}

func TestTT(t *testing.T) {
	s := `.[]()?@"'\n\t\r中`
	rs := []rune(s)
	for _, v := range rs {
		if v < '1' {
		}
		fmt.Printf("v: %v\n", v)
	}
	i := strings.IndexAny(s, "文")
	fmt.Printf("i: %v\n", i)
	fmt.Printf("s[i:]: %s\n", []byte{s[i], s[i+1]})
	// @..*
}

///
