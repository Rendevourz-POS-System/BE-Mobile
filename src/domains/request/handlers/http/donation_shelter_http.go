package http

import (
	"github.com/gin-gonic/gin"
	_const "main.go/configs/const"
	"main.go/configs/database"
	"main.go/domains/request/interfaces"
	"main.go/domains/request/interfaces/impl/repository"
	"main.go/domains/request/interfaces/impl/usecase"
)

type DonationShelterHttp struct {
	donationShelterUsecase interfaces.DonationShelterUsecase
}

func NewDonationShelterHttp(router *gin.Engine) *DonationShelterHttp {
	handlers := &DonationShelterHttp{
		donationShelterUsecase: usecase.NewDonationShelterUsecase(repository.NewDonationShelterRepository(database.GetDatabase(_const.DB_SHELTER_APP))),
	}
	return handlers
}
