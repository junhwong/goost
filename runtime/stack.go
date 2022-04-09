package runtime

import (
	"errors"
	"runtime"
)

type wrappedStackError struct {
	Err   error
	Stack []byte

	depth int
}

func (err *wrappedStackError) GetCallStack() []byte {
	return err.Stack
}

func (err *wrappedStackError) Unwrap() error {
	return err.Err
}
func (err *wrappedStackError) Error() string {
	return err.Err.Error()
}

func WrapCallStack(err error, depth int, forceWrap ...bool) (ex *wrappedStackError) {
	if err == nil {
		panic("err cannot nil")
	}
	f := false
	if i := len(forceWrap); i > 0 {
		f = forceWrap[i-1]
	}
	if errors.As(err, &ex) && !f {
		return
	}
	ex = &wrappedStackError{
		Err:   err,
		depth: depth + 1,
	}
	const size = 64 << 10
	stacktrace := make([]byte, size)
	ex.Stack = stacktrace[:runtime.Stack(stacktrace, false)]
	return
}
