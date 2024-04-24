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

type ShelterFavoriteHttp struct {
	shelterFavoriteUseCase interfaces.ShelterFavoriteUseCase
}

func NewShelterFavoriteHttp(router *gin.Engine) *ShelterFavoriteHttp {
	handler := &ShelterFavoriteHttp{
		shelterFavoriteUseCase: usecase.NewShelterFavoriteUseCase(repository.NewShelterFavoriteRepository(database.GetDatabase(_const.DB_SHELTER_APP))),
	}
	user := router.Group("/shelter_favorite", middlewares.JwtAuthMiddleware(app.GetConfig().AccessToken.AccessTokenSecret))
	{
		user.POST("update", handler.UpdateData)
	}
	return handler
}

func (shelterFavorite *ShelterFavoriteHttp) UpdateData(ctx *gin.Context) {
	data := &Shelter.ShelterFavoriteCreate{}
	if err := ctx.ShouldBindJSON(data); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed to process request", Error: err.Error()})
		return
	}
	data.UserId = helpers.ParseObjectIdToString(helpers.GetUserId(ctx))
	err := shelterFavorite.shelterFavoriteUseCase.UpdateIsFavoriteShelter(ctx, data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed to update data", Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, errors.SuccessWrapper{Message: "Success update data", Data: data})
}
