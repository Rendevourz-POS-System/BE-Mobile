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
	"strings"
)

type PetHttp struct {
	petUsecase  interfaces.PetUseCase
	shelterHttp *ShelterHttp
}

func NewPetHttp(router *gin.Engine, shelterHttp *ShelterHttp) *PetHttp {
	handler := &PetHttp{
		petUsecase:  usecase.NewPetUseCase(repository.NewPetRepository(database.GetDatabase(_const.DB_SHELTER_APP))),
		shelterHttp: shelterHttp,
	}
	guest := router.Group("/pet")
	{
		guest.GET("", handler.GetAllPets)
		guest.GET("/:id", handler.FindPetById)
	}
	user := router.Group(guest.BasePath(), middlewares.JwtAuthMiddleware(app.GetConfig().AccessToken.AccessTokenSecret, "user", "admin"))
	{
		user.POST("/create", handler.CreatePet)
		user.GET("/favorite", handler.FindAllFavorite)
		user.DELETE("/delete", handler.DeletePetByUser)
		user.PUT("/update", handler.UpdatePet)
	}
	admin := router.Group("/admin"+guest.BasePath(), middlewares.JwtAuthMiddleware(app.GetConfig().AccessToken.AccessTokenSecret, "admin"))
	{
		admin.DELETE("/delete/:id", handler.DeletePetByAdmin)
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
		go image_helpers.RemoveTempImagePath(filesName)
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Create Pet ! ", ErrorS: errs})
		return
	}
	if filesName != nil {
		if pet.Pet.ShelterId == nil {
			pet, err = image_helpers.MoveUploadedFile(ctx, filesName, pet, app.GetConfig().Image.PetPath, data.ID.Hex())
		} else {
			pet, err = image_helpers.MoveUploadedFile(ctx, filesName, pet, app.GetConfig().Image.UserPath, app.GetConfig().Image.ShelterPath, pet.Pet.ShelterId.Hex(), app.GetConfig().Image.PetPath, data.ID.Hex())
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
		ShelterId:   ctx.Query("shelter_id"),
	}
	data, err := h.petUsecase.GetAllPets(ctx, search)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, data)
}

func (h *PetHttp) FindAllFavorite(c *gin.Context) {
	search := &Pet.PetSearch{
		Search:   c.Query("search"),
		Page:     helpers.ParseStringToInt(c.Query("page")),
		PageSize: helpers.ParseStringToInt(c.Query("page_size")),
		Sort:     c.Query("sort"),
		OrderBy:  c.Query("order_by"),
		UserId:   helpers.GetUserId(c),
	}
	data, err := h.petUsecase.GetAllPets(c, search)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Get Shelter Data ! ", Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *PetHttp) FindPetById(ctx *gin.Context) {
	Id := helpers.ParseStringToObjectId(ctx.Param("id"))
	data, err := h.petUsecase.GetPetById(ctx, &Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed to find pet !", Errors: err})
		return
	}
	ctx.JSON(http.StatusOK, errors.SuccessWrapper{Data: data, Message: "Success Get Pet Detail ! "})
}

func (h *PetHttp) UpdatePet(ctx *gin.Context) {
	pet := &Pet.PetUpdatePayload{}
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
	findPet, err := h.petUsecase.GetPetById(ctx, &pet.Pet.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Pet Id Not Found ! ", Error: err.Error()})
		return
	}
	if strings.ToLower(helpers.GetRoleFromContext(ctx)) == "user" {
		if findPet.ShelterId != nil && len(findPet.ShelterId) > 0 {
			data, _ := h.shelterHttp.shelterUsecase.GetOneDataById(ctx, &Pet.ShelterSearch{
				ShelterId: *findPet.ShelterId,
			})
			if data.UserId != helpers.GetUserId(ctx) {
				ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "You Can Only Update Your Own Pet"})
				return
			}
		}
	}
	if form.File != nil {
		pet, _ = image_helpers.UploadPet(ctx, form, pet)
	}
	updatedPet, errUpdatePet := h.petUsecase.UpdatePet(ctx, &pet.Pet.ID, &pet.Pet)
	if errUpdatePet != nil {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Update Pet with Image Paths ! ", Error: errUpdatePet.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updatedPet)
}

func (h *PetHttp) DeletePetByAdmin(ctx *gin.Context) {
	Id := helpers.ParseStringToObjectId(ctx.Param("id"))
	data, err := h.petUsecase.DeletePetByAdmin(ctx, &Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed to find pet !", Errors: err})
		return
	}
	ctx.JSON(http.StatusOK, errors.SuccessWrapper{Data: data, Message: "Success Get Pet Detail ! "})
}

func (h *PetHttp) DeletePetByUser(ctx *gin.Context) {
	userDeleteReq := Pet.PetDeletePayload{}
	if err := ctx.ShouldBindJSON(&userDeleteReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Bad Request !", ErrorS: []string{err.Error()}})
		return
	}
	userDeleteReq.UserId = helpers.GetUserId(ctx)
	data, err := h.petUsecase.DeletePetByUser(ctx, userDeleteReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed to find pet !", ErrorS: err})
		return
	}
	ctx.JSON(http.StatusOK, errors.SuccessWrapper{Data: data, Message: "Success Get Pet Detail ! "})
}
