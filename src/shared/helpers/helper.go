package helpers

import (
	"fmt"
	Request "main.go/src/domains/request/entities"
	Pet "main.go/src/domains/shelter/entities"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/matthewhartstonge/argon2"
	"github.com/nanorand/nanorand"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"main.go/src/configs/app"
	_const "main.go/src/configs/const"
	ShelterPresistence "main.go/src/domains/shelter/presistence"
	"main.go/src/domains/user/presistence"
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

func GetVerifiedUrl(secretCode string, Otp *int) string {
	return app.GetConfig().Domain.Protocol + "://exp" + app.GetConfig().Domain.Name + ":" + app.GetConfig().Domain.Port + "/--" + app.GetConfig().Domain.FrontendPath + "/" + secretCode
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

func GenerateOTP(length int) *int {
	code, err := nanorand.Gen(length)
	if err != nil {
		return nil
	}
	codes, _ := strconv.Atoi(code)
	return &codes
}

func ToString(value interface{}) string {
	return fmt.Sprintf("%v", value)
}

func PartIntToString(value int) string {
	return strconv.Itoa(value)
}

func ParsePointerIntToString(value *int) string {
	return strconv.Itoa(*value)
}

func CheckStaffStatus(value string) bool {
	if presistence.Role(value) == presistence.StaffRole {
		return true
	}
	return false
}

func GetRole(value string) string {
	if presistence.Role(value) == presistence.StaffRole {
		return ToString(presistence.StaffRole)
	}
	return ToString(presistence.UserRole)
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

func GetRoleFromContext(c *gin.Context) string {
	userRole, _ := c.MustGet("x-user-role").(string)
	return userRole
}

func ParseStringToObjectId(value string) primitive.ObjectID {
	objectId, err := primitive.ObjectIDFromHex(value)
	if err != nil {
		logrus.Warnf("Failed to parse string to object id !")
	}
	return objectId
}

func ParseStringToObjectIdAddress(value string) *primitive.ObjectID {
	if value == "" {
		return nil
	}
	objectId, err := primitive.ObjectIDFromHex(value)
	if err != nil {
		logrus.Warnf("Failed to parse string to object id !")
	}
	return &objectId
}

func GetAddressString(value string) *string {
	return &value
}

func ArrayAddress(value []string) *[]string {
	return &value
}

func ParseStringToBoolean(value string) *bool {
	if value == "" {
		return nil
	}
	values, err := strconv.ParseBool(value)
	if err != nil {
		logrus.Warnf("Failed to parse string to boolean !")
	}
	return &values
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

func FillRequestData(req *Pet.Pet, ctx *gin.Context) (res *Request.Request) {
	userId := GetUserId(ctx)
	res = &Request.Request{
		PetId: &req.ID,
		//ShelterId: *req.ShelterId,
		UserId: userId,
	}
	return res
}
