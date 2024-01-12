package jsonpath

import (
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
			desc: "中文",
			ex:   "f:中文",
		},
		{
			desc: "a",
			ex:   "f:a",
		},
		{
			desc: "a.b",
			ex:   "m:(f:a,f:b)",
		},
		{
			desc: "a..b",
			ex:   "m:(f:a,s:..,f:b)",
		},
		{
			desc: "c#1",
			ex:   "m:(f:c,i:1)",
		},
		{
			desc: "c#1#2",
			ex:   "m:(f:c,i:1,i:2)",
		},
		{
			desc: "c#1.2",
			ex:   "m:(f:c,i:1,f:2)",
		},
		{
			desc: "c#1.2#3",
			ex:   "m:(f:c,i:1,f:2,i:3)",
		},
		{
			desc: "d[1]",
			ex:   "m:(f:d,i:1)",
		},
		{
			desc: `"f"`,
			ex:   `q:"f"`,
		},
		{
			desc: `"f".g`,
			ex:   `m:(q:"f",f:g)`,
		},
		{
			desc: `a"f".g`,
			err:  true,
		},
		{
			desc: `"f"a.g`,
			err:  true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			r, err := Parse(tC.desc)
			if tC.err {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tC.ex, String(r))
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

}
