package interfaces

import (
	"context"
	"main.go/domains/payment/interfaces"
	Request "main.go/domains/request/entities"
)

type RequestRepository interface {
	StoreOneRequest(ctx context.Context, req *Request.Request) (*Request.Request, error)
	FindAllRequest(ctx context.Context, req *Request.SearchRequestPayload) ([]Request.Request, error)
}

type RequestUsecase interface {
	CreateRequest(ctx context.Context, req *Request.Request, payment interfaces.MidtransUsecase) (res *Request.Request, err []string)
	CreateDonationRequest(ctx context.Context, req *Request.DonationPayload) (res *Request.RequestResponse, err []string)
	GetAllData(ctx context.Context, req *Request.SearchRequestPayload) (res []Request.Request, err error)
}
