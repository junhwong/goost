package apm

import (
	"fmt"
	"runtime/debug"
	"testing"
)

func TestStack(t *testing.T) {
	s := debug.Stack()
	fmt.Printf("s: %s\n", s)
	fmt.Println("======")
	ci := StackToCallerInfo(s)
	fmt.Printf("ci: %+v\n", ci)
}
