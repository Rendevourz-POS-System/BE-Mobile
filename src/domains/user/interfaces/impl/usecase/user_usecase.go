package usecase

import (
	"context"
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

func (u userUsecase) RegisterUser(ctx context.Context, user *User.User) (*User.User, error) {
	validate := helpers.NewValidator()
	if err := validate.Struct(user); err != nil {
		return nil, err
	}
	data, err := u.userRepo.StoreOne(ctx, user)
	if err != nil {
		return nil, err
	}
	return data, nil
}
