package interfaces

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"main.go/domains/payment/interfaces"
	Request "main.go/domains/request/entities"
)

type RequestRepository interface {
	StoreOneRequest(ctx context.Context, req *Request.Request) (*Request.Request, error)
	FindAllRequest(ctx context.Context, req *Request.SearchRequestPayload) ([]Request.Request, error)
	FindOneRequestByData(ctx context.Context, req *bson.M) (res *Request.Request, err error)
	FindOneRequestById(ctx context.Context, Id *primitive.ObjectID) (res *Request.Request, err error)
	PutStatusRequestRescueOrSurrender(ctx context.Context, req *Request.UpdateRescueAndSurrenderRequestStatus) (res *Request.UpdateRescueAndSurrenderRequestStatusResponse, err []string)
	PutStatusRequestAdoption(ctx context.Context, req *Request.UpdateAdoptionRequestStatus) (res *Request.UpdateRescueAndSurrenderRequestStatusResponse, err []string)
}

type RequestUsecase interface {
	CreateRequest(ctx context.Context, req *Request.Request, payment interfaces.MidtransUsecase) (res *Request.Request, err []string)
	CreateDonationRequest(ctx context.Context, req *Request.DonationPayload) (res *Request.RequestResponse, err []string)
	GetAllData(ctx context.Context, req *Request.SearchRequestPayload) (res []Request.Request, err error)
	GetOneRequestByData(ctx context.Context, req *Request.UpdateRescueAndSurrenderRequestStatus) (res *Request.Request, err []string)
	GetOneRequestById(ctx context.Context, Id *primitive.ObjectID) (res *Request.Request, err error)
	UpdateStatusRequestRescueOrSurrender(ctx context.Context, req *Request.UpdateRescueAndSurrenderRequestStatus) (res *Request.UpdateRescueAndSurrenderRequestStatusResponse, err []string)
	UpdateStatusRequestAdoption(ctx context.Context, req *Request.UpdateAdoptionRequestStatus) (res *Request.UpdateRescueAndSurrenderRequestStatusResponse, err []string)
}
