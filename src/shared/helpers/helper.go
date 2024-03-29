package helpers

import (
	"github.com/matthewhartstonge/argon2"
	"main.go/configs/app"
	_const "main.go/configs/const"
	"math/rand"
	"time"
)

var (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	argon   *argon2.Config
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

func GetVerifiedUrl(secretCode, email string) string {
	return app.GetConfig().Domain.Protocol + "://" + app.GetConfig().Domain.Name + ":" + app.GetConfig().Domain.Port + app.GetConfig().Domain.FrontendPath + "/" + secretCode
}

func GetCurrentTime(minute *int) *time.Time {
	if minute != nil {
		times := time.Now().Add(time.Minute * time.Duration(*minute))
		return &times
	}
	times := time.Now()
	return &times
}

func GenerateRandomString(length int) string {
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
