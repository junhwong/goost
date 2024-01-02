package field

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceToGoTimeTempl(t *testing.T) {
	testCases := []struct {
		layout   string
		expected string
	}{
		{
			layout:   "yyyy.MM.dd HH:mm:ss.fffffffff",
			expected: "2006.01.02 15:04:05.999999999",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.layout, func(t *testing.T) {
			v := replaceToGoTimeTempl(tC.layout)
			assert.Equal(t, tC.expected, v)
		})
	}
}
