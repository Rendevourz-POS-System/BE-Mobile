package http

import (
	"fmt"
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
	"net/http"
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
		guest.GET("data", middlewares.JwtAuthMiddleware(app.GetConfig().AccessToken.AccessTokenSecret), handler.FindUserByToken)
	}
	return handler
}

func (userHttp *UserHttp) FindAll(c *gin.Context) {
	data, err := userHttp.userUsecase.GetAllData(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (userHttp *UserHttp) LoginUsers(c *gin.Context) {
	user := &User.LoginPayload{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := userHttp.userUsecase.LoginUser(c, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (userHttp *UserHttp) RegisterUsers(c *gin.Context) {
	user := &User.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := userHttp.userUsecase.RegisterUser(c, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": err,
		})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (userHttp *UserHttp) FindUserByToken(c *gin.Context) {
	userId := helpers.ToString(c.MustGet("x-user-id"))
	fmt.Println("UserId: ", userId)
	data, err := userHttp.userUsecase.GetUserByUserId(c, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, data)
}
