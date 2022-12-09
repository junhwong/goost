package zap

import (
	"testing"

	"github.com/junhwong/goost/apm"
	"github.com/junhwong/goost/apm/deflog"
)

func TestLog(t *testing.T) {

	provider := &impl{}
	provider.AddHandlers(&deflog.ConsoleHandler{Formatter: deflog.NewTextFormatter()})

	apm.SetDefault(provider)

	apm.Default().Debug("hello")
	apm.Default().Debug()
}
