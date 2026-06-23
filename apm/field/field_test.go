package field

import (
	"slices"
	"testing"
)

func TestLessWithName(t *testing.T) {
	f := Any("", map[string]any{
		"sort": "",
		"page": "",
		"date": "",
		"9":    "",
		"10":   "",
	})
	slices.SortFunc(f.Items, LessWithName(false))
	for _, it := range f.Items {
		println(it.GetName())
	}
}
