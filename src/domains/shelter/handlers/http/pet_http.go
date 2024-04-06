package http

import (
	"github.com/gin-gonic/gin"
	_const "main.go/configs/const"
	"main.go/configs/database"
	Pet "main.go/domains/shelter/entities"
	"main.go/domains/shelter/interfaces"
	"main.go/domains/shelter/interfaces/impl/repository"
	"main.go/domains/shelter/interfaces/impl/usecase"
	"main.go/shared/helpers"
	"net/http"
)

type PetHttp struct {
	petUsecase interfaces.PetUseCase
}

func NewPetHttp(router *gin.Engine) *PetHttp {
	handler := &PetHttp{
		petUsecase: usecase.NewPetUseCase(repository.NewPetRepository(database.GetDatabase(_const.DB_SHELTER_APP))),
	}
	guest := router.Group("/pet")
	{
		guest.GET("", handler.GetAllPets)
	}
	return handler
}

func (h *PetHttp) GetAllPets(ctx *gin.Context) {
	search := &Pet.PetSearch{
		Search:   ctx.Query("search"),
		Page:     helpers.ParseStringToInt(ctx.Query("page")),
		PageSize: helpers.ParseStringToInt(ctx.Query("page_size")),
		Gender:   helpers.CheckPetGender(ctx.Query("gender")),
		Type:     ctx.Query("type"),
		Sort:     ctx.Query("sort"),
		OrderBy:  ctx.Query("order_by"),
		Location: ctx.Query("location"),
		AgeStart: helpers.ParseStringToInt(ctx.Query("age_start")),
		AgeEnd:   helpers.ParseStringToInt(ctx.Query("age_end")),
	}
	data, err := h.petUsecase.GetAllPets(search)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, data)
}
