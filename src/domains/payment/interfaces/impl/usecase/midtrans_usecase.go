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
	fmt.Println("Data Req --> ", req.UserTarget.ID, req.UserTarget.Username)
	chargeReq := &coreapi.ChargeReq{
		PaymentType: coreapi.CoreapiPaymentType(strings.ToLower(req.DonationPayload.PaymentType)),
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  req.Donation.Id.Hex(),
			GrossAmt: req.Donation.Amount,
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:           req.Request.Id.Hex(),
				Name:         req.User.Username,
				Price:        req.Donation.Amount,
				Qty:          1,
				Category:     req.Request.Type,
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
				Postcode:    helpers.PartIntToString(req.User.PostalCode),
				CountryCode: "IDN",
			},
			ShipAddr: &midtrans.CustomerAddress{
				FName:       req.UserTarget.Username,
				Phone:       req.UserTarget.PhoneNumber,
				Address:     req.UserTarget.Address,
				City:        req.UserTarget.City,
				Postcode:    helpers.PartIntToString(req.UserTarget.PostalCode),
				CountryCode: "IDN",
			},
		},
	}
	chargeReq, _ = m.paymentTypeSelector(chargeReq, req)
	midtransResponse, midtransErr := m.midtransCoreServices.CreateChargeRequest(chargeReq)
	if midtransErr != nil {
		return nil, midtransErr
	}
	return midtransResponse, nil
}

func (m *midtransUsecase) paymentTypeSelector(chargeReq *coreapi.ChargeReq, req *Request.RequestResponse) (*coreapi.ChargeReq, error) {
	//, coreapi.PaymentTypeGopay, coreapi.PaymentTypeShopeepay, coreapi.PaymentTypeQris
	chargeReq.CustomExpiry = &coreapi.CustomExpiry{
		OrderTime:      req.Donation.TransactionDate.Format("2006-01-02 15:04:05 -0700"),
		ExpiryDuration: 30,
		Unit:           "minutes",
	}
	switch coreapi.CoreapiPaymentType(strings.ToLower(req.DonationPayload.PaymentType)) {
	case coreapi.PaymentTypeBankTransfer:
		switch midtrans.Bank(strings.ToLower(*req.DonationPayload.PaymentChannel)) {
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
	case coreapi.PaymentTypeGopay:
		chargeReq.Gopay = &coreapi.GopayDetails{
			EnableCallback: false,
			CallbackUrl:    "",
			PreAuth:        false,
		}
		return chargeReq, nil
	case coreapi.PaymentTypeShopeepay:
		chargeReq.ShopeePay = &coreapi.ShopeePayDetails{CallbackUrl: ""}
		return chargeReq, nil
	}
	return chargeReq, nil
}
