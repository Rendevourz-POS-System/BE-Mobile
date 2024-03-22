package usecase

import "main.go/domains/user/interfaces"

type userTokenUsecase struct {
	userTokenRepo interfaces.UserTokenRepository
}

func NewUserTokenUsecase(userTokenRepo interfaces.UserTokenRepository) *userTokenUsecase {
	return &userTokenUsecase{userTokenRepo}
}

func (u *userTokenUsecase) GenerateToken() (string, error) {
	return "", nil
}
