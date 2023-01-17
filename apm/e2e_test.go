package apm_test

import (
	"errors"
	"testing"

	"github.com/junhwong/goost/apm"
)

func TestLog(t *testing.T) {
	t.Cleanup(apm.Flush)

	// apm.UseAsyncDispatcher()

	apm.Default().WithFields(apm.LogComponent("test")).Debug("hello")
	apm.Default().Debug(apm.WrapCallStack(errors.New("hhh")))
}
