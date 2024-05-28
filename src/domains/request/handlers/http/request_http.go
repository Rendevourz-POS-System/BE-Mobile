package http

import (
	"github.com/gin-gonic/gin"
	_const "main.go/configs/const"
	"main.go/configs/database"
	"main.go/domains/request/interfaces"
	"main.go/domains/request/interfaces/impl/repository"
	"main.go/domains/request/interfaces/impl/usecase"
)

type RequestHttp struct {
	requestUsecase interfaces.RequestUsecase
}

func NewRequestHttp(rotuer *gin.Engine) *RequestHttp {
	handlers := &RequestHttp{
		requestUsecase: usecase.NewRequestUsecase(repository.NewRequestRepository(database.GetDatabase(_const.DB_SHELTER_APP))),
	}
	return handlers
}
