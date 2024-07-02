package http

import (
	"github.com/gin-gonic/gin"
	_const "main.go/src/configs/const"
	"main.go/src/configs/database"
	"main.go/src/domains/request/interfaces"
	"main.go/src/domains/request/interfaces/impl/repository"
	"main.go/src/domains/request/interfaces/impl/usecase"
)

type DonationShelterHttp struct {
	donationShelterUsecase interfaces.DonationShelterUsecase
}

func NewDonationShelterHttp(router *gin.Engine) *DonationShelterHttp {
	handlers := &DonationShelterHttp{
		donationShelterUsecase: usecase.NewDonationShelterUsecase(repository.NewDonationShelterRepository(database.GetDatabase(_const.DB_SHELTER_APP))),
	}
	guest := router.Group("/donation")
	{
		guest.GET("/")
	}
	return handlers
}
