package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main.go/configs/app"
	_const "main.go/configs/const"
	"main.go/configs/database"
	Shelter "main.go/domains/shelter/entities"
	"main.go/domains/shelter/interfaces"
	"main.go/domains/shelter/interfaces/impl/repository"
	"main.go/domains/shelter/interfaces/impl/usecase"
	"main.go/middlewares"
	"main.go/shared/helpers"
	"main.go/shared/message/errors"
	"net/http"
)

type ShelterHttp struct {
	shelterUsecase interfaces.ShelterUsecase
}

func NewShelterHttp(router *gin.Engine) *ShelterHttp {
	handler := &ShelterHttp{
		shelterUsecase: usecase.NewShelterUsecase(repository.
			NewShelterRepository(database.GetDatabase(_const.DB_SHELTER_APP))),
	}
	guest := router.Group("/shelter")
	{
		guest.GET("", handler.FindAll)
	}
	user := router.Group("/shelter", middlewares.JwtAuthMiddleware(app.GetConfig().AccessToken.AccessTokenSecret))
	{
		user.GET("/my-shelter", handler.FindOneByUserId)
		user.GET("/:id", handler.FindOneById)
		user.POST("/register", handler.RegisterShelter)
	}
	return handler
}

func (shelterHttp *ShelterHttp) FindAll(c *gin.Context) {
	search := &Shelter.ShelterSearch{
		Search:   c.Query("search"),
		Page:     helpers.ParseStringToInt(c.Query("page")),
		PageSize: helpers.ParseStringToInt(c.Query("page_size")),
		Sort:     c.Query("sort"),
		OrderBy:  c.Query("order_by"),
	}
	data, err := shelterHttp.shelterUsecase.GetAllData(c, search)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Get Shelter Data ! ", Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (shelterHttp *ShelterHttp) RegisterShelter(c *gin.Context) {
	shelter := &Shelter.Shelter{}
	if err := c.ShouldBindJSON(&shelter); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Bind JSON Request ! ", Error: err.Error()})
		return
	}
	shelter.UserId = helpers.GetUserId(c)
	res, err := shelterHttp.shelterUsecase.RegisterShelter(c, shelter)
	if err != nil {
		if res != nil {
			c.JSON(http.StatusOK, errors.ErrorWrapper{Message: err.Error(), Data: res})
			return
		}
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Register Shelter ! ", Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, errors.SuccessWrapper{Message: "Success Register Shelter", Data: res})
}

func (shelterHttp *ShelterHttp) FindOneByUserId(c *gin.Context) {
	search := &Shelter.ShelterSearch{
		ShelterId: helpers.GetUserId(c),
	}
	data, err := shelterHttp.shelterUsecase.GetOneDataByUserId(c, search)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Get Shelter ! ", Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, errors.SuccessWrapper{Message: "Success Get Shelter By User Id ! ", Data: data})
}

func (shelterHttp *ShelterHttp) FindOneById(c *gin.Context) {
	fmt.Println("ID : ", c.Param("id"))
	search := &Shelter.ShelterSearch{
		ShelterId: helpers.ParseStringToObjectId(c.Param("id")),
	}
	data, err := shelterHttp.shelterUsecase.GetOneDataById(c, search)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Get Shelter ! ", Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, errors.SuccessWrapper{Message: "Success Get Shelter By Id ! ", Data: data})
}
