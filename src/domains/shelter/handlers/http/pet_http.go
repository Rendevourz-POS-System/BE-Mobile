package http

import (
	"encoding/json"
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
	"main.go/shared/helpers/image_helpers"
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
	// Parse the multipart form with a maximum of 30 MB memory
	if err := ctx.Request.ParseMultipartForm(30 << 20); err != nil { // 30 MB max memory
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Parse MultiPartForm Request ! ", Error: err.Error()})
		return
	}
	// Get the multipart form data
	form, _ := ctx.MultipartForm()
	// Unmarshal the JSON data into the Pet struct
	jsonData := form.Value["data"][0]
	if err := json.Unmarshal([]byte(jsonData), &pet.Pet); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Bind JSON Request ! ", Error: err.Error()})
		return
	}
	filesName, err := image_helpers.SaveImageToTemp(ctx, form)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Move Image ! ", Error: err.Error()})
	}
	// Create the pet data
	data, errs := h.petUsecase.CreatePets(ctx, pet)
	if errs != nil {
		// Delete temporary files if pet creation fails
		image_helpers.RemoveTempImagePath(filesName)
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Create Pet ! ", ErrorS: errs})
		return
	}
	if filesName != nil {
		if pet.Pet.ShelterId.Hex() == "" {
			pet, err = image_helpers.MoveUploadedFile(ctx, filesName, data, pet, app.GetConfig().Image.PetPath)
		} else {
			pet, err = image_helpers.MoveUploadedFile(ctx, filesName, data, pet, app.GetConfig().Image.UserPath, app.GetConfig().Image.PetPath, helpers.GetUserId(ctx).Hex())
		}
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Move Image Pet ! ", Error: err.Error()})
	}
	data.Image = pet.Pet.Image
	// Update the pet entity with the image paths
	_, err = h.petUsecase.UpdatePet(ctx, &data.ID, data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Update Pet with Image Paths ! ", Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, data)
}

func (h *PetHttp) GetAllPets(ctx *gin.Context) {
	search := &Pet.PetSearch{
		Search:      ctx.Query("search"),
		Page:        helpers.ParseStringToInt(ctx.Query("page")),
		PageSize:    helpers.ParseStringToInt(ctx.Query("page_size")),
		Gender:      helpers.CheckPetGender(ctx.Query("gender")),
		Type:        ctx.Query("type"),
		Sort:        ctx.Query("sort"),
		OrderBy:     ctx.Query("order_by"),
		Location:    ctx.Query("location"),
		AgeStart:    helpers.ParseStringToInt(ctx.Query("age_start")),
		AgeEnd:      helpers.ParseStringToInt(ctx.Query("age_end")),
		ShelterName: ctx.Query("shelter_name"),
	}
	data, err := h.petUsecase.GetAllPets(ctx, search)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, data)
}
