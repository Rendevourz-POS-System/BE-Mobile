package interfaces

import (
	"context"
	"main.go/domains/payment/interfaces"
	Request "main.go/domains/request/entities"
)

type RequestRepository interface {
	StoreOneRequest(ctx context.Context, req *Request.Request) (*Request.Request, error)
}

type RequestUsecase interface {
	CreateRequest(ctx context.Context, req *Request.Request, payment interfaces.MidtransUsecase) (res *Request.Request, err []string)
	CreateDonationRequest(ctx context.Context, req *Request.DonationPayload, payment interfaces.MidtransUsecase) (res *Request.Request, err []string)
}
