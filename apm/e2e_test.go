package apm_test

import (
	"errors"
	"testing"

	"github.com/junhwong/goost/apm"
)

func TestLog(t *testing.T) {
	t.Cleanup(apm.Flush)

	// apm.UseAsyncDispatcher()
	apm.LogComponent("test")
	apm.Default().WithFields().Debug("hello")
	apm.Default().Debug(apm.WrapCallStack(errors.New("hhh")))
}
