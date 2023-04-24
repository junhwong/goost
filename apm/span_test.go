package apm

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestSpanCaller(t *testing.T) {
	t.Cleanup(Flush)

	var sb strings.Builder

	sd := dispatcher.(*syncDispatcher)
	sd.handlers = handlerSlice{&SimpleHandler{
		IsEnd:     true,
		Formatter: &TextFormatter{},
		Out:       &sb,
	}}

	var span Span

	{
		_, span = Default().NewSpan(context.TODO())
		span.End()

	}
	{
		_, span = Start(context.TODO())
		span.End()
	}
	{
		_, span = Default().WithFields(LogComponent("t")).NewSpan(context.TODO())
		span.End()
	}

	if strings.Count(sb.String(), "apm.TestSpanCaller") != 3 {
		t.Fatal(sb.String())
	}
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
	if !bytes.Equal(id.High(), id2.High()) || !bytes.Equal(id.Low(), id2.Low()) {
		t.Fatal()
	}
	id3, err := ParseHexID(id2.Low().String())
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(id.Low(), id3.Low()) {
		t.Fatal()
	}

}

func TestHexIDNil(t *testing.T) {
	var id HexID

	if id.String() != "" {
		t.Fatal()
	}
}

func TestContextKey(t *testing.T) {
	ctx := context.TODO()
	var (
		s1 = ""
		s2 = ""
	)
	c1 := context.WithValue(ctx, &s1, "s1")
	c2 := context.WithValue(ctx, &s2, "s2")

	k := ""

	if c1.Value(k) == "s1" {
		t.Fatal()
	}
	if c1.Value(&k) == c1.Value(&s1) {
		t.Fatal(c1.Value(&k))
	}
	if c1.Value(&s1) == c2.Value(&s1) {
		t.Fatal()
	}
	if c1.Value(&s1) == c2.Value(&s2) {
		t.Fatal()
	}
}
