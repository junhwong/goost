package runtime

import (
	"fmt"
	"io"
)

type Finalizer interface {
	Finalize() error
}

func Finalize(finalizer interface{}) {
	if finalizer == nil {
		panic("finalizer cannot be nil")
	}
	var err error
	switch ins := finalizer.(type) {
	case func() error:
		err = ins()
	case io.Closer:
		err = ins.Close()
	case Finalizer:
		err = ins.Finalize()
	default:
		//todo:
	}

	if err != nil {
		fmt.Println(err) //log
	}
}
