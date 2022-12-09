package apm_test

import (
	"errors"
	"testing"

	"github.com/junhwong/goost/apm"
	"github.com/junhwong/goost/apm/deflog"
)

func TestLog(t *testing.T) {
	t.Cleanup(apm.Flush)

	apm.AddHandlers(&deflog.ConsoleHandler{
		Formatter: deflog.NewTextFormatter(),
	})
	apm.UseAsyncDispatcher()

	apm.Default().Debug("hello")
	apm.Default().Debug(apm.WrapCallStack(errors.New("hhh")))
}
