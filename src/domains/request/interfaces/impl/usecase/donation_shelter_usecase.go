package usecase

import (
	"context"
	"errors"
	"github.com/midtrans/midtrans-go/coreapi"
	midtrans_interfaces "main.go/src/domains/payment/interfaces"
	Request "main.go/src/domains/request/entities"
	"main.go/src/domains/request/interfaces"
	"main.go/src/domains/request/presistence"
	"main.go/src/shared/helpers"
	"strings"
)

type donationShelterUsecase struct {
	donationShelteRepo interfaces.DonationShelterRepository
}

func NewDonationShelterUsecase(donationShelter interfaces.DonationShelterRepository) *donationShelterUsecase {
	return &donationShelterUsecase{donationShelter}
}

func (u *donationShelterUsecase) CreateDonation(ctx context.Context, req *Request.RequestResponse, midtranValidator midtrans_interfaces.MidtransUsecase) (*coreapi.ChargeResponse, error) {
	req.DonationPayload.PaymentType = strings.ToLower(req.DonationPayload.PaymentType)
	donation := &Request.DonationShelter{
		RequestId:         req.Request.Id,
		Amount:            req.DonationPayload.Amount,
		TransactionDate:   helpers.GetCurrentTime(nil),
		StatusTransaction: "new",
		Channel:           *req.DonationPayload.PaymentChannel,
		PaymentType:       req.DonationPayload.PaymentType,
		CreatedAt:         helpers.GetCurrentTime(nil),
	}
	res, errs := u.donationShelteRepo.StoreOneDonation(ctx, donation)
	if errs != nil {
		return nil, errs
	}
	req.Donation = res
	if presistence.Type(req.DonationPayload.Type) != presistence.Donation {
		return nil, errors.New("Type must be donation !")
	}
	response, err := midtranValidator.ChargeRequest(req)
	if err != nil {
		return nil, err
	}
	return response, nil
}
