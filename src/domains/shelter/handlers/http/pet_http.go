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
	"main.go/shared/message/errors"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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
	jsonData := form.Value["data"]
	if err := json.Unmarshal([]byte(jsonData[0]), &pet.Pet); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Bind JSON Request ! ", Error: err.Error()})
		return
	}
	tempFilePaths, err := h.saveImageToTemp(ctx, form)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Move Image ! ", Error: err.Error()})
	}
	// Create the pet data
	data, errs := h.petUsecase.CreatePets(ctx, pet)
	if errs != nil {
		// Delete temporary files if pet creation fails
		for _, tempFilePath := range tempFilePaths {
			_ = os.Remove(tempFilePath)
		}
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Create Pet ! ", ErrorS: errs})
		return
	}
	pet, err = h.moveUploadedFile(ctx, tempFilePaths, data, pet)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Move Image Pet ! ", Error: err.Error()})
	}
	data.ImagePath = pet.Pet.ImagePath
	// Update the pet entity with the image paths
	_, err = h.petUsecase.UpdatePet(ctx, &data.ID, data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Update Pet with Image Paths ! ", Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, data)
}

func (h *PetHttp) saveImageToTemp(ctx *gin.Context, form *multipart.Form) (res []string, err error) {
	// Get the uploaded files
	files := form.File["files"]

	// Save the uploaded files to a temporary directory
	tempFilePaths := make([]string, len(files))
	for i, file := range files {
		// Construct the temporary file path
		tempFilePath := filepath.Join("uploads", "temp", file.Filename)

		// Save the uploaded file with the temporary path
		if err := ctx.SaveUploadedFile(file, tempFilePath); err != nil {
			return nil, err
		}
		tempFilePaths[i] = tempFilePath
	}
	return tempFilePaths, nil
}

func (h *PetHttp) moveUploadedFile(ctx *gin.Context, tempFilePaths []string, data *Pet.Pet, pet *Pet.PetCreate) (res *Pet.PetCreate, err error) {
	// Move the uploaded files to their final location with the data.ID in the path
	for _, tempFilePath := range tempFilePaths {
		// Construct the final file path
		finalFilePath := filepath.Join("uploads", data.ID.Hex(), "pets", filepath.Base(tempFilePath))

		// Create directories if they don't exist
		if err = os.MkdirAll(filepath.Dir(finalFilePath), 0755); err != nil {
			ctx.JSON(http.StatusInternalServerError, errors.ErrorWrapper{Message: "Failed To Create Directories ! ", Errors: err})
			return
		}
		// Move the temporary file to the final location
		if err = os.Rename(tempFilePath, finalFilePath); err != nil {
			ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Move File ! ", Errors: err})
			return
		}
		// Update the pet's image path
		pet.Pet.ImagePath = append(pet.Pet.ImagePath, finalFilePath)
	}
	return pet, nil
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
