package field

import "fmt"

type logger interface {
	Error(...any)
}

var log logger = &plog{}

type plog struct {
}

func (p *plog) Error(args ...any) {
	fmt.Println(args...)
}
