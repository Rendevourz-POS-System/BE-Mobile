package usecase

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
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

func (u userUsecase) RegisterUser(ctx context.Context, user *User.User) (res *User.User, err error) {
	validate := helpers.NewValidator()
	if err := validate.Struct(user); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err.Namespace(), err.Field(), err.StructNamespace(), err.StructField(), err.Tag(), err.ActualTag(), err.Kind(), err.Type(), err.Value(), err.Param())
		}
		return nil, err
	}
	user.IsActive = false
	user.Password, err = helpers.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	data, checkUserData, err2 := u.userRepo.StoreOne(ctx, user)
	if err2 != nil {
		return nil, err2
	}
	if checkUserData {
		// Email already exists
	} else {

	}
	return data, nil
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
		return nil, fmt.Errorf("password not match ! ")
	}
	//token, err := helpers.GenerateToken()
	//if err != nil {
	//	return nil, err
	//}
	//res.Token = token
	return res, nil
}
