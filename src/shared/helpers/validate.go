package helpers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	ShelterConst "main.go/domains/shelter/presistence"
	"main.go/domains/user/presistence"
	"regexp"
)

var (
	errorMsg = map[interface{}]string{
		"required":        "The %s field is required",
		"email":           "The %s field must be a valid email address",
		"min":             "The %s field must be at least %s characters",
		"alphanum_symbol": "The %s field must contain at least one letter, one number, and one symbol",
		"number":          "The %s field must be a number",
		"max":             "The %s field must be at most %s characters",
		"role":            "The %s field must be a valid be either Staff or User",
		"pet-gender":      "The %s field must be a valid be either Male or Female",
		"pet-age":         "The %s field must be a valid number and greater than or equal to 0",
	}
	validate *validator.Validate
)

func NewValidator() *validator.Validate {
	if validate == nil {
		validate = validator.New()
	}
	err := validate.RegisterValidation("alphanum_symbol", isAlphanumericAndSymbol)
	if err != nil {
		panic(err)
	}
	err = validate.RegisterValidation("role", roleValidation)
	err = validate.RegisterValidation("pet-gender", petGenderValidation)
	err = validate.RegisterValidation("pet-age", petAgeValidation)
	return validate
}

func petAgeValidation(fl validator.FieldLevel) bool {
	age := fl.Field().Int()
	return age >= 0
}

func petGenderValidation(fl validator.FieldLevel) bool {
	gender := fl.Field().String()
	return gender == ShelterConst.PetGenderMale || gender == ShelterConst.PetGenderFemale || gender == ShelterConst.PetGenderUnknown
}

func roleValidation(fl validator.FieldLevel) bool {
	role := fl.Field().String()
	return role == presistence.StaffRole || role == presistence.UserRole || role == ""
}

func isAlphanumericAndSymbol(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	hasAlphaNumeric := regexp.MustCompile(`[a-zA-Z0-9]`).MatchString(field)
	hasSymbol := regexp.MustCompile(`[^a-zA-Z0-9\s]`).MatchString(field) // \s allows spaces; remove \s if spaces should count as symbols
	return hasAlphaNumeric && hasSymbol
}

func CustomError(err error) (errsMsg []string) {
	for _, e := range err.(validator.ValidationErrors) {
		if errMsg, ok := errorMsg[e.Tag()]; ok {
			if e.Param() != "" {
				errsMsg = append(errsMsg, fmt.Sprintf(errMsg, e.Field(), e.Param()))
			} else {
				errsMsg = append(errsMsg, fmt.Sprintf(errMsg, e.Field()))
			}
		} else {
			errsMsg = append(errsMsg, fmt.Sprintf("The %s field is invalid", e.Field()))
		}
	}
	//for _, err := range err.(validator.ValidationErrors) {
	//	fmt.Println(err.Namespace(), err.Field(), err.StructNamespace(), err.StructField(), err.Tag(), err.ActualTag(), err.Kind(), err.Type(), err.Value(), err.Param())
	//}
	return errsMsg
}
