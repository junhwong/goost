package field

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/junhwong/goost/jsonpath"
)

func TestParseTime(t *testing.T) {
	// 1680452604227276000
	// 1680452349917244000
	//2023.04.02 16:19:09.917244

	x, err := time.Parse("2006.01.02 15:04:05.999999999", "2023.04.02 16:19:09.917244")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("x: %v\n", x)
	fmt.Printf("x.UnixNano(): %v\n", x.UnixNano())

	var i = int64(1680452349917244000)

	y := time.Unix(0, i)
	fmt.Printf("y: %v\n", y)
	fmt.Printf("y.UnixNano(): %v\n", y.UnixNano())

	// 2023.04.02 16:11:15.847959
	z, err := time.Parse(time.RFC3339Nano, "2023-04-03T00:11:15.847959+09:00")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(z.Zone())
	fmt.Printf("z: %v\n", z.UTC())
	fmt.Printf("z.UnixNano(): %v\n", z.UnixNano())
}

func TestXXX(t *testing.T) {

	s := `{
		"store": {
			"book": [{
					"category": "reference",
					"author": "Nigel Rees",
					"title": "Sayings of the Century",
					"price": 8.95
				}, {
					"category": "fiction",
					"author": "Evelyn Waugh",
					"title": "Sword of Honour",
					"price": 12.99
				}, {
					"category": "fiction",
					"author": "Herman Melville",
					"title": "Moby Dick",
					"isbn": "0-553-21311-3",
					"price": 8.99
				}, {
					"category": "fiction",
					"author": "J. R. R. Tolkien",
					"title": "The Lord of the Rings",
					"isbn": "0-395-19395-8",
					"price": 22.99
				}
			],
			"bicycle": {
				"color": "red",
				"price": 19.95
			}
		}
	}`

	var v any
	err := json.Unmarshal([]byte(s), &v)
	if err != nil {
		t.Fatal(err)
	}

	// fmt.Printf("v: %v                                                                                                   \n", v)
	root := MakeRoot()
	f := Any("x", v)
	root.Set(f)

	// fmt.Printf("f: %#v\n", f)

	{
		r, err := Find(root, "x.store.book[0].title")
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("r: %#v\n", r[0])
		return
	}
	// {
	// 	expr, _ := jsonpath.Parse("x2.xxxxxxxx1")
	// 	err := Apply(expr, root, func(f *Field) {
	// 		f.SetString("hello")
	// 	})
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	fmt.Printf("root: %#v\n", root)
	// }

	{
		expr, _ := jsonpath.Parse("x3.[]")
		err := Apply(expr, root, func(f *Field) {
			f.SetString("hello")
		})
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("root: %#v\n", root)
	}
}
