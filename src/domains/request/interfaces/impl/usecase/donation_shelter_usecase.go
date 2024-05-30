package usecase

import (
	"context"
	"github.com/midtrans/midtrans-go/coreapi"
	midtrans_interfaces "main.go/domains/payment/interfaces"
	Request "main.go/domains/request/entities"
	"main.go/domains/request/interfaces"
	"main.go/domains/request/presistence"
	"main.go/shared/helpers"
)

type donationShelterUsecase struct {
	donationShelteRepo interfaces.DonationShelterRepository
}

func NewDonationShelterUsecase(donationShelter interfaces.DonationShelterRepository) *donationShelterUsecase {
	return &donationShelterUsecase{donationShelter}
}

func (u *donationShelterUsecase) CreateDonation(ctx context.Context, req *Request.RequestResponse, midtranValidator midtrans_interfaces.MidtransUsecase) (response *coreapi.ChargeResponse, err error) {
	donation := &Request.DonationShelter{
		RequestId:         req.Request.Id,
		Amount:            req.Donation.Amount,
		TransactionDate:   helpers.GetCurrentTime(nil),
		StatusTransaction: "new",
		PaymentType:       req.DonationPayload.PaymentType,
	}
	res, err := u.donationShelteRepo.StoreOneDonation(ctx, donation)
	if err != nil {
		return nil, err
	}
	req.Donation = res
	if presistence.Type(req.Request.Type) == presistence.Donation {
		response, err = midtranValidator.ChargeRequest(req)
		if err != nil {
			return nil, err
		}
	}
	return response, nil
}
