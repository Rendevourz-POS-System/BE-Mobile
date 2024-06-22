package interfaces

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"main.go/domains/payment/interfaces"
	Request "main.go/domains/request/entities"
)

type RequestRepository interface {
	StoreOneRequest(ctx context.Context, req *Request.Request) (*Request.Request, error)
	FindAllRequest(ctx context.Context, req *Request.SearchRequestPayload) ([]Request.Request, error)
	FindOneRequestByData(ctx context.Context, req *bson.M) (res *Request.Request, err error)
}

type RequestUsecase interface {
	CreateRequest(ctx context.Context, req *Request.Request, payment interfaces.MidtransUsecase) (res *Request.Request, err []string)
	CreateDonationRequest(ctx context.Context, req *Request.DonationPayload) (res *Request.RequestResponse, err []string)
	GetAllData(ctx context.Context, req *Request.SearchRequestPayload) (res []Request.Request, err error)
	GetOneRequestByData(ctx context.Context, req *Request.UpdateRescueAndSurrenderRequestStatus) (res *Request.Request, err []string)
	//UpdateStatusRequest(ctx context.Context, req *Request.UpdateRescueAndSurrenderRequestStatus) (res *Request.Request, err []string)
}
