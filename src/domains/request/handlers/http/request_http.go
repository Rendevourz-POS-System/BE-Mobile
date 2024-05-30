package http

import (
	"github.com/gin-gonic/gin"
	"main.go/configs/app"
	_const "main.go/configs/const"
	"main.go/configs/database"
	midtrans_interfaces "main.go/domains/payment/interfaces"
	Request "main.go/domains/request/entities"
	"main.go/domains/request/interfaces"
	"main.go/domains/request/interfaces/impl/repository"
	"main.go/domains/request/interfaces/impl/usecase"
	Shelter "main.go/domains/shelter/handlers/http"
	UserHttp "main.go/domains/user/handlers/http"
	"main.go/middlewares"
	"main.go/shared/helpers"
	"main.go/shared/message/errors"
	"net/http"
)

type RequestHttp struct {
	requestUsecase   interfaces.RequestUsecase
	midtransUsecase  midtrans_interfaces.MidtransUsecase
	donationHandlers *DonationShelterHttp
	adoptionHandlers *AdoptionShelterHttp
	userHandlers     *UserHttp.UserHttp
	shelterHandler   *Shelter.ShelterHttp
}

func NewRequestHttp(router *gin.Engine, midtrans midtrans_interfaces.MidtransUsecase, donationHandlers *DonationShelterHttp, adoptionHandlers *AdoptionShelterHttp, userHandlers *UserHttp.UserHttp, shelterHandlers *Shelter.ShelterHttp) *RequestHttp {
	handlers := &RequestHttp{
		requestUsecase: usecase.NewRequestUsecase(
			repository.NewRequestRepository(database.GetDatabase(_const.DB_SHELTER_APP)),
		),
		adoptionHandlers: adoptionHandlers,
		donationHandlers: donationHandlers,
		userHandlers:     userHandlers,
		shelterHandler:   shelterHandlers,
		midtransUsecase:  midtrans,
	}
	guest := router.Group("/request")
	{
		guest.GET("/")
	}
	user := router.Group(guest.BasePath(), middlewares.JwtAuthMiddleware(app.GetConfig().AccessToken.AccessTokenSecret))
	{
		user.POST("/create", handlers.CreateRequest)
		user.POST("/donation", handlers.CreateDonationRequest)
	}
	return handlers
}

func (RequestHttp *RequestHttp) CreateRequest(ctx *gin.Context) {
	req := &Request.Request{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Bad request Data !", Error: err.Error()})
		return
	}
	data, err := RequestHttp.requestUsecase.CreateRequest(ctx, req, RequestHttp.midtransUsecase)
	if err != nil {
		ctx.JSON(http.StatusExpectationFailed, errors.ErrorWrapper{Message: "Failed to create request ! ", ErrorS: err})
		return
	}
	ctx.JSON(http.StatusOK, errors.SuccessWrapper{Data: data, Message: "Created Request Successfully !"})
}

func (RequestHttp *RequestHttp) CreateDonationRequest(ctx *gin.Context) {
	req := &Request.DonationPayload{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Bad request Data !", Error: err.Error()})
		return
	}
	data, err := RequestHttp.requestUsecase.CreateDonationRequest(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusExpectationFailed, errors.ErrorWrapper{Message: "Failed to create request ! ", ErrorS: err})
		return
	}
	data.User = RequestHttp.userHandlers.FindUserByIdForRequest(ctx, helpers.GetUserId(ctx))
	data.UserTarget = RequestHttp.userHandlers.FindUserByIdForRequest(ctx, RequestHttp.shelterHandler.FindOneByShelterId(ctx, req.ShelterId))
	res, errDonation := RequestHttp.donationHandlers.donationShelterUsecase.CreateDonation(ctx, data, RequestHttp.midtransUsecase)
	if errDonation != nil {
		ctx.JSON(http.StatusExpectationFailed, errors.ErrorWrapper{Message: "Donation Request Failed ! ", ErrorS: err})
		return
	}
	ctx.JSON(http.StatusOK, errors.SuccessWrapper{Data: res, Message: "Created Request Successfully !"})
}
