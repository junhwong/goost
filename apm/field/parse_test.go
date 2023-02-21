package field

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLabel(t *testing.T) {
	testCases := []struct {
		desc      string
		expectErr bool
		expectArr []string
	}{
		{
			desc: "",
		},
		{
			desc: ",",
		},
		{
			desc:      "a;b:c d;e",
			expectArr: []string{"a", "b", "c", "d", "e"},
		},
		{
			desc:      "a,",
			expectArr: []string{"a"},
		},
		{
			desc:      "%-1.2",
			expectArr: []string{"%-1.2"},
		},
		{
			desc:      "%",
			expectErr: true,
		},
		{
			desc:      "中文",
			expectErr: true,
		},
		{
			desc:      "#",
			expectErr: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			r, err := ParseLabelValue(tC.desc)
			if err != nil {
				if !tC.expectErr {
					t.Fatal(err)
				}
				return
			}
			assert.Equal(t, tC.expectArr, r)
		})
	}
}
