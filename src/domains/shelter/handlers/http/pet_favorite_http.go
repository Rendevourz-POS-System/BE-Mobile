package http

import (
	"github.com/gin-gonic/gin"
	"main.go/configs/app"
	_const "main.go/configs/const"
	"main.go/configs/database"
	Shelter "main.go/domains/shelter/entities"
	"main.go/domains/shelter/interfaces"
	"main.go/domains/shelter/interfaces/impl/repository"
	"main.go/domains/shelter/interfaces/impl/usecase"
	"main.go/middlewares"
	"main.go/shared/helpers"
	"main.go/shared/message/errors"
	"net/http"
)

type PetFavoriteHttp struct {
	petFavoriteUsecase interfaces.PetFavoriteUseCase
}

func NewPetFavoriteHttp(router *gin.Engine) *PetFavoriteHttp {
	handler := &PetFavoriteHttp{
		petFavoriteUsecase: usecase.NewPetFavoriteUseCase(repository.NewPetFavoriteRepository(database.GetDatabase(_const.DB_SHELTER_APP))),
	}
	user := router.Group("/pet_favorite", middlewares.JwtAuthMiddleware(app.GetConfig().AccessToken.AccessTokenSecret, "user"))
	{
		user.POST("update", handler.UpdateData)
	}
	return handler
}

func (petFavorite *PetFavoriteHttp) UpdateData(ctx *gin.Context) {
	data := &Shelter.PetFavoriteCreate{}
	if err := ctx.ShouldBindJSON(data); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed to process request", Error: err.Error()})
		return
	}
	data.UserId = helpers.GetUserId(ctx)
	err := petFavorite.petFavoriteUsecase.UpdateIsFavoritePet(ctx, data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed to update data", Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, errors.SuccessWrapper{Message: "Success update data", Data: data})
}
