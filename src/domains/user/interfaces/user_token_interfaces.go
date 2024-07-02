package interfaces

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	JwtEmailClaims "main.go/src/domains/user/entities"
)

type UserTokenRepository interface {
	StoreToken(token string) error
	FindOneUserTokenByNonce(ctx context.Context, claims *JwtEmailClaims.JwtEmailClaims) (*primitive.ObjectID, error)
	FindValidTokenByUserId(ctx context.Context, userId *primitive.ObjectID, Otp *int) (*primitive.ObjectID, error)
	//VerifiedOtpByUserId(ctx context.Context, userId *primitive.ObjectID, Otp *int) (*primitive.ObjectID, error)
}

type UserTokenUsecase interface {
	GenerateToken() (string, error)
	FindValidToken(ctx context.Context, claims *JwtEmailClaims.JwtEmailClaims) (*primitive.ObjectID, error)
	FindValidTokenByUserId(ctx context.Context, userId *primitive.ObjectID, Otp *int) (*primitive.ObjectID, error)
}
