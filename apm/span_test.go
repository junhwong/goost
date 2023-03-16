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

func TestHexID(t *testing.T) {
	id := NewHexID()
	if len(id.String()) != 32 {
		t.Fatal()
	}
	id2, err := ParseHexID(id.String())
	if err != nil {
		t.Fatal(err)
	}
	if id.High != id2.High || id.Low != id2.Low {
		t.Fatal()
	}
	id.High = 0
	if len(id.String()) != 16 {
		t.Fatal()
	}
	id2, err = ParseHexID(id.String())
	if err != nil {
		t.Fatal(err)
	}
	if id.High != id2.High || id.Low != id2.Low {
		t.Fatal()
	}
	id.Low = 0
	if len(id.String()) != 0 {
		t.Fatal()
	}
}
