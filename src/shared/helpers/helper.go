package helpers

import (
	"github.com/matthewhartstonge/argon2"
	_const "main.go/configs/const"
)

var (
	argon *argon2.Config
)

func ParseDatabase(database string) string {
	if database == _const.DB_SHELTER_APP {
		return _const.DB_SHELTER_APP
	}
	return _const.DB_SHELTER_APP
}

func HashPassword(password string) (string, error) {
	if argon == nil {
		config := argon2.DefaultConfig()
		argon = &config
	}
	hashedPassword, err := argon.HashEncoded([]byte(password))
	if err != nil {
		return "", err // Return the error to be handled by the caller
	}
	return string(hashedPassword), nil
}

func ComparePassword(hashedPassword, password string) bool {
	if argon == nil {
		*argon = argon2.DefaultConfig()
	}
	ok, err := argon2.VerifyEncoded([]byte(password), []byte(hashedPassword))
	if err != nil {
		panic(err) // 💥
	}
	return ok
}
