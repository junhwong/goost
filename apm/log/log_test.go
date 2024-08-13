package log

import (
	"testing"

	"github.com/junhwong/goost/apm"
)

func TestLog(t *testing.T) {
	t.Cleanup(apm.Flush)

	Debug("aaa")
}
