package http

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"main.go/src/configs/app"
	_const "main.go/src/configs/const"
	"main.go/src/configs/database"
	Shelter "main.go/src/domains/shelter/entities"
	"main.go/src/domains/shelter/interfaces"
	"main.go/src/domains/shelter/interfaces/impl/repository"
	"main.go/src/domains/shelter/interfaces/impl/usecase"
	"main.go/src/middlewares"
	"main.go/src/shared/helpers"
	"main.go/src/shared/helpers/image_helpers"
	"main.go/src/shared/message/errors"
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
	user := router.Group(guest.BasePath(), middlewares.JwtAuthMiddleware(app.GetConfig().AccessToken.AccessTokenSecret, "user", "admin"))
	{
		user.GET("/:id", handler.FindOneById)
		user.GET("/my-shelter", handler.FindOneByUserId)
		user.POST("/register", handler.RegisterShelter)
		user.GET("/favorite", handler.FindAllFavorite)
		user.PUT("/update", handler.UpdateShelter)
	}
	admin := router.Group("/admin"+guest.BasePath(), middlewares.JwtAuthMiddleware(app.GetConfig().AccessToken.AccessTokenSecret, "user", "admin"))
	{
		admin.PUT("/update/:id", handler.UpdateShelterByAdmin)
		admin.Group("/delete/")
	}
	return handler
}

func (shelterHttp *ShelterHttp) FindAllFavorite(c *gin.Context) {
	search := &Shelter.ShelterSearch{
		Search:   c.Query("search"),
		Page:     helpers.ParseStringToInt(c.Query("page")),
		PageSize: helpers.ParseStringToInt(c.Query("page_size")),
		Sort:     c.Query("sort"),
		OrderBy:  c.Query("order_by"),
		UserId:   helpers.GetUserId(c),
	}
	data, err := shelterHttp.shelterUsecase.GetAllData(c, search)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Get Shelter Data ! ", Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (shelterHttp *ShelterHttp) FindAll(c *gin.Context) {
	search := &Shelter.ShelterSearch{
		Search:              c.Query("search"),
		Page:                helpers.ParseStringToInt(c.Query("page")),
		PageSize:            helpers.ParseStringToInt(c.Query("page_size")),
		Sort:                c.Query("sort"),
		ShelterLocationName: c.Query("location_name"),
		OrderBy:             c.Query("order_by"),
	}
	data, err := shelterHttp.shelterUsecase.GetAllData(c, search)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Get Shelter Data ! ", Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (shelterHttp *ShelterHttp) RegisterShelter(c *gin.Context) {
	shelterCreate := &Shelter.ShelterCreate{}
	// Parse the multipart form with a maximum of 30 MB memory
	if err := c.Request.ParseMultipartForm(30 << 20); err != nil { // 30 MB max memory
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Parse MultiPartForm Request ! ", Error: err.Error()})
		return
	}
	form, _ := c.MultipartForm()
	jsonData := form.Value["data"][0]
	shelter := &Shelter.Shelter{}
	if err := json.Unmarshal([]byte(jsonData), &shelter); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Marshal Request ! ", Error: err.Error()})
		return
	}
	shelter.UserId = helpers.GetUserId(c)
	tempFilePaths, errs := image_helpers.SaveImageToTemp(c, form)
	if errs != nil {
		// Delete temporary files if pet creation fails
		for _, tempFilePath := range tempFilePaths {
			_ = os.Remove(tempFilePath)
		}
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Move Image ! ", Error: errs.Error()})
	}
	res, err := shelterHttp.shelterUsecase.RegisterShelter(c, shelter)
	if err != nil {
		if res != nil {
			c.JSON(http.StatusOK, errors.ErrorWrapper{ErrorS: err, Data: res})
			return
		}
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Register Shelter ! ", ErrorS: err})
		return
	}
	if tempFilePaths != nil {
		shelterCreate, _ = image_helpers.MoveUploadedShelterFile(c, tempFilePaths, shelter, shelterCreate, res.ID.Hex())
	}
	shelter.Image = shelterCreate.Shelter.Image
	shelter, _ = shelterHttp.shelterUsecase.UpdateShelterById(c, &shelter.ID, shelter)
	c.JSON(http.StatusOK, errors.SuccessWrapper{Message: "Success Register Shelter", Data: res})
}

func (shelterHttp *ShelterHttp) FindOneByUserId(c *gin.Context) {
	search := &Shelter.ShelterSearch{
		UserId: helpers.GetUserId(c),
	}
	data, err := shelterHttp.shelterUsecase.GetOneDataByUserId(c, search)
	if err != nil {
		c.JSON(http.StatusOK, errors.ErrorWrapper{Message: "Failed To Get Shelter ! ", Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, errors.SuccessWrapper{Message: "Success Get Shelter By User Id ! ", Data: data})
}

func (shelterHttp *ShelterHttp) FindOneByUserIdForRequest(c *gin.Context) (*Shelter.ShelterResponsePayload, error) {
	search := &Shelter.ShelterSearch{
		UserId: helpers.GetUserId(c),
	}
	data, err := shelterHttp.shelterUsecase.GetOneDataByUserId(c, search)
	if err != nil {
		return nil, err
	}
	return data, nil
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

func (shelterHttp *ShelterHttp) UpdateShelter(c *gin.Context) {
	shelterReq := &Shelter.ShelterUpdate{}
	// Parse the multipart form with a maximum of 30 MB memory
	if err := c.Request.ParseMultipartForm(30 << 20); err != nil { // 30 MB max memory
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Parse MultiPartForm Request ! ", Error: err.Error()})
		return
	}
	form, _ := c.MultipartForm()
	jsonData := form.Value["data"][0]
	shelter := &Shelter.Shelter{}
	if err := json.Unmarshal([]byte(jsonData), &shelter); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Marshal Shelter Update Request ! ", Error: err.Error()})
		return
	}
	shelter.UserId = helpers.GetUserId(c)
	shelterReq.Shelter = shelter
	findShelter, err := shelterHttp.shelterUsecase.GetOneDataByUserId(c, &Shelter.ShelterSearch{
		UserId: shelter.UserId,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed to Get Shelter Data ! ", Error: err.Error()})
		return
	}
	if helpers.GetRoleFromContext(c) == "user" {
		if findShelter.UserId != shelter.UserId {
			c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "You Can Only Update Your Own Shelter"})
			return
		}
	}
	shelterReq.Shelter.ID = findShelter.ID
	shelterReq.Shelter.Image = findShelter.Image
	if form.File != nil {
		shelterReq, _ = image_helpers.UploadShelter(c, form, shelterReq)
	}
	res, err := shelterHttp.shelterUsecase.UpdateShelterById(c, &findShelter.ID, shelterReq.Shelter)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed to Update Shelter ! ", Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, errors.SuccessWrapper{Message: "Shelter updated successfully ! ", Data: res})
}

func (shelterHttp *ShelterHttp) FindOneByShelterId(c *gin.Context, Id primitive.ObjectID) primitive.ObjectID {
	search := &Shelter.ShelterSearch{
		ShelterId: Id,
	}
	data, err := shelterHttp.shelterUsecase.GetOneDataByIdForRequest(c, search)
	if err != nil {
		logrus.Warnf("Failed to get data shelter for request %v: %v", search.ShelterId, err)
	}
	return data.UserId
}

func (shelterHttp *ShelterHttp) UpdateShelterByAdmin(c *gin.Context) {
	shelterReq := &Shelter.ShelterUpdate{}
	// Parse the multipart form with a maximum of 30 MB memory
	if err := c.Request.ParseMultipartForm(30 << 20); err != nil { // 30 MB max memory
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Parse MultiPartForm Request ! ", Error: err.Error()})
		return
	}
	form, _ := c.MultipartForm()
	jsonData := form.Value["data"][0]
	shelter := &Shelter.Shelter{}
	if err := json.Unmarshal([]byte(jsonData), &shelter); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Marshal Shelter Update Request ! ", Error: err.Error()})
		return
	}
	shelter.ID = helpers.ParseStringToObjectId(c.Param("id"))
	shelterReq.Shelter = shelter
	findShelter, err := shelterHttp.shelterUsecase.GetOneDataByIdForRequest(c, &Shelter.ShelterSearch{
		ShelterId: shelter.ID,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed to Get Shelter Data ! ", Error: err.Error()})
		return
	}
	shelterReq.Shelter.UserId = findShelter.UserId
	if form.File != nil {
		shelterReq, _ = image_helpers.UploadShelter(c, form, shelterReq)
	}
	res, err := shelterHttp.shelterUsecase.UpdateShelterById(c, &findShelter.ID, shelterReq.Shelter)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed to Update Shelter ! ", Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, errors.SuccessWrapper{Message: "Shelter updated successfully ! ", Data: res})
}
