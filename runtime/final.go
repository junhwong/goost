package runtime

import (
	"fmt"
	"io"
)

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
	default:
		//todo:
	}

	if err != nil {
		fmt.Println(err) //log
	}
}
