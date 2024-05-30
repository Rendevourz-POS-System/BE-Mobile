package interfaces

import (
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	Request "main.go/domains/request/entities"
)

type MidtransCoreServices interface {
	CreateChargeRequest(reqMap *coreapi.ChargeReq) (*coreapi.ChargeResponse, *midtrans.Error)
}

type MidtransUsecase interface {
	ChargeRequest(chargeReq *Request.RequestResponse) (*coreapi.ChargeResponse, *midtrans.Error)
}
