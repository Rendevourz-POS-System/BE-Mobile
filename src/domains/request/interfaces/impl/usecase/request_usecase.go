package usecase

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	Midtrans "main.go/domains/payment/interfaces"
	Request "main.go/domains/request/entities"
	"main.go/domains/request/interfaces"
	"main.go/shared/helpers"
	"strings"
)

type requestUsecase struct {
	requestRepo interfaces.RequestRepository
}

func NewRequestUsecase(requestRepo interfaces.RequestRepository) *requestUsecase {
	return &requestUsecase{requestRepo}
}

func (u *requestUsecase) CreateRequest(ctx context.Context, req *Request.Request, midtranValidator Midtrans.MidtransUsecase) (res *Request.Request, err []string) {
	validate := helpers.NewValidator()
	if errs := validate.Struct(req); errs != nil {
		err = helpers.CustomError(errs)
		return nil, err
	}
	res, failedSendReq := u.requestRepo.StoreOneRequest(ctx, u.fillDefaultRequest(req))
	if failedSendReq != nil {
		err = append(err, failedSendReq.Error())
		return nil, err
	}
	return res, nil
}

func (u *requestUsecase) CreateDonationRequest(ctx context.Context, req *Request.DonationPayload) (res *Request.RequestResponse, err []string) {
	validate := helpers.NewValidator()
	if errs := validate.Struct(req); errs != nil {
		err = helpers.CustomError(errs)
		return nil, err
	}
	req.Type = strings.ToLower(req.Type)
	requestRes, failedSendReq := u.requestRepo.StoreOneRequest(ctx, &Request.Request{
		UserId:      req.UserId,
		ShelterId:   req.ShelterId,
		Type:        req.Type,
		Status:      "Ongoing",
		Reason:      req.Reason,
		RequestedAt: helpers.GetCurrentTime(nil),
	})
	if failedSendReq != nil {
		err = append(err, failedSendReq.Error())
		return nil, err
	}
	// Initialize res before using it
	res = &Request.RequestResponse{}
	res.Request = requestRes
	res.DonationPayload = req
	return res, nil
}

func (u *requestUsecase) GetAllData(ctx context.Context, req *Request.SearchRequestPayload) (res []Request.Request, err error) {
	res, err = u.requestRepo.FindAllRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u *requestUsecase) fillDefaultRequest(req *Request.Request) *Request.Request {
	var petId *primitive.ObjectID
	if req.PetId.Hex() != "" {
		petId = req.PetId
	}
	request := "Ongoing"
	if req.Status != "" {
		request = req.Status
	}
	return &Request.Request{
		UserId:      req.UserId,
		ShelterId:   req.ShelterId,
		Type:        req.Type,
		Status:      request,
		Reason:      req.Reason,
		PetId:       petId,
		RequestedAt: helpers.GetCurrentTime(nil),
	}
}

func (u *requestUsecase) GetOneRequestByData(ctx context.Context, req *Request.UpdateRescueAndSurrenderRequestStatus) (res *Request.Request, err []string) {
	validate := helpers.NewValidator()
	if errs := validate.Struct(req); errs != nil {
		err = helpers.CustomError(errs)
		return nil, err
	}
	data := &bson.M{
		"_id":    req.RequestId,
		"status": "Ongoing",
	}
	responseData, errs := u.requestRepo.FindOneRequestByData(ctx, data)
	if errs != nil {
		err = append(err, errs.Error())
		return nil, err
	}
	return responseData, nil
}

func (u *requestUsecase) UpdateStatusRequest(ctx context.Context, req *Request.UpdateRescueAndSurrenderRequestStatus) (res *Request.UpdateRescueAndSurrenderRequestStatusResponse, err []string) {
	res, err = u.requestRepo.PutStatusRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
