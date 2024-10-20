package field

import (
	"fmt"
	"testing"
)

// func TestSetSort(t *testing.T) {
// 	var fs FieldSet
// 	// fs = append(fs, &Field{Key: "a"})
// 	// fs = append(fs, &Field{Key: "e"})
// 	// fs = append(fs, &Field{Key: "c"})
// 	fs.Sort()

// 	var s string
// 	for _, v := range fs {
// 		s += v.Key
// 	}
// 	if s != "ace" {
// 		t.Fatal()
// 	}
// }

// func TestSetRemove(t *testing.T) {
// 	var fs FieldSet
// 	// fs = append(fs, &Field{Key: "a"})
// 	// fs = append(fs, &Field{Key: "e"})
// 	// fs = append(fs, &Field{Key: "c"})
// 	// fs = append(fs, &Field{Key: "e"})

// 	if fs.Remove("e") == nil {
// 		t.Fatal()
// 	}
// 	if x := fs.Remove("e"); x != nil {
// 		t.Fatal()
// 	}
// 	var s string
// 	for _, v := range fs {
// 		s += v.Key
// 	}
// 	if s != "ac" {
// 		t.Fatal()
// 	}
// }

type tmeta struct {
	key  string
	kind string
	its  []*tmeta
}

type tval struct {
	sv  string
	its []*tval
}

type tfld struct {
	*tmeta
	val *tval
	prt *tfld
}

func TestTV(t *testing.T) {

	hy := tfld{
		tmeta: &tmeta{
			key:  "合约",
			kind: "record",
			its: []*tmeta{
				{key: "交易品种"},
				{key: "报价单位"},
			},
		},
		val: &tval{
			its: []*tval{
				{sv: "黄金"},
				{sv: "克"},
			},
		},
	}

	fmt.Printf("hy: %#v\n", hy)

	k := tfld{
		tmeta: hy.tmeta.its[0],
		val:   hy.val.its[0],
		prt:   &hy,
	}

	fmt.Printf("k: %+v\n", k)

}
