package usecase

import (
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"main.go/domains/payment/interfaces"
)

type midtransUsecase struct {
	midtransCoreServices interfaces.MidtransCoreServices
}

func NewMidtransUsecase(coreServices interfaces.MidtransCoreServices) *midtransUsecase {
	return &midtransUsecase{coreServices}
}

func (m *midtransUsecase) ChargeRequest(chargeReq *coreapi.ChargeReq) (*coreapi.ChargeResponse, *midtrans.Error) {

	return nil, nil
}
