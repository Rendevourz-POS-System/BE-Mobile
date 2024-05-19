package http

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"main.go/configs/app"
	_const "main.go/configs/const"
	"main.go/configs/database"
	User "main.go/domains/user/entities"
	"main.go/domains/user/interfaces"
	"main.go/domains/user/interfaces/impl/repository"
	"main.go/domains/user/interfaces/impl/usecase"
	"main.go/middlewares"
	"main.go/shared/helpers"
	"main.go/shared/message/errors"
	"net/http"
	"path/filepath"
)

type UserHttp struct {
	userUsecase interfaces.UserUsecase
}

func NewUserHttp(router *gin.Engine) *UserHttp {
	handler := &UserHttp{
		userUsecase: usecase.NewUserUsecase(repository.NewUserRepository(
			database.GetDatabase(_const.DB_SHELTER_APP))),
	}
	guest := router.Group("/user")
	{
		guest.GET("/", handler.FindAll)
		guest.POST("/register", handler.RegisterUsers)
		guest.POST("/login", handler.LoginUsers)
	}
	user := router.Group("/user", middlewares.JwtAuthMiddleware(app.GetConfig().AccessToken.AccessTokenSecret))
	{
		user.GET("data", handler.FindUserByToken)
		user.PUT("/update", handler.UpdateUser)
		user.PUT("/update-pw", handler.UpdatePassword)
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
	data, err := userHttp.userUsecase.GetUserByUserId(c, helpers.ToString(c.MustGet("x-user-id")))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
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
		FilePath := filepath.Join(app.GetConfig().Image.Folder, app.GetConfig().Image.UserPath, app.GetConfig().Image.ProfilePath, data.ID.Hex(), file.Filename)
		// Save the uploaded file with the temporary path
		if err := c.SaveUploadedFile(file, FilePath); err != nil {
			c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed to Upload Image !", Error: err.Error()})
			return
		}
		data.ImagePath = FilePath
	}
	res, err := userHttp.userUsecase.UpdateUserData(c, data)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed to Update User ! ", ErrorS: err})
		return
	}
	c.JSON(http.StatusOK, res)
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
}
