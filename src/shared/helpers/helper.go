package helpers

import (
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
	"strings"
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
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
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
	value = strings.ToLower(value)
	if value == ShelterPresistence.PetGenderMale || value == ShelterPresistence.PetGenderFemale {
		return value
	}
	return ""
}

func GenerateFileName(filename string) string {
	return GenerateRandomString(10) + "_" + filename
}

func RegexCaseInsensitivePattern(pattern interface{}) *bson.M {
	// Convert the input to a string, escaping any special regex characters to avoid issues in pattern matching
	safePattern := regexp.QuoteMeta(ToString(pattern))
	// Adjust the pattern to match any part of the string (i.e., contains, not just starts with or exact match)
	regexPattern := ".*" + safePattern + ".*"

	return &bson.M{"$regex": primitive.Regex{
		Pattern: regexPattern,
		Options: "i", // Case-insensitive
	}}
}
