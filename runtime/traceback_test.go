package runtime

import (
	"strings"
	"testing"
)

func TestCaller(t *testing.T) {
	v := Caller(0)
	if v.Method != "runtime.TestCaller" {
		t.Fail()
	}
	if v.File != "traceback_test.go" {
		t.Fail()
	}
	if !(strings.Contains(v.Path, "/runtime") && strings.Contains(v.Package, "/goost")) {
		t.Fail()
	}
}
