package interfaces

type UserTokenUsecase interface {
	GenerateToken() (string, error)
}

type UserTokenRepository interface {
	StoreToken(token string) error
}
