package helpers

import (
	"fmt"
	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/go-playground/validator/v10"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"log"
	RequestPersistence "main.go/domains/request/presistence"
	ShelterConst "main.go/domains/shelter/presistence"
	"main.go/domains/user/presistence"
	"net"
	"net/smtp"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var (
	once     sync.Once
	verifier *emailverifier.Verifier
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
		"is-vaccinated":       "The %s field must be Vaccinated Or Not Vaccinated !",
		"request-type":        "The %s field must be [adoption, donation, publish, rescue], but got '%s'",
		"payment_type":        "The %s field must be [gopay, shopeepay, bank_transfer,qris]",
		"bank_type":           "The %s field must be [bni, bca, bri, mandiri, cimb, maybank, mega, permata]",
		//"valid-email":         "The %s field must be a valid email address, %s",
		//"valid-domain":        "The %s field must be a valid domain",
	}
	validate *validator.Validate
)

func NewValidator() *validator.Validate {
	//if validate == nil {
	//	validate = validator.New()
	//}
	//if verifier == nil {
	//	verifier = emailverifier.NewVerifier()
	//}
	once.Do(func() {
		validate = validator.New()
		verifier = emailverifier.NewVerifier()
		// This block will only be executed once, regardless of how many times NewValidator is called
	})
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
	if err := validate.RegisterValidation("is-vaccinated", petVaccinated); err != nil {
		panic(err)
	}
	if err := validate.RegisterValidation("request-type", requestTypeValidation); err != nil {
		panic(err)
	}
	if err := validate.RegisterValidation("payment_type", paymentTypeValidation); err != nil {
		panic(err)
	}
	if err := validate.RegisterValidation("bank_type", bankTypeValidation); err != nil {
		panic(err)
	}

	//if err := validate.RegisterValidation("valid-email", checkEmailReachable); err != nil {
	//	panic(err)
	//}
	//if err := validate.RegisterValidation("valid-domain", checkDomainValid); err != nil {
	//	panic(err)
	//}
	return validate
}

func bankTypeValidation(fl validator.FieldLevel) bool {
	reqType := midtrans.Bank(strings.ToLower(fl.Field().String()))
	switch reqType {
	case midtrans.BankBni, midtrans.BankBca, midtrans.BankBri, midtrans.BankMandiri, midtrans.BankCimb, midtrans.BankMaybank, midtrans.BankMega, midtrans.BankPermata:
		return true
	default:
		return false
	}
}

func paymentTypeValidation(fl validator.FieldLevel) bool {
	reqType := coreapi.CoreapiPaymentType(strings.ToLower(fl.Field().String()))
	switch reqType {
	case coreapi.PaymentTypeBankTransfer, coreapi.PaymentTypeGopay, coreapi.PaymentTypeShopeepay, coreapi.PaymentTypeQris:
		return true
	default:
		return false
	}
}

func requestTypeValidation(fl validator.FieldLevel) bool {
	reqType := RequestPersistence.Type(strings.ToLower(fl.Field().String()))
	switch reqType {
	case RequestPersistence.Adoption, RequestPersistence.Donation, RequestPersistence.Rescue, RequestPersistence.Monitoring, RequestPersistence.Publish:
		return true
	default:
		return false
	}
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

func petVaccinated(fl validator.FieldLevel) bool {
	data := fl.Field().String()
	return data == "Vaccinated" || data == "Not Vaccinated"
}

func petAgeValidation(fl validator.FieldLevel) bool {
	age := fl.Field().Int()
	return age >= 0
}

func petGenderValidation(fl validator.FieldLevel) bool {
	gender := fl.Field().String()
	gender = strings.ToLower(gender)
	return gender == ShelterConst.PetGenderMale || gender == ShelterConst.PetGenderFemale || gender == ShelterConst.PetGenderUnknown
}

func roleValidation(fl validator.FieldLevel) bool {
	role := fl.Field().String()
	role = strings.ToLower(role)
	return role == presistence.StaffRole || role == presistence.UserRole || role == ""
}

func isAlphanumericAndSymbol(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	hasAlphaNumeric := regexp.MustCompile(`[a-zA-Z0-9]`).MatchString(field)
	hasSymbol := regexp.MustCompile(`[^a-zA-Z0-9\s]`).MatchString(field) // \s allows spaces; remove \s if spaces should count as symbols
	return hasAlphaNumeric && hasSymbol
}

func checkDomainValid(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	domain := parts[1]
	mxRecords, err := net.LookupMX(domain)
	if err != nil || len(mxRecords) == 0 {
		return false
	}
	return true
}

func checkEmailReachable(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	domain := parts[1]

	mxRecords, err := net.LookupMX(domain)
	if err != nil || len(mxRecords) == 0 {
		return false
	}

	client, err := smtp.Dial(mxRecords[0].Host + ":25")
	if err != nil {
		return false
	}
	defer func(client *smtp.Client) {
		err := client.Close()
		if err != nil {
			log.Fatalf("error closing")
		}
	}(client)

	err = client.Hello("localhost")
	if err != nil {
		return false
	}
	err = client.Mail("test@example.com")
	if err != nil {
		return false
	}
	err = client.Rcpt(email)
	if err != nil {
		return false
	}

	return true
}

func CustomError(err error) (errsMsg []string) {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return []string{"Unexpected error type"}
	}
	for _, e := range validationErrors {
		if errMsg, ok := errorMsg[e.Tag()]; ok {
			if e.Tag() == "request-type" {
				errsMsg = append(errsMsg, fmt.Sprintf(errMsg, e.Field(), e.Value()))
			} else if e.Param() != "" {
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
