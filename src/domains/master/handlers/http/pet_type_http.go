package http

import (
	"github.com/gin-gonic/gin"
	_const "main.go/configs/const"
	"main.go/configs/database"
	PetType "main.go/domains/master/entities"
	"main.go/domains/master/interfaces"
	"main.go/domains/master/interfaces/impl/repository"
	"main.go/domains/master/interfaces/impl/usecase"
	"main.go/shared/message/errors"
	"net/http"
)

type PetTypeHttp struct {
	petTypeUsecase interfaces.PetTypeUsecase
}

func NewPetTypeHttp(router *gin.Engine) *PetTypeHttp {
	handler := &PetTypeHttp{
		petTypeUsecase: usecase.NewPetTypeUsecase(repository.NewPetTypeRepo(database.GetDatabase(_const.DB_SHELTER_APP))),
	}
	guest := router.Group("/pet-types")
	{
		guest.GET("/", handler.GetAllPetTypes)
		guest.POST("/create", handler.CreatePetType)
	}
	return handler
}

func (petTypeHttp *PetTypeHttp) CreatePetType(ctx *gin.Context) {
	req := &PetType.PetType{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Bad request", Error: err.Error()})
		return
	}
	data, err := petTypeHttp.petTypeUsecase.CreatePetType(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Create Pet Type", Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, data)
}

func (petTypeHttp *PetTypeHttp) GetAllPetTypes(ctx *gin.Context) {
	data, err := petTypeHttp.petTypeUsecase.GetAllPetTypes(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed to get all pet types !", Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, data)
}
