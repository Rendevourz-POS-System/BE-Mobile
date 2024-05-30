package usecase

import (
	"fmt"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"main.go/domains/payment/interfaces"
	Request "main.go/domains/request/entities"
	"main.go/shared/helpers"
	"strings"
)

type midtransUsecase struct {
	midtransCoreServices interfaces.MidtransCoreServices
}

func NewMidtransUsecase(coreServices interfaces.MidtransCoreServices) *midtransUsecase {
	return &midtransUsecase{coreServices}
}

func (m *midtransUsecase) ChargeRequest(req *Request.RequestResponse) (*coreapi.ChargeResponse, *midtrans.Error) {
	chargeReq := &coreapi.ChargeReq{
		PaymentType: coreapi.CoreapiPaymentType(req.DonationPayload.PaymentType),
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  req.Donation.Id.Hex(),
			GrossAmt: req.Donation.Amount,
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:           req.Request.Id.Hex(),
				Name:         "Donation",
				Price:        req.Donation.Amount,
				Qty:          1,
				Category:     "Donation",
				MerchantName: "Shelter-apps",
			},
		},
		CustomerDetails: &midtrans.CustomerDetails{
			FName: req.User.Username,
			Email: req.User.Email,
			Phone: req.User.PhoneNumber,
			BillAddr: &midtrans.CustomerAddress{
				FName:       req.User.Email,
				Phone:       req.User.PhoneNumber,
				Address:     req.User.Address,
				City:        req.User.City,
				Postcode:    helpers.ToString(req.User.PostalCode),
				CountryCode: "+62",
			},
			ShipAddr: &midtrans.CustomerAddress{
				FName:       req.UserTarget.Username,
				LName:       "",
				Phone:       req.UserTarget.PhoneNumber,
				Address:     req.UserTarget.Address,
				City:        req.UserTarget.City,
				Postcode:    helpers.ToString(req.UserTarget.PostalCode),
				CountryCode: "+62",
			},
		},
	}
	midtransResponse, midtransErr := m.midtransCoreServices.CreateChargeRequest(chargeReq)
	if midtransErr != nil {
		return nil, midtransErr
	}
	fmt.Println("Midtrans Response: ", midtransResponse)
	return midtransResponse, nil
}

func (m *midtransUsecase) paymentTypeSelector(chargeReq *coreapi.ChargeReq, req *Request.RequestResponse) (*coreapi.ChargeReq, error) {
	//, coreapi.PaymentTypeGopay, coreapi.PaymentTypeShopeepay, coreapi.PaymentTypeQris
	switch coreapi.CoreapiPaymentType(req.DonationPayload.PaymentType) {
	case coreapi.PaymentTypeBankTransfer:
		switch midtrans.Bank(strings.ToLower(*req.DonationPayload.BankType)) {
		case midtrans.BankBca:
			BankTransfer := &coreapi.BankTransferDetails{
				Bank: midtrans.BankBca,
			}
			chargeReq.BankTransfer = BankTransfer
			return chargeReq, nil
		case midtrans.BankBni:
			BankTransfer := &coreapi.BankTransferDetails{
				Bank: midtrans.BankBni,
			}
			chargeReq.BankTransfer = BankTransfer
			return chargeReq, nil
		case midtrans.BankBri:
			BankTransfer := &coreapi.BankTransferDetails{
				Bank: midtrans.BankBri,
			}
			chargeReq.BankTransfer = BankTransfer
			return chargeReq, nil
		case midtrans.BankMandiri:
			BankTransfer := &coreapi.BankTransferDetails{
				Bank: midtrans.BankMandiri,
			}
			chargeReq.BankTransfer = BankTransfer
			return chargeReq, nil
		case midtrans.BankCimb:
			BankTransfer := &coreapi.BankTransferDetails{
				Bank: midtrans.BankCimb,
			}
			chargeReq.BankTransfer = BankTransfer
			return chargeReq, nil
		case midtrans.BankMaybank:
			BankTransfer := &coreapi.BankTransferDetails{
				Bank: midtrans.BankMaybank,
			}
			chargeReq.BankTransfer = BankTransfer
			return chargeReq, nil
		case midtrans.BankPermata:
			BankTransfer := &coreapi.BankTransferDetails{
				Bank: midtrans.BankPermata,
			}
			chargeReq.BankTransfer = BankTransfer
			return chargeReq, nil
		case midtrans.BankMega:
			BankTransfer := &coreapi.BankTransferDetails{
				Bank: midtrans.BankMega,
			}
			chargeReq.BankTransfer = BankTransfer
			return chargeReq, nil
		}
	}
	return nil, nil
}
