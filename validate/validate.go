package validate

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type (
	Validator  = *validator.Validate
	FieldLevel = validator.FieldLevel
)

var (
	validate                    = validator.New()
	RegisterTagNameFunc         = validate.RegisterTagNameFunc
	RegisterValidation          = validate.RegisterValidation
	RegisterValidationCtx       = validate.RegisterValidationCtx
	RegisterAlias               = validate.RegisterAlias
	RegisterStructValidation    = validate.RegisterStructValidation
	RegisterStructValidationCtx = validate.RegisterStructValidationCtx
	RegisterCustomTypeFunc      = validate.RegisterCustomTypeFunc
	RegisterTranslation         = validate.RegisterTranslation
	Struct                      = validate.Struct
	StructCtx                   = validate.StructCtx
	StructFiltered              = validate.StructFiltered
	StructFilteredCtx           = validate.StructFilteredCtx
	StructPartial               = validate.StructPartial
	StructPartialCtx            = validate.StructPartialCtx
	StructExcept                = validate.StructExcept
	StructExceptCtx             = validate.StructExceptCtx
	Var                         = validate.Var
	VarCtx                      = validate.VarCtx
	VarWithValue                = validate.VarWithValue
	VarWithValueCtx             = validate.VarWithValueCtx
)

func init() {
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}
