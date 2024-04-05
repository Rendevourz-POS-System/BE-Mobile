package usecase

import (
	"context"
	"errors"
	"fmt"
	User "main.go/domains/user/entities"
	"main.go/domains/user/interfaces"
	"main.go/domains/user/mail/controller"
	"main.go/shared/helpers"
)

type userUsecase struct {
	userRepo interfaces.UserRepository
}

func NewUserUsecase(userRepo interfaces.UserRepository) *userUsecase {
	return &userUsecase{userRepo}
}

func (u *userUsecase) GetAllData(ctx context.Context) ([]User.User, error) {
	data, err := u.userRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (u *userUsecase) RegisterUser(ctx context.Context, user *User.User) (res *User.User, errs []string) {
	var err error
	validate := helpers.NewValidator()
	if err = validate.Struct(user); err != nil {
		errs = helpers.CustomError(err)
		return nil, errs
	}
	datas := u.setDefaultUserData(user)
	if err != nil {
		errs = append(errs, err.Error())
		return nil, errs
	}
	resData, checkUserData, err2 := u.userRepo.StoreOne(ctx, datas)
	if err2 != nil {
		errs = append(errs, err2.Error())
		return nil, errs
	}
	if resData != nil && !checkUserData {
		secretCode, err := u.userRepo.GenerateAndStoreToken(ctx, resData.ID, resData.Email)
		if err != nil {
			errs = append(errs, err.Error())
			return nil, errs
		}
		_, SendEmailVerification := u.SendEmailVerification(ctx, user, secretCode)
		if SendEmailVerification != nil {
			errs = append(errs, SendEmailVerification.Error())
			return nil, errs
		}
		return resData, nil
	}
	return resData, nil
}

func (u *userUsecase) setDefaultUserData(user *User.User) *User.User {
	StaffSatus := helpers.CheckStaffStatus(user.Role)
	return &User.User{
		Nik:                user.Nik,
		Email:              user.Email,
		Username:           user.Username,
		PostalCode:         user.PostalCode,
		Province:           user.Province,
		PhoneNumber:        user.PhoneNumber,
		Password:           helpers.HashPassword(user.Password),
		Address:            user.Address,
		City:               user.City,
		Verified:           false,
		ShelterIsActivated: false,
		StaffStatus:        StaffSatus,
		Role:               helpers.GetRole(user.Role),
		CreatedAt:          helpers.GetCurrentTime(nil),
	}
}

func (u *userUsecase) SendEmailVerification(ctx context.Context, data *User.User, secretCode string) (res *User.User, err error) {
	// send email verification
	ok := controller.SendEmail(&User.MailSend{
		To:      data.Email,
		Subject: "Email Verification",
		Content: helpers.GetVerifiedUrl(secretCode, data.Email),
		Cc:      "",
		Bcc:     "",
		Attach:  "readme.txt",
	})
	if ok != nil {
		return nil, fmt.Errorf("failed to send email verification ! ")
	}
	return res, nil
}

func (u *userUsecase) LoginUser(ctx context.Context, userReq *User.LoginPayload) (*User.LoginResponse, error) {
	var res = &User.LoginResponse{}
	user, err := u.userRepo.FindByEmail(ctx, userReq.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user has not register yet ! ")
	}
	ok := helpers.ComparePassword(userReq.Password, user.Password)
	if !ok {
		return nil, errors.New("password or email doesn't match ! ")
	}
	if !user.Verified {
		return nil, errors.New("user is not active ! ")
	}
	token, err := helpers.GenerateToken(user)
	if err != nil {
		return nil, err
	}
	res.Token = token
	res.Username = user.Username
	return res, nil
}

func (u *userUsecase) GetUserByUserId(ctx context.Context, id string) (*User.User, error) {
	user, err := u.userRepo.FindUserById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
