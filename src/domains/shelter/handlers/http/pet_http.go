package http

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"main.go/configs/app"
	_const "main.go/configs/const"
	"main.go/configs/database"
	Pet "main.go/domains/shelter/entities"
	"main.go/domains/shelter/interfaces"
	"main.go/domains/shelter/interfaces/impl/repository"
	"main.go/domains/shelter/interfaces/impl/usecase"
	"main.go/middlewares"
	"main.go/shared/helpers"
	"main.go/shared/message/errors"
	"net/http"
)

type PetHttp struct {
	petUsecase interfaces.PetUseCase
}

func NewPetHttp(router *gin.Engine) *PetHttp {
	handler := &PetHttp{
		petUsecase: usecase.NewPetUseCase(repository.NewPetRepository(database.GetDatabase(_const.DB_SHELTER_APP))),
	}
	guest := router.Group("/pet")
	{
		guest.GET("", handler.GetAllPets)
	}
	user := router.Group("/pet", middlewares.JwtAuthMiddleware(app.GetConfig().AccessToken.AccessTokenSecret))
	{
		user.POST("/create", handler.CreatePet)
	}
	return handler
}

func (h *PetHttp) CreatePet(ctx *gin.Context) {
	pet := &Pet.PetCreate{}
	if err := ctx.Request.ParseMultipartForm(30 << 40); err != nil { // 30 MB max memory
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Parse MultiPartForm Request ! ", Error: err.Error()})
		return
	}
	form, _ := ctx.MultipartForm()
	jsonData := form.Value["data"]
	if err := json.Unmarshal([]byte(jsonData[0]), &pet.Pet); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Bind JSON Request ! ", Error: err.Error()})
		return
	}
	data, err := h.petUsecase.CreatePets(ctx, pet)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Create Pet ! ", ErrorS: err})
		return
	}
	files := form.File["files"]
	for _, file := range files {
		err := ctx.SaveUploadedFile(file, fmt.Sprintf("uploads/%s/pets/%s", data.ID, file.Filename))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Save File ! ", Errors: err})
			return
		} else {
			pet.Pet.ImagePath = append(pet.Pet.ImagePath, file.Filename)
		}
	}
	ctx.JSON(http.StatusOK, data)
}

func (h *PetHttp) GetAllPets(ctx *gin.Context) {
	search := &Pet.PetSearch{
		Search:   ctx.Query("search"),
		Page:     helpers.ParseStringToInt(ctx.Query("page")),
		PageSize: helpers.ParseStringToInt(ctx.Query("page_size")),
		Gender:   helpers.CheckPetGender(ctx.Query("gender")),
		Type:     ctx.Query("type"),
		Sort:     ctx.Query("sort"),
		OrderBy:  ctx.Query("order_by"),
		Location: ctx.Query("location"),
		AgeStart: helpers.ParseStringToInt(ctx.Query("age_start")),
		AgeEnd:   helpers.ParseStringToInt(ctx.Query("age_end")),
	}
	data, err := h.petUsecase.GetAllPets(ctx, search)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, data)
}
