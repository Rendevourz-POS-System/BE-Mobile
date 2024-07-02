package http

import (
	"github.com/gin-gonic/gin"
	_const "main.go/src/configs/const"
	"main.go/src/configs/database"
	"main.go/src/domains/user/interfaces"
	"main.go/src/domains/user/interfaces/impl/repository"
	"main.go/src/domains/user/interfaces/impl/usecase"
	"net/http"
)

type UserTokenHttp struct {
	userTokenRepo interfaces.UserTokenUsecase
}

func NewUserTokenHttp(router *gin.Engine) *UserTokenHttp {
	handler := &UserTokenHttp{
		userTokenRepo: usecase.NewUserTokenUsecase(repository.NewUserTokenRepository(
			database.GetDatabase(_const.DB_SHELTER_APP))),
	}
	guest := router.Group("/user-token")
	{
		guest.POST("/token", handler.GenerateToken)
	}
	return handler
}

func (userTokenHttp *UserTokenHttp) GenerateToken(c *gin.Context) {
	token, err := userTokenHttp.userTokenRepo.GenerateToken()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
