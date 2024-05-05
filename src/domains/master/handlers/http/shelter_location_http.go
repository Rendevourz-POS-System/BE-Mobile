package http

import (
	"github.com/gin-gonic/gin"
	_const "main.go/configs/const"
	"main.go/configs/database"
	"main.go/domains/master/interfaces"
	"main.go/domains/master/interfaces/impl/repository"
	"main.go/domains/master/interfaces/impl/usecase"
	"main.go/shared/message/errors"
	"net/http"
)

type shelterLocationHttp struct {
	shelterLocationUsecase interfaces.ShelterLocationUsecase
}

func NewShelterLocationHttp(router *gin.Engine) *shelterLocationHttp {
	handler := &shelterLocationHttp{
		usecase.NewShelterLocationUsecase(
			repository.NewShelterLocationRepository(database.GetDatabase(_const.DB_SHELTER_APP))),
	}
	guest := router.Group("/shelter-location")
	{
		guest.GET("/", handler.GetAllLocation)
	}
	return handler
}

func (shelterLocationHttp *shelterLocationHttp) GetAllLocation(ctx *gin.Context) {
	data, err := shelterLocationHttp.shelterLocationUsecase.GetAllShelterLocation(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Get All Shelter Location Data", Error: err.Error()})
	}
	ctx.JSON(http.StatusOK, data)
}
