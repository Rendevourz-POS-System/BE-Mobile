package helpers

import "github.com/go-playground/validator/v10"

var (
	validate *validator.Validate
)

func NewValidator() *validator.Validate {
	if validate == nil {
		validate = validator.New()
	}
	return validate
}

func ValidateStruct(s interface{}) error {
	return NewValidator().Struct(s)
}
