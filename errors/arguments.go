package errors

import (
	"net/http"
	"strings"
)

// 表示错误明细
type FieldError struct {
	Err     error  `json:"-" xml:"-"`
	Code    string `json:"code,omitempty" xml:"code,omitempty"`
	Field   string `json:"field,omitempty" xml:"field,omitempty"`
	Message string `json:"message,omitempty" xml:"message,omitempty"`
}

func (err *FieldError) Unwrap() error {
	return err.Err
}
func (err *FieldError) Error() string {
	s := ""

	if err.Code != "" {
		s += err.Code

	}
	msg := err.Message
	if err.Err != nil {
		if msg != "" {
			msg += ". "
		}
		msg += err.Err.Error()
	}
	if msg != "" {
		if s != "" && err.Code != "" {
			s += ", "
		}
		s += msg
	}
	if err.Field != "" {
		if s != "" {
			s = err.Field + ": " + s
		} else {
			s += err.Field
		}
	}
	if s == "" {
		return "field_error"
	}
	return s
}

func (err *FieldError) ResponseStatusCode() int {
	return http.StatusUnprocessableEntity
}
func (err *FieldError) ResponseData() interface{} {
	code := err.Code
	if code == "" {
		code = "invalid_argument"
	}
	return ErrorDetail{
		Code:    code,
		Field:   err.Field,
		Message: err.Message,
	}
}
func (err *FieldError) ResponseDetailData() ErrorDetail {
	data, _ := (err.ResponseData()).(ErrorDetail)
	return data
}

// 参数错误
type ArgumentsError struct {
	Err     error
	Code    string // missing missing_field invalid already_exists unprocessable custom
	Message string
	Details []error
}

func (err *ArgumentsError) Unwrap() error {
	return err.Err
}
func (err *ArgumentsError) Error() string {
	s := "ArgumentError: "

	arr := []string{}
	for _, e := range err.Details {
		arr = append(arr, e.Error())
	}

	return s + strings.Join(arr, ",")
}
func (err *ArgumentsError) ResponseStatusCode() int {
	return http.StatusUnprocessableEntity
}
func (err *ArgumentsError) ResponseData() interface{} {
	code := err.Code
	if code == "" {
		code = "invalid_arguments"
	}
	details := []ErrorDetail{}
	for _, e := range err.Details {
		if rd, ok := e.(ErrorResponseDetail); ok {
			details = append(details, rd.ResponseDetailData())
		}
	}
	if len(details) == 0 {
		details = nil
	}
	return ErrorResult{
		Code:    code,
		Message: err.Message,
		Details: details,
	}
}

func AsArgumentsError(err error) (target *ArgumentsError) {
	if !As(err, &target) {
		return nil
	}
	return
}
