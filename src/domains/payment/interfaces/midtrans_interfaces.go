package interfaces

import (
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type MidtransUsecase interface {
	ChargeRequest(chargeReq *coreapi.ChargeReq) (*coreapi.ChargeResponse, *midtrans.Error)
}

type MidtransCoreServices interface {
	CreateChargeRequest(reqMap *coreapi.ChargeReqWithMap) (coreapi.ResponseWithMap, *midtrans.Error)
}
