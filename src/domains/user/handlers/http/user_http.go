package http

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"main.go/configs/app"
	_const "main.go/configs/const"
	"main.go/configs/database"
	User "main.go/domains/user/entities"
	"main.go/domains/user/interfaces"
	"main.go/domains/user/interfaces/impl/repository"
	"main.go/domains/user/interfaces/impl/usecase"
	"main.go/middlewares"
	"main.go/shared/helpers"
	"main.go/shared/helpers/image_helpers"
	"main.go/shared/message/errors"
	"net/http"
)

type UserHttp struct {
	userUsecase   interfaces.UserUsecase
	userTokenHttp *UserTokenHttp
}

func NewUserHttp(router *gin.Engine, tokenHttp *UserTokenHttp) *UserHttp {
	handler := &UserHttp{
		userUsecase: usecase.NewUserUsecase(repository.NewUserRepository(
			database.GetDatabase(_const.DB_SHELTER_APP))),
		userTokenHttp: tokenHttp,
	}
	guest := router.Group("/user")
	{
		guest.GET("/", handler.FindAll)
		guest.POST("/register", handler.RegisterUsers)
		guest.POST("/login", handler.LoginUsers)
		guest.POST("/verify-email", handler.AccountVerification)
		guest.POST("/resend-otp", handler.ResendVerificationOtp)
	}
	user := router.Group(guest.BasePath(), middlewares.JwtAuthMiddleware(app.GetConfig().AccessToken.AccessTokenSecret, "user"))
	{
		user.GET("/data", handler.FindUserByToken)
		user.PUT("/update", handler.UpdateUser)
		user.PUT("/update-pw", handler.UpdatePassword)
		user.DELETE("/delete/account", handler.DeleteUserAccount)
	}
	userAndAdmin := router.Group(guest.BasePath(), middlewares.JwtAuthMiddleware(app.GetConfig().AccessToken.AccessTokenSecret, ""))
	{
		userAndAdmin.GET("/details/:id", handler.FindUserDetailById)
	}
	admin := router.Group("/admin"+guest.BasePath(), middlewares.JwtAuthMiddleware(app.GetConfig().AccessToken.AccessTokenSecret, "admin"))
	{
		admin.DELETE("/delete/:id", handler.DeleteUserByAdmin)
	}
	return handler
}

func (userHttp *UserHttp) FindAll(c *gin.Context) {
	data, err := userHttp.userUsecase.GetAllData(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.ErrorWrapper{
			Message: "Failed To Get All Data ! ", Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (userHttp *UserHttp) LoginUsers(c *gin.Context) {
	user := &User.LoginPayload{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Parse Request ! ", Error: err.Error()})
		return
	}
	res, err := userHttp.userUsecase.LoginUser(c, user)
	if err != nil {
		if res != nil {
			c.JSON(http.StatusOK, errors.ErrorWrapper{Message: "Please Activated Your Account ! ", Error: err.Error(), Data: res.User.ID})
			return
		}
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Login ! ", Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (userHttp *UserHttp) RegisterUsers(c *gin.Context) {
	user := &User.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Error: err.Error()})
		return
	}
	res, err := userHttp.userUsecase.RegisterUser(c, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{ErrorS: err})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (userHttp *UserHttp) FindUserByToken(c *gin.Context) {
	userId := helpers.GetUserId(c)
	fmt.Println("UID: ", userId)
	data, err := userHttp.userUsecase.GetUserByUserId(c, userId.Hex())
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
	return
}

func (userHttp *UserHttp) FindUserDetailById(c *gin.Context) {
	userId := c.Param("id")
	data, err := userHttp.userUsecase.GetUserByUserId(c, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
	return
}

func (userHttp *UserHttp) UpdateUser(c *gin.Context) {
	// Parse the multipart form data
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Error parsing multipart form", Error: err.Error()})
		return
	}
	// Retrieve the file from the multipart form
	file, errFile := c.FormFile("file")
	if errFile != nil {
		if errFile != http.ErrMissingFile {
			c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed to get file", Error: errFile.Error()})
			return
		}
	}
	req, errs := c.GetPostForm("data")
	if !errs {
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed to Update User Bad Request"})
		return
	}
	data := &User.UpdateProfilePayload{}
	if err := json.Unmarshal([]byte(req), data); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed to Update User Bad Request", Error: err.Error()})
		return
	}
	data.ID = helpers.GetUserId(c)
	if file != nil {
		data, _ = image_helpers.UploadProfile(c, file, data)
	}
	res, err := userHttp.userUsecase.UpdateUserData(c, data)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed to Update User ! ", ErrorS: err})
		return
	}
	c.JSON(http.StatusOK, res)
	return
}

func (userHttp *UserHttp) UpdatePassword(c *gin.Context) {
	data := &User.UpdatePasswordPayload{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Update Password ! ", Error: err.Error()})
		return
	}
	data.Id = helpers.GetUserId(c)
	err := userHttp.userUsecase.UpdatePassword(c, data)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Update Password ! ", Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, errors.SuccessWrapper{Message: "Success To Update Password ! "})
	return
}

func (userHttp *UserHttp) AccountVerification(c *gin.Context) {
	data := &User.EmailVerifiedPayload{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Invalid Data Or Bad Request ! ", ErrorS: []string{err.Error()}})
		return
	}
	res, err := userHttp.userUsecase.VerifyEmailVerification(c, data, userHttp.userTokenHttp.userTokenRepo)
	if err != nil {
		if res != nil {
			c.JSON(http.StatusOK, errors.SuccessWrapper{Message: "Email Verified Already ! ", Data: res})
			return
		}
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Invalid To Verified Email, Otp Is Not Valid ! ", ErrorS: err})
		return
	}
	c.JSON(http.StatusOK, errors.SuccessWrapper{Message: "Success To Verify Email ! ", Data: res})
	return
}

func (userHttp *UserHttp) ResendVerificationOtp(c *gin.Context) {
	req := &User.ResendVerificationPayload{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Bad Data Request ! ", Error: err.Error()})
		return
	}
	data, err := userHttp.userUsecase.ResendVerificationRequest(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Invalid To Resend Otp ! ", ErrorS: err})
		return
	}
	c.JSON(http.StatusOK, errors.SuccessWrapper{Message: "Success To Resend Otp ! ", Data: data})
	return
}

func (userHttp *UserHttp) FindUserByIdForRequest(c *gin.Context, Id primitive.ObjectID) *User.User {
	userId := Id
	data, err := userHttp.userUsecase.GetUserByUserId(c, userId.Hex())
	if err != nil {
		logrus.Warnf("Failed to retrieve user for request, err : %v", err)
		return nil
	}
	return data
}

func (userHttp *UserHttp) DeleteUserByAdmin(c *gin.Context) {
	userId := helpers.ParseStringToObjectId(c.Param("id"))
	data, err := userHttp.userUsecase.DeleteUserById(c, &userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
	return
}

func (userHttp *UserHttp) DeleteUserAccount(c *gin.Context) {
	userId := helpers.GetUserId(c)
	data, err := userHttp.userUsecase.DeleteUserById(c, &userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
	return
}
