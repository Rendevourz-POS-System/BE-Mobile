package image_helpers

import (
	"github.com/gin-gonic/gin"
	"main.go/configs/app"
	Pet "main.go/domains/shelter/entities"
	"main.go/shared/helpers"
	"main.go/shared/message/errors"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func SaveImageToTemp(ctx *gin.Context, form *multipart.Form) (res []string, err error) {
	// Get the uploaded files
	files := form.File["files"]

	// Save the uploaded files to a temporary directory
	tempFilePaths := make([]string, len(files))
	for i, file := range files {
		// Construct the temporary file path
		tempFilePath := filepath.Join(app.GetConfig().Image.Folder, app.GetConfig().Image.TempPath, file.Filename)

		// Save the uploaded file with the temporary path
		if err := ctx.SaveUploadedFile(file, tempFilePath); err != nil {
			return nil, err
		}
		tempFilePaths[i] = tempFilePath
	}
	return tempFilePaths, nil
}

func MoveUploadedFile(ctx *gin.Context, tempFilePaths []string, data *Pet.Pet, pet *Pet.PetCreate, path string) (res *Pet.PetCreate, err error) {
	// Move the uploaded files to their final location with the data.ID in the path
	for _, tempFilePath := range tempFilePaths {
		// Construct the final file path
		finalFilePath := filepath.Join(app.GetConfig().Image.Folder, helpers.ToString(path), data.ID.Hex(), filepath.Base(tempFilePath))

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

func MoveUploadedShelterFile(ctx *gin.Context, tempFilePaths []string, data *Pet.Shelter, shelter *Pet.ShelterCreate, path string) (res *Pet.ShelterCreate, err error) {
	// Move the uploaded files to their final location with the data.ID in the path
	for _, tempFilePath := range tempFilePaths {
		// Construct the final file path
		finalFilePath := filepath.Join(app.GetConfig().Image.Folder, app.GetConfig().Image.UserPath, app.GetConfig().Image.ShelterPath, data.ID.Hex(), filepath.Base(tempFilePath))

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
		shelter.Shelter.ImagePath = append(shelter.Shelter.ImagePath, finalFilePath)
	}
	return shelter, nil
}
