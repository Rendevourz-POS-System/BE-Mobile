package interfaces

import (
	"context"
	User "main.go/domains/user/entities"
)

type UserUsecase interface {
	GetAllData(ctx context.Context) ([]User.User, error)
	RegisterUser(ctx context.Context, user *User.User) (*User.User, []string)
	LoginUser(ctx context.Context, userReq *User.LoginPayload) (*User.LoginResponse, error)
	SendEmailVerification(ctx context.Context, data *User.User) (res *User.User, err error)
}

type UserRepository interface {
	FindAll(ctx context.Context) ([]User.User, error)
	StoreOne(ctx context.Context, user *User.User) (*User.User, bool, error)
	FindByEmail(ctx context.Context, email string) (*User.User, error)
}
