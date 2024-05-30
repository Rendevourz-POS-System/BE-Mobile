package usecase

import (
	"context"
	Midtrans "main.go/domains/payment/interfaces"
	Request "main.go/domains/request/entities"
	"main.go/domains/request/interfaces"
	"main.go/shared/helpers"
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

func (u *requestUsecase) CreateDonationRequest(ctx context.Context, req *Request.DonationPayload) (response *Request.RequestResponse, err []string) {
	validate := helpers.NewValidator()
	if errs := validate.Struct(req); errs != nil {
		err = helpers.CustomError(errs)
		return nil, err
	}
	res, failedSendReq := u.requestRepo.StoreOneRequest(ctx, &Request.Request{
		UserId:      req.UserId,
		ShelterId:   req.ShelterId,
		Type:        req.Type,
		Status:      req.Status,
		Reason:      req.Reason,
		RequestedAt: helpers.GetCurrentTime(nil),
	})
	response.Request = res
	response.DonationPayload = req
	if failedSendReq != nil {
		err = append(err, failedSendReq.Error())
		return nil, err
	}
	return response, nil
}

func (u *requestUsecase) fillDefaultRequest(req *Request.Request) *Request.Request {
	return &Request.Request{
		UserId:      req.UserId,
		ShelterId:   req.ShelterId,
		Type:        req.Type,
		Status:      req.Status,
		Reason:      req.Reason,
		RequestedAt: helpers.GetCurrentTime(nil),
	}
}
