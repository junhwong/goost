package zap

import (
	"testing"

	"github.com/junhwong/goost/apm"
)

func TestLog(t *testing.T) {

	apm.Default().Debug("hello")
	apm.Default().Debug()
}
