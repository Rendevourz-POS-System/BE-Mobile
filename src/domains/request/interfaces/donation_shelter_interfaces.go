package interfaces

import (
	"context"
	"github.com/midtrans/midtrans-go/coreapi"
	midtrans_interfaces "main.go/domains/payment/interfaces"
	Request "main.go/domains/request/entities"
)

type DonationShelterRepository interface {
	StoreOneDonation(ctx context.Context, req *Request.DonationShelter) (*Request.DonationShelter, error)
}

type DonationShelterUsecase interface {
	CreateDonation(ctx context.Context, req *Request.RequestResponse, midtransValidator midtrans_interfaces.MidtransUsecase) (*coreapi.ChargeResponse, error)
}
