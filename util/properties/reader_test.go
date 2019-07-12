package properties_test

import (
	"io"
	"strings"
	"testing"

	"github.com/junhwong/goboot/util/properties"
)

func TestReaderWithCorrect(t *testing.T) {
	text := `# doc comment,
	# doc comment.
	
	# comment line
	b:true
	num=1.2
	base64_val = base64:dGhpcyUyMGlzJTIwYmFzZTY0 # the value with base64
	Chinese=你好世界!
	en-str="hello \"word\"!"
	en-str2="hello \
word2!"
	$key=hel\
lo
	$key2=hel\lo2
	empty=
	empty= #empty value
	path.file = file:///var/tmps`

	rd := properties.NewReader(strings.NewReader(text))
	i := 0
	for {
		i++
		k, v, err := rd.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Fatal(i, ": ", err)
		}
		t.Log(i, ": ", string(k), "=", string(v))

	}
}

func TestReaderWithIncorrect(t *testing.T) {
	testCases := []struct {
		desc string
		s    string
	}{
		{
			desc: "m-comment",
			s:    "/* ...",
		},
		{
			desc: "s-comment",
			s:    "//...",
		},
		{
			desc: "empty key",
			s:    "=v",
		},
		{
			desc: "bad key1",
			s:    "@key=v",
		},
		{
			desc: "bad key2",
			s:    "k ey=v",
		},
		{
			desc: "bad key3",
			s:    "1key=v",
		},
		{
			desc: "bad key4",
			s:    "-key=v",
		},
		{
			desc: "bad key4",
			s:    `k\tey=v`,
		},
		{
			desc: "bad val1",
			s:    "key=v v",
		},
		{
			desc: "bad val2",
			s:    "key=\"val",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			rd := properties.NewReader(strings.NewReader(tC.s))
			for {
				k, v, err := rd.Next()
				if err == io.EOF {
					break
				} else if err == nil {
					t.Fatal(string(k), "=", string(v))
				}
			}
		})
	}
}
