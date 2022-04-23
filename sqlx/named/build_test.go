package named

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildSQL(t *testing.T) {
	testCases := []struct {
		tpl      string
		expected string
		c        int
	}{
		{
			tpl:      ":p",
			expected: "?",
			c:        1,
		},
		{
			tpl:      "{:p}",
			expected: "?",
			c:        1,
		},
		{
			tpl:      "{:p=abc}",
			expected: "?",
			c:        1,
		},
		{
			tpl:      "{:p=123}",
			expected: "?",
			c:        1,
		},
		{
			tpl:      "{:p=abc123}",
			expected: "?",
			c:        1,
		},
		{
			tpl:      "{:p=abc123_}",
			expected: "?",
			c:        1,
		},
		{
			tpl:      "{:p|pipe|pipe2}",
			expected: "?",
			c:        1,
		},
		{
			tpl:      "{:p=val|pipe|pipe2}",
			expected: "?",
			c:        1,
		},
		{
			tpl:      "{@time_now}",
			expected: "?",
			c:        1,
		},
		{
			tpl:      "{@time_now|pipe|pipe2}",
			expected: "?",
			c:        1,
		},
	}

	pipes["pipe"] = func(in interface{}) (out interface{}, err error) {
		return
	}
	pipes["pipe2"] = func(in interface{}) (out interface{}, err error) {
		return
	}

	for _, tC := range testCases {
		t.Run(tC.tpl, func(t *testing.T) {
			sql, p, err := BuildNamedQuery(tC.tpl)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tC.expected, sql)
			assert.Equal(t, tC.c, len(p))
		})
	}

}

func TestR(t *testing.T) {
	reg := regexp.MustCompile(`\:\w+[^\W]|\{(\:\w+(\=\w+)?)|(\@\w+)(\|[a-zA-Z_]\w+)*\}`)

	t.Log(reg.FindAllString("{:name=simple_default_value|pipe|pipe2}", -1))
	t.Log(reg.FindAllString(":name=simple_default_value|pipe|pipe2", -1))
}
func TestR2(t *testing.T) {
	reg := regexp.MustCompile(`\:\w+[^\W]`)

	t.Log(reg.FindAllString("p=:name", -1))
	t.Log(reg.FindAllString("=:name=simple_default_value|pipe|pipe2", -1))
}
func TestR3(t *testing.T) {
	reg := parameterExtract

	t.Log(reg.FindAllString("p=:p c", -1))
	t.Log(reg.FindAllString("{:name=simple_default_value}", -1))
	t.Log(reg.FindAllString("{:name=simple_default_value|pipe|pipe2}", -1))
	t.Log(reg.FindAllString("{@name=simple_default_value|pipe|pipe2}", -1))
	t.Log(reg.FindAllString("{@name|pipe|pipe2} end", -1))
}
