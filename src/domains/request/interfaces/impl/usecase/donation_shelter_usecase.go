package usecase

import (
	"context"
	"errors"
	"fmt"
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

func (u *donationShelterUsecase) CreateDonation(ctx context.Context, req *Request.RequestResponse, midtranValidator midtrans_interfaces.MidtransUsecase) (*coreapi.ChargeResponse, error) {
	req.DonationPayload.PaymentType = strings.ToLower(req.DonationPayload.PaymentType)
	donation := &Request.DonationShelter{
		RequestId:         req.Request.Id,
		Amount:            req.DonationPayload.Amount,
		TransactionDate:   helpers.GetCurrentTime(nil),
		StatusTransaction: "new",
		PaymentType:       req.DonationPayload.PaymentType,
		CreatedAt:         helpers.GetCurrentTime(nil),
	}
	switch coreapi.CoreapiPaymentType(req.DonationPayload.PaymentType) {
	case coreapi.PaymentTypeBankTransfer:
		donation.Channel = *req.DonationPayload.BankType
		break
	default:
		donation.Channel = *req.DonationPayload.EWalletType
	}
	res, errs := u.donationShelteRepo.StoreOneDonation(ctx, donation)
	if errs != nil {
		return nil, errs
	}
	req.Donation = res
	fmt.Println("PaymentType: ", req.DonationPayload.PaymentType)
	if presistence.Type(req.DonationPayload.Type) != presistence.Donation {
		return nil, errors.New("Type must be donation !")
	}
	response, err := midtranValidator.ChargeRequest(req)
	if err != nil {
		return nil, err
	}
	return response, nil
}
