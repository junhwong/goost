package errors

import (
	"fmt"
	"testing"
)

func TestAs(t *testing.T) {

	// err := &AccessDeniedError{}
	// var err2 *AccessDeniedError // = &AccessDeniedError{}
	// if As(err, &err2) {
	// 	t.Log("ok", err2)
	// }
	//te := new(AccessDeniedError);
}

func TestErrorDetail(t *testing.T) {
	err := &FieldError{
		Code:    "code",
		Field:   "field",
		Message: "message",
		Err:     fmt.Errorf("extra err"),
	}
	t.Log(err.Error())
}

func TestCodeError(t *testing.T) {
	err := New("TEST", "MSG")
	t.Log(err)
	t.Log(GetErrorCode(err))
	t.Log(GetErrorMessage(err))
}
