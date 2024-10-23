package time

import (
	"fmt"
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

func TestParseMoment(t *testing.T) {
	// 13天22小时
	// last year
	// last month
	// last week
	// 3 years ago
	// 2 months ago
	// 5 days ago
	// 16 hours ago
	// yesterday
	// on Mar 20
	// on Oct 6, 2023
	testCases := []struct {
		desc string
	}{
		{
			desc: "刚刚",
		},
		{
			desc: "last week",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			v, err := ParseMomentDuration(tC.desc)
			assert.NoError(t, err)
			fmt.Printf("v: %v\n", v)
		})
	}
}
