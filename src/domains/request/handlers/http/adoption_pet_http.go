package http

import (
	"github.com/gin-gonic/gin"
	_const "main.go/configs/const"
	"main.go/configs/database"
	"main.go/domains/request/interfaces"
	"main.go/domains/request/interfaces/impl/repository"
	"main.go/domains/request/interfaces/impl/usecase"
)

type AdoptionShelterHttp struct {
	adoptionShelterUsecase interfaces.AdoptionPetUsecase
}

func NewAdoptionShelterHttp(router *gin.Engine) *AdoptionShelterHttp {
	hanlders := &AdoptionShelterHttp{
		adoptionShelterUsecase: usecase.NewAdoptionPetUsecase(repository.NewAdoptionPetRepository(database.GetDatabase(_const.DB_SHELTER_APP))),
	}
	return hanlders
}
