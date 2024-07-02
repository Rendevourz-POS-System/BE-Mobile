package http

import (
	"github.com/gin-gonic/gin"
	_const "main.go/src/configs/const"
	"main.go/src/configs/database"
	ShelterLocation "main.go/src/domains/master/entities"
	"main.go/src/domains/master/interfaces"
	"main.go/src/domains/master/interfaces/impl/repository"
	"main.go/src/domains/master/interfaces/impl/usecase"
	"main.go/src/shared/message/errors"
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
		guest.POST("/create", handler.CreateShelterLocation)
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

func (shelterLocationHttp *shelterLocationHttp) CreateShelterLocation(ctx *gin.Context) {
	req := []ShelterLocation.ShelterLocation{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Get All Shelter Location Data", Error: err.Error()})
		return
	}
	data, err := shelterLocationHttp.shelterLocationUsecase.CreateShelterLocation(ctx, req)
	if len(err) > 0 {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Create Shelter Location", ErrorS: err})
		return
	}
	ctx.JSON(http.StatusOK, data)
}
