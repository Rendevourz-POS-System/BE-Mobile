package usecase

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	JwtEmailClaims "main.go/domains/user/entities"
	"main.go/domains/user/interfaces"
)

type userTokenUsecase struct {
	userTokenRepo interfaces.UserTokenRepository
}

func NewUserTokenUsecase(userTokenRepo interfaces.UserTokenRepository) *userTokenUsecase {
	return &userTokenUsecase{userTokenRepo}
}

func (u *userTokenUsecase) GenerateToken() (string, error) {
	return "", nil
}

func (u *userTokenUsecase) FindValidToken(ctx context.Context, claims *JwtEmailClaims.JwtEmailClaims) (*primitive.ObjectID, error) {
	userId, err := u.userTokenRepo.FindOneUserTokenByNonce(ctx, claims)
	if err != nil {
		return nil, err
	}
	return userId, nil
}

func (u *userTokenUsecase) FindValidTokenByUserId(ctx context.Context, userId *primitive.ObjectID, Otp *int) (*primitive.ObjectID, error) {
	findUserId, err := u.userTokenRepo.FindValidTokenByUserId(ctx, userId, Otp)
	if err != nil {
		return nil, err
	}
	return findUserId, nil
}
