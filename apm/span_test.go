package apm

import (
	"context"
	"testing"
)

func TestSpanCaller(t *testing.T) {
	t.Cleanup(Flush)

	_, span := Default().NewSpan(context.TODO())
	// _, span = Start(context.TODO())
	_, span = Default().WithFields(LogComponent("t")).NewSpan(context.TODO())
	span.End()
}
