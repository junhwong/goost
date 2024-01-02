package jsonpath

import (
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
