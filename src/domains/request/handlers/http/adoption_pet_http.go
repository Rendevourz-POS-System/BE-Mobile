package http

import (
	"github.com/gin-gonic/gin"
	"main.go/configs/app"
	_const "main.go/configs/const"
	"main.go/configs/database"
	"main.go/domains/request/interfaces"
	"main.go/domains/request/interfaces/impl/repository"
	"main.go/domains/request/interfaces/impl/usecase"
	"main.go/middlewares"
)

type AdoptionShelterHttp struct {
	adoptionShelterUsecase interfaces.AdoptionPetUsecase
}

func NewAdoptionShelterHttp(router *gin.Engine) *AdoptionShelterHttp {
	handlers := &AdoptionShelterHttp{
		adoptionShelterUsecase: usecase.NewAdoptionPetUsecase(repository.NewAdoptionPetRepository(database.GetDatabase(_const.DB_SHELTER_APP))),
	}
	guest := router.Group("/adoption")
	{
		guest.GET("/")
	}
	user := router.Group(guest.BasePath(), middlewares.JwtAuthMiddleware(app.GetConfig().AccessToken.AccessTokenSecret, "user", "admin"))
	{
		user.POST("/create")
	}
	return handlers
}
