package validate

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/junhwong/goost/errors"
)

type Validatable interface {
	Validate() error
}
type ValidationError validator.ValidationErrors

func (err ValidationError) Error() string {
	return (validator.ValidationErrors(err)).Error()
}
func (err ValidationError) ResponseStatusCode() int {
	return http.StatusUnprocessableEntity
}
func (err ValidationError) ResponseData() interface{} {
	details := []errors.ErrorDetail{}
	for _, it := range err {

		details = append(details, errors.ErrorDetail{
			Code:  it.Tag(),
			Field: it.Field(),
			// Message: it.Error(),
		})
	}
	if len(details) == 0 {
		details = nil
	}
	return errors.ErrorResult{
		// Code:    code,
		Message: "Validation Failed",
		Details: details,
	}
}

func IsValidationError(err error) bool {
	return AsValidationError(err) != nil
}
func AsValidationError(err error) (target ValidationError) {
	var v validator.ValidationErrors
	if errors.As(err, &v) {
		target = ValidationError(v)
	}
	return
}
