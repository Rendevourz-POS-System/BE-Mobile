package usecase

import (
	"context"
	"fmt"
	"github.com/midtrans/midtrans-go/coreapi"
	Midtrans "main.go/domains/payment/interfaces"
	Request "main.go/domains/request/entities"
	"main.go/domains/request/interfaces"
	"main.go/domains/request/presistence"
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

func (u *requestUsecase) CreateDonationRequest(ctx context.Context, req *Request.DonationPayload, midtranValidator Midtrans.MidtransUsecase) (res *Request.Request, err []string) {
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
	if failedSendReq != nil {
		err = append(err, failedSendReq.Error())
		return nil, err
	}

	if presistence.Type(req.Type) == presistence.Donation {
		chargeReq := &coreapi.ChargeReq{
			PaymentType: coreapi.CoreapiPaymentType(req.PaymentType),
		}
		midtransResponse, midtransErr := midtranValidator.ChargeRequest(chargeReq)
		if midtransErr != nil {
			err = append(err, midtransErr.Error())
			return nil, err
		}
		fmt.Println("Midtrans Response: ", midtransResponse)
	}
	return res, nil
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
