package http

import (
	"encoding/json"
	"fmt"
	"main.go/domains/request/presistence"
	Pet "main.go/domains/shelter/entities"
	"net/http"

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
)

type RequestHttp struct {
	requestUsecase   interfaces.RequestUsecase
	midtransUsecase  midtrans_interfaces.MidtransUsecase
	donationHandlers *DonationShelterHttp
	adoptionHandlers *AdoptionShelterHttp
	userHandlers     *UserHttp.UserHttp
	shelterHandler   *Shelter.ShelterHttp
	petHttp          *Shelter.PetHttp
}

func NewRequestHttp(router *gin.Engine, midtrans midtrans_interfaces.MidtransUsecase, donationHandlers *DonationShelterHttp, adoptionHandlers *AdoptionShelterHttp, userHandlers *UserHttp.UserHttp, shelterHandlers *Shelter.ShelterHttp, petHttp *Shelter.PetHttp) *RequestHttp {
	handlers := &RequestHttp{
		requestUsecase: usecase.NewRequestUsecase(
			repository.NewRequestRepository(database.GetDatabase(_const.DB_SHELTER_APP)),
		),
		adoptionHandlers: adoptionHandlers,
		donationHandlers: donationHandlers,
		userHandlers:     userHandlers,
		shelterHandler:   shelterHandlers,
		midtransUsecase:  midtrans,
		petHttp:          petHttp,
	}
	guest := router.Group("/request")
	{
		guest.GET("/")
	}
	user := router.Group(guest.BasePath(), middlewares.JwtAuthMiddleware(app.GetConfig().AccessToken.AccessTokenSecret, "user", "admin"))
	{
		user.GET("/find", handlers.FindAll)
		user.POST("/create", handlers.CreateRequest)
		user.POST("/donation", handlers.CreateDonationRequest)
		user.POST("/rescue_or_surrender", handlers.CreateRescueAndSurrender)
	}
	return handlers

}

func (RequestHttp *RequestHttp) FindAll(ctx *gin.Context) {
	searchReq := &Request.SearchRequestPayload{
		RequestId: helpers.ParseStringToObjectIdAddress(ctx.Query("request_id")),
		UserId:    helpers.ParseStringToObjectIdAddress(ctx.Query("user_id")),
		ShelterId: helpers.ParseStringToObjectIdAddress(ctx.Query("shelter_id")),
		Type:      helpers.ArrayAddress(ctx.QueryArray("type")),
		Reason:    nil,
		Status:    helpers.GetAddressString(ctx.Query("status")),
		Search:    helpers.GetAddressString(ctx.Query("search")),
		Page:      helpers.ParseStringToInt(ctx.Query("page")),
		PageSize:  helpers.ParseStringToInt(ctx.Query("page_size")),
	}
	data, err := RequestHttp.requestUsecase.GetAllData(ctx, searchReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed to get all data !", Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, errors.SuccessWrapper{Data: data})
}

func (RequestHttp *RequestHttp) CreateRequest(ctx *gin.Context) {
	req := &Request.Request{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Bad request Data !", Error: err.Error()})
		return
	}
	req.UserId = helpers.GetUserId(ctx)
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
	req.UserId = helpers.GetUserId(ctx)
	data, err := RequestHttp.requestUsecase.CreateDonationRequest(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusExpectationFailed, errors.ErrorWrapper{Message: "Failed to create request ! ", ErrorS: err})
		return
	}
	data.User = RequestHttp.userHandlers.FindUserByIdForRequest(ctx, req.UserId)
	data.UserTarget = RequestHttp.userHandlers.FindUserByIdForRequest(ctx, RequestHttp.shelterHandler.FindOneByShelterId(ctx, req.ShelterId))
	res, errDonation := RequestHttp.donationHandlers.donationShelterUsecase.CreateDonation(ctx, data, RequestHttp.midtransUsecase)
	if errDonation != nil {
		ctx.JSON(http.StatusExpectationFailed, errors.ErrorWrapper{Message: "Donation Request Failed ! ", ErrorS: err})
		return
	}
	ctx.JSON(http.StatusOK, errors.SuccessWrapper{Data: res, Message: res.StatusMessage})
}

func (RequestHttp *RequestHttp) CreateRescueAndSurrender(ctx *gin.Context) {
	request := &Request.CreateRescueAndSurrenderRequestPayload{}
	form, _ := ctx.MultipartForm()
	// Unmarshal the JSON data into the Pet struct
	jsonData := form.Value["request"][0]
	if errParse := json.Unmarshal([]byte(jsonData), &request.Request); errParse != nil {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Bind JSON Request ! ", Error: errParse.Error()})
		return
	}
	var (
		pet *Pet.Pet
		err error
	)
	if presistence.Type(request.Request.Type) == presistence.Rescue {
		pet, err = RequestHttp.petHttp.CreatePetForRescueAndSurenderPet(ctx)
		if err != nil {
			return
		}
	}
	request.Pet = pet
	request.Request = helpers.FillRequestData(pet, ctx)

	request.Request.UserId = helpers.GetUserId(ctx)
	data, errCreateReq := RequestHttp.requestUsecase.CreateRequest(ctx, request.Request, nil)
	if errCreateReq != nil {
		ctx.JSON(http.StatusExpectationFailed, errors.ErrorWrapper{Message: "Failed to create request ! ", ErrorS: errCreateReq})
		return
	}
	request.Request = data
	response := &Request.RescueAndSurrenderResponse{
		Pet:     request.Pet,
		Request: request.Request,
	}
	ctx.JSON(http.StatusOK, errors.SuccessWrapper{Message: fmt.Sprintf("Success Create %s Request !", request.Request.Type), Data: response})
}
