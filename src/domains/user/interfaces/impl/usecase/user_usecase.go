package usecase

import (
	"context"
	"fmt"
	User "main.go/domains/user/entities"
	"main.go/domains/user/interfaces"
	"main.go/shared/helpers"
)

type userUsecase struct {
	userRepo interfaces.UserRepository
}

func NewUserUsecase(userRepo interfaces.UserRepository) *userUsecase {
	return &userUsecase{userRepo}
}

func (u userUsecase) GetAllData(ctx context.Context) ([]User.User, error) {
	data, err := u.userRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (u userUsecase) RegisterUser(ctx context.Context, user *User.User) (res *User.User, errs []string) {
	var err error
	validate := helpers.NewValidator()
	if err = validate.Struct(user); err != nil {
		errs = helpers.CustomError(err)
		return nil, errs
	}
	user.IsActive = false
	user.Password, err = helpers.HashPassword(user.Password)
	if err != nil {
		errs = append(errs, err.Error())
		return nil, errs
	}
	data, checkUserData, err2 := u.userRepo.StoreOne(ctx, user)
	if err2 != nil {
		errs = append(errs, err2.Error())
		return nil, errs
	}
	if data != nil && !checkUserData {
		_, SendEmailVerification := u.SendEmailVerification(ctx, user)
		if SendEmailVerification != nil {
			errs = append(errs, SendEmailVerification.Error())
			return nil, errs
		}
		return data, nil
	} else if data != nil {
		return data, nil
	}
	return data, nil
}

func (u userUsecase) SendEmailVerification(ctx context.Context, data *User.User) (res *User.User, err error) {
	// send email verification
	return nil, nil
}

func (u userUsecase) LoginUser(ctx context.Context, userReq *User.LoginPayload) (res *User.LoginResponse, err error) {
	user, err := u.userRepo.FindByEmail(ctx, userReq.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user has not register yet ! ")
	}
	ok := helpers.ComparePassword(userReq.Password, user.Password)
	if !ok {
		return nil, fmt.Errorf("password or email doesn't match ! ")
	}
	if !user.IsActive {
		return nil, fmt.Errorf("user is not active ! ")
	}
	token, err := helpers.GenerateToken(user)
	if err != nil {
		return nil, err
	}
	res.Token = token
	res.Username = user.Username
	return res, nil
}
