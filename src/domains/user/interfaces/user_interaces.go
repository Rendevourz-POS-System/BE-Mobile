package interfaces

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	User "main.go/domains/user/entities"
)

type UserRepository interface {
	FindAll(ctx context.Context) ([]User.User, error)
	StoreOne(ctx context.Context, user *User.User) (*User.User, bool, error)
	FindByEmail(ctx context.Context, email string) (*User.User, error)
	GenerateAndStoreToken(ctx context.Context, userId primitive.ObjectID, email string) (string, error)
	FindUserById(ctx context.Context, userId string) (*User.User, error)
	PutUser(ctx context.Context, user *User.User) (*User.User, error)
	PutUserPassword(ctx context.Context, user *User.UpdatePasswordPayload) error
}

type UserUsecase interface {
	GetAllData(ctx context.Context) ([]User.User, error)
	RegisterUser(ctx context.Context, user *User.User) (*User.User, []string)
	LoginUser(ctx context.Context, userReq *User.LoginPayload) (*User.LoginResponse, error)
	SendEmailVerification(ctx context.Context, data *User.User, secretCode string) (res *User.User, err error)
	GetUserByUserId(ctx context.Context, token string) (*User.User, error)
	UpdateUserData(ctx context.Context, user *User.UpdateProfilePayload) (*User.User, []string)
	UpdatePassword(ctx context.Context, user *User.UpdatePasswordPayload) error
}
