package helpers

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

var (
	validate *validator.Validate
)

func NewValidator() *validator.Validate {
	if validate == nil {
		validate = validator.New()
	}
	validate.RegisterValidation("alphanum_symbol", isAlphanumericAndSymbol)
	return validate
}

func ValidateStruct(s interface{}) error {
	return NewValidator().Struct(s)
}

func RegisterValidation(tag string, fn validator.Func, callValidationEvenIfNull ...bool) error {
	return NewValidator().RegisterValidation(tag, fn, callValidationEvenIfNull...)
}

func isAlphanumericAndSymbol(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	hasAlphaNumeric := regexp.MustCompile(`[a-zA-Z0-9]`).MatchString(field)
	hasSymbol := regexp.MustCompile(`[^a-zA-Z0-9\s]`).MatchString(field) // \s allows spaces; remove \s if spaces should count as symbols

	return hasAlphaNumeric && hasSymbol
}
