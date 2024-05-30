package usecase

import (
	"context"
	"github.com/midtrans/midtrans-go/coreapi"
	midtrans_interfaces "main.go/domains/payment/interfaces"
	Request "main.go/domains/request/entities"
	"main.go/domains/request/interfaces"
	"main.go/domains/request/presistence"
	"main.go/shared/helpers"
	"strings"
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
		Amount:            req.DonationPayload.Amount,
		TransactionDate:   helpers.GetCurrentTime(nil),
		StatusTransaction: "new",
		PaymentType:       req.DonationPayload.PaymentType,
	}
	switch coreapi.CoreapiPaymentType(strings.ToLower(req.DonationPayload.PaymentType)) {
	case coreapi.PaymentTypeBankTransfer:
		donation.Channel = *req.DonationPayload.BankType
		break
	default:
		donation.Channel = *req.DonationPayload.EWalletType
	}
	res, err := u.donationShelteRepo.StoreOneDonation(ctx, donation)
	if err != nil {
		return nil, err
	}
	req.Donation = res
	if presistence.Type(strings.ToLower(req.DonationPayload.Type)) == presistence.Donation {
		response, err = midtranValidator.ChargeRequest(req)
		if err != nil {
			return nil, err
		}
	}
	return response, nil
}
