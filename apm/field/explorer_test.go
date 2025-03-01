package field

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/junhwong/goost/jsonpath"
	"github.com/stretchr/testify/assert"
)

// $.store.book[*].author 不准确

func TestRxplorerRead(t *testing.T) {
	var obj any
	err := json.Unmarshal([]byte(`{ "store": {
    "book": [ 
      { "category": "reference",
        "author": "Nigel Rees",
        "title": "Sayings of the Century",
        "price": 8.95
      },
      { "category": "fiction",
        "author": "Evelyn Waugh",
        "title": "Sword of Honour",
        "price": 12.99
      },
      { "category": "fiction",
        "author": "Herman Melville",
        "title": "Moby Dick",
        "isbn": "0-553-21311-3",
        "price": 8.99
      },
      { "category": "fiction",
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
}`), &obj)
	if !assert.NoError(t, err) {
		return
	}
	root := Any("", obj)
	t.Logf("root: %#v\n", root)

	testCases := []struct {
		desc     string
		err      bool
		vadidate func(t *testing.T, fs []*Field)
	}{
		{
			desc: "$.store.book.*.author",
			vadidate: func(t *testing.T, fs []*Field) {
				assert.Equal(t, 4, len(fs))
			},
		},
		{
			desc: "$.store.book.*.author[0:1]",
			vadidate: func(t *testing.T, fs []*Field) {
				if !assert.Equal(t, 1, len(fs)) {
					return
				}
			},
		},
		{
			desc: "$.store.book.*.author[:].len()",
			vadidate: func(t *testing.T, fs []*Field) {
				if !assert.Equal(t, 1, len(fs)) {
					return
				}
				fmt.Printf("fs[0]: %#v\n", fs[0])
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			exp, _, err := jsonpath.Parse(tC.desc)
			if !assert.NoError(t, err) {
				return
			}
			v := &explorer{readonly: true, root: root, current: []*Field{root}, parent: []*Field{}, getCall: getCall}
			v.visit = func(e jsonpath.Expr) {
				jsonpath.Visit(e, v, v.setError)
			}
			v.visit(exp)
			if tC.err {
				assert.Error(t, err)
				return
			}
			tC.vadidate(t, v.current)
		})
	}
}
