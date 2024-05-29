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
	"main.go/middlewares"
	"main.go/shared/message/errors"
	"net/http"
)

type RequestHttp struct {
	requestUsecase  interfaces.RequestUsecase
	midtransUsecase midtrans_interfaces.MidtransUsecase
}

func NewRequestHttp(router *gin.Engine, midtrans midtrans_interfaces.MidtransUsecase) *RequestHttp {
	handlers := &RequestHttp{
		requestUsecase: usecase.NewRequestUsecase(
			repository.NewRequestRepository(database.GetDatabase(_const.DB_SHELTER_APP)),
		),
		midtransUsecase: midtrans,
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
	data, err := RequestHttp.requestUsecase.CreateDonationRequest(ctx, req, RequestHttp.midtransUsecase)
	if err != nil {
		ctx.JSON(http.StatusExpectationFailed, errors.ErrorWrapper{Message: "Failed to create request ! ", ErrorS: err})
		return
	}
	ctx.JSON(http.StatusOK, errors.SuccessWrapper{Data: data, Message: "Created Request Successfully !"})
}
