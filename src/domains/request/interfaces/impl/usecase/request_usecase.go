package usecase

import "main.go/domains/request/interfaces"

type requestUsecase struct {
	requestRepo interfaces.RequestRepository
}

func NewRequestUsecase(requestRepo interfaces.RequestRepository) *requestUsecase {
	return &requestUsecase{requestRepo}
}
