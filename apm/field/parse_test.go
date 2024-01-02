package field

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

func TestSplitPath(t *testing.T) {
	testCases := []struct {
		desc string
		err  bool
	}{
		{
			desc: "a",
		},
		{
			desc: "a.b",
		},
		{
			desc: "a..b",
			err:  true,
		},
		{
			desc: "c#1",
		},
		{
			desc: "c#1.2",
		},
		{
			desc: "c#1.2#3",
		},
		{
			desc: "d[e]",
		},
		{
			desc: `"f"`,
		},
		{
			desc: `"f".g`,
		},
		{
			desc: `a"f".g`,
			err:  true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			r, err := SplitPath(tC.desc)
			if tC.err {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
			fmt.Printf("r: %v\n", r)
		})
	}
}
