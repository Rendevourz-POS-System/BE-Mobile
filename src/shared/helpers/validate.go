package helpers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	ShelterConst "main.go/domains/shelter/presistence"
	"main.go/domains/user/presistence"
	"reflect"
	"regexp"
	"strconv"
)

var (
	errorMsg = map[interface{}]string{
		"required":            "The %s field is required",
		"email":               "The %s field must be a valid email address",
		"min":                 "The %s field must be at least %s characters",
		"alphanum_symbol":     "The %s field must contain at least one letter, one number, and one symbol",
		"number":              "The %s field must be a number",
		"max":                 "The %s field must be at most %s characters",
		"role":                "The %s field must be a valid be either Staff or User",
		"pet-gender":          "The %s field must be a valid be either Male or Female",
		"pet-age":             "The %s field must be a valid number and greater than or equal to 0",
		"pet-accepted-min":    "The %s field must be a valid and greater than or equal to 0",
		"min-location-length": "The %s field must be at least %s characters",
	}
	validate *validator.Validate
)

func NewValidator() *validator.Validate {
	if validate == nil {
		validate = validator.New()
	}
	if err := validate.RegisterValidation("alphanum_symbol", isAlphanumericAndSymbol); err != nil {
		panic(err)
	}
	if err := validate.RegisterValidation("role", roleValidation); err != nil {
		panic(err)
	}
	if err := validate.RegisterValidation("pet-gender", petGenderValidation); err != nil {
		panic(err)
	}
	if err := validate.RegisterValidation("pet-age", petAgeValidation); err != nil {
		panic(err)
	}
	if err := validate.RegisterValidation("pet-accepted-min", petTypeAcceptedMin); err != nil {
		panic(err)
	}
	if err := validate.RegisterValidation("min-location-length", shelterLocationMinLength); err != nil {
		panic(err)
	}
	return validate
}

func shelterLocationMinLength(fl validator.FieldLevel) bool {
	param := fl.Param() // Get parameter from the tag, e.g., "3"
	if minLen, err := strconv.Atoi(param); err == nil {
		field := fl.Field()
		if field.Kind() == reflect.String {
			return len(field.String()) >= minLen
		} else if field.Kind() == reflect.Slice || field.Kind() == reflect.Array {
			return field.Len() >= minLen
		}
	}
	return false
}

func petTypeAcceptedMin(fl validator.FieldLevel) bool {
	field := fl.Field()
	// Check if the field is actually a slice or array
	if field.Kind() == reflect.Slice || field.Kind() == reflect.Array {
		length := field.Len()
		return length > 0
	}
	// If it's not a slice or array, the validation does not apply
	return false
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
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return []string{"Unexpected error type"}
	}
	for _, e := range validationErrors {
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
	return errsMsg
}
