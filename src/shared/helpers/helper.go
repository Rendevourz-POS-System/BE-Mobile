package helpers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/matthewhartstonge/argon2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"main.go/configs/app"
	_const "main.go/configs/const"
	ShelterPresistence "main.go/domains/shelter/presistence"
	"main.go/domains/user/presistence"
	"math/rand"
	"regexp"
	"strconv"
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

func HashPassword(password string) string {
	if argon == nil {
		config := argon2.DefaultConfig()
		argon = &config
	}
	hashedPassword, err := argon.HashEncoded([]byte(password))
	if err != nil {
		return ""
	}
	return string(hashedPassword)
}

func ComparePassword(hashedPassword, password string) bool {
	ok, err := argon2.VerifyEncoded([]byte(password), []byte(hashedPassword))
	if err != nil {
		return false
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

func ToString(value interface{}) string {
	return value.(string)
}

func CheckStaffStatus(value string) bool {
	if value == presistence.StaffRole {
		return true
	}
	return false
}

func GetRole(value string) string {
	if value == presistence.StaffRole {
		return presistence.StaffRole
	}
	return presistence.UserRole
}

func ParseStringToInt(value string) int {
	result, _ := strconv.Atoi(value)
	return result
}

func GetUserId(c *gin.Context) primitive.ObjectID {
	userId, _ := c.MustGet("x-user-id").(string)
	userID, err := primitive.ObjectIDFromHex(userId)
	fmt.Println(userID)
	if err != nil {
		panic("Failed to get user id from middlewares !")
	}
	return userID
}

func ParseStringToObjectId(value string) primitive.ObjectID {
	objectId, err := primitive.ObjectIDFromHex(value)
	if err != nil {
		panic("Failed to parse string to object id !")
	}
	return objectId
}

func ParseObjectIdToString(value primitive.ObjectID) string {
	return value.Hex()
}

func CheckPetGender(value string) string {
	if value == ShelterPresistence.PetGenderMale || value == ShelterPresistence.PetGenderFemale {
		return value
	}
	return ""
}

func GenerateFileName(filename string) string {
	return GenerateRandomString(10) + "_" + filename
}

func RegexPattern(pattern interface{}) *bson.M {
	return &bson.M{"$regex": primitive.Regex{
		Pattern: "^" + regexp.QuoteMeta(ToString(pattern)) + "$", // Exact match, case insensitive
		Options: "i",                                             // Case-insensitive
	}}
}
