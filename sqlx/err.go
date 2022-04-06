package sqlx

import (
	"errors"
	"fmt"
)

var (
	ErrRecordNotFound = wrapErr(errors.New("record not found"))
)

type SQLError struct {
	Err error
}

func (err *SQLError) Unwrap() error {
	return err.Err
}

func (err *SQLError) Error() string {
	return fmt.Sprintf("stmt: SQLError: %v", err.Err)
}

func wrapErr(err error) *SQLError {
	if err == nil {
		return nil
	}
	return &SQLError{Err: err}
}
