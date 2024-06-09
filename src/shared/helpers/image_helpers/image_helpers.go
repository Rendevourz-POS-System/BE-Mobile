package image_helpers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"main.go/configs/app"
	Pet "main.go/domains/shelter/entities"
	"main.go/domains/user/entities"
	"main.go/shared/message/errors"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func GenerateImagePath(data ...string) string {
	return filepath.Join(app.GetConfig().Image.Folder, filepath.Join(data...))
}

func RemoveTempImagePath(fileName []string) {
	for _, tempFilePath := range fileName {
		err := os.Remove(GenerateImagePath(app.GetConfig().Image.TempPath, tempFilePath))
		fmt.Println(err)
	}
}

func SaveImageToTemp(ctx *gin.Context, form *multipart.Form) (res []string, err error) {
	// Get the uploaded files
	files := form.File["files"]
	if len(files) <= 0 {
		return nil, nil
	}
	// Save the uploaded files to a temporary directory
	tempFilePaths := make([]string, len(files))
	for i, file := range files {
		// Construct the temporary file path
		tempFilePath := GenerateImagePath(app.GetConfig().Image.TempPath, file.Filename)

		// Save the uploaded file with the temporary path
		if err := ctx.SaveUploadedFile(file, tempFilePath); err != nil {
			return nil, err
		}
		tempFilePaths[i] = file.Filename
	}
	return tempFilePaths, nil
}

func MoveUploadedShelterFile(ctx *gin.Context, filesName []string, shelter *Pet.Shelter, shelterCreate *Pet.ShelterCreate, shelterId string) (res *Pet.ShelterCreate, err error) {
	// Move the uploaded files to their final location with the data.ID in the path
	for _, RealFileName := range filesName {
		// Construct the final file path
		tempPath := GenerateImagePath(app.GetConfig().Image.TempPath, RealFileName)
		finalFilePath := GenerateImagePath(app.GetConfig().Image.UserPath, app.GetConfig().Image.ShelterPath, shelterId, RealFileName)

		// Create directories if they don't exist
		if err = os.MkdirAll(filepath.Dir(finalFilePath), 0755); err != nil {
			ctx.JSON(http.StatusInternalServerError, errors.ErrorWrapper{Message: "Failed To Create Directories ! ", Errors: err})
			return
		}
		// Move the temporary file to the final location
		if err = os.Rename(tempPath, finalFilePath); err != nil {
			ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Move File ! ", Errors: err})
			return
		}
		// Update the pet's image path
		shelterCreate.Shelter.Image = append(shelterCreate.Shelter.Image, RealFileName)
	}
	return shelterCreate, nil
}

func UploadProfile(ctx *gin.Context, file *multipart.FileHeader, data *entities.UpdateProfilePayload) (res *entities.UpdateProfilePayload, err error) {
	FilePath := GenerateImagePath(app.GetConfig().Image.UserPath, app.GetConfig().Image.ProfilePath, data.ID.Hex(), file.Filename)
	if data.OldImageName != "" {
		OldFilePath := GenerateImagePath(app.GetConfig().Image.UserPath, app.GetConfig().Image.ProfilePath, data.ID.Hex(), data.OldImageName)
		// Check if a file already exists at the FilePath
		if _, err = os.Stat(OldFilePath); err == nil {
			// File exists, attempt to remove it
			if err = os.Remove(OldFilePath); err != nil {
				log.Printf("Failed to remove file: %s, error: %v", OldFilePath, err) // More detailed logging
				ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed to Remove Existing File!", Error: err.Error()})
				return nil, err
			}
		} else if !os.IsNotExist(err) {
			// An error other than "not existing" occurred
			log.Printf("Error checking file: %s, error: %v", OldFilePath, err) // Log unexpected errors
			ctx.JSON(http.StatusInternalServerError, errors.ErrorWrapper{Message: "Error Checking Existing File!", Error: err.Error()})
			return nil, err
		}
	}
	// Checking the directory if it doesn't exist
	if err = os.MkdirAll(filepath.Dir(FilePath), 0755); err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.ErrorWrapper{Message: "Failed To Create Directories ! ", Errors: err})
		return
	}
	// Save the uploaded file with the temporary path
	if err = ctx.SaveUploadedFile(file, FilePath); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed to Upload Image !", Error: err.Error()})
		return
	}
	data.Image = file.Filename
	return data, nil
}

func UploadShelter(ctx *gin.Context, form *multipart.Form, shelterUpdate *Pet.ShelterUpdate) (res *Pet.ShelterUpdate, err error) {
	// Get the uploaded files
	files := form.File["files"]
	if shelterUpdate.Shelter.OldImage != nil && len(shelterUpdate.Shelter.OldImage) > 0 {
		// Move the uploaded files to their final location with the data.ID in the path
		for _, RealFileName := range shelterUpdate.Shelter.OldImage {
			// Construct the final file path
			OldFilePath := GenerateImagePath(app.GetConfig().Image.UserPath, shelterUpdate.Shelter.ID.Hex(), app.GetConfig().Image.ShelterPath, RealFileName)
			// Check if a file already exists at the FilePath
			if _, err = os.Stat(OldFilePath); err == nil {
				// File exists, attempt to remove it
				if err = os.Remove(OldFilePath); err != nil {
					log.Printf("Failed to remove file: %s, error: %v", OldFilePath, err) // More detailed logging
					ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed to Remove Existing File!", Error: err.Error()})
					return nil, err
				}
			} else if !os.IsNotExist(err) {
				// An error other than "not existing" occurred
				log.Printf("Error checking file: %s, error: %v", OldFilePath, err) // Log unexpected errors
				ctx.JSON(http.StatusInternalServerError, errors.ErrorWrapper{Message: "Error Checking Existing File!", Error: err.Error()})
				return nil, err
			}
		}
	}
	if files == nil || len(files) <= 0 {
		return shelterUpdate, nil
	}
	for _, File := range files {
		FilePath := GenerateImagePath(app.GetConfig().Image.UserPath, shelterUpdate.Shelter.ID.Hex(), app.GetConfig().Image.ShelterPath, File.Filename)
		// Create directories if they don't exist
		if err = os.MkdirAll(filepath.Dir(FilePath), 0755); err != nil {
			ctx.JSON(http.StatusInternalServerError, errors.ErrorWrapper{Message: "Failed To Create Directories ! ", Errors: err})
			return
		}
		// Save the uploaded file with the temporary path
		if err = ctx.SaveUploadedFile(File, FilePath); err != nil {
			ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed to Upload Image !", Error: err.Error()})
			return
		}
		shelterUpdate.Shelter.Image = append(shelterUpdate.Shelter.Image, File.Filename)
	}
	return shelterUpdate, nil
}

func MoveUploadedFile(ctx *gin.Context, filesName []string, pet *Pet.PetCreate, data ...string) (res *Pet.PetCreate, err error) {
	// Move the uploaded files to their final location with the data.ID in the path
	for _, RealFileName := range filesName {
		// Construct the final file path
		newData := data
		newData = append(newData, RealFileName)
		finalFilePath := GenerateImagePath(newData...)
		// Create directories if they don't exist
		if err = os.MkdirAll(filepath.Dir(finalFilePath), 0755); err != nil {
			ctx.JSON(http.StatusInternalServerError, errors.ErrorWrapper{Message: "Failed To Create Directories ! ", Errors: err})
			return
		}
		// Move the temporary file to the final location
		if err = os.Rename(GenerateImagePath(app.GetConfig().Image.TempPath, RealFileName), finalFilePath); err != nil {
			ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed To Move File ! ", Errors: err})
			return
		}
		// Update the pet's image path
		pet.Pet.Image = append(pet.Pet.Image, RealFileName)
	}
	return pet, nil
}

func UploadPet(ctx *gin.Context, form *multipart.Form, petUpdate *Pet.PetUpdatePayload) (res *Pet.PetUpdatePayload, err error) {
	files := form.File["files"]
	if petUpdate.Pet.OldImage != nil && len(petUpdate.Pet.OldImage) > 0 {
		// Move the uploaded files to their final location with the data.ID in the path
		for _, RealFileName := range petUpdate.Pet.OldImage {
			// Construct the final file path
			var OldFilePath string
			if petUpdate.Pet.ShelterId == nil {
				OldFilePath = GenerateImagePath(app.GetConfig().Image.PetPath, petUpdate.Pet.ID.Hex(), RealFileName)
			} else {
				OldFilePath = GenerateImagePath(app.GetConfig().Image.UserPath, app.GetConfig().Image.ShelterPath, petUpdate.Pet.ShelterId.Hex(), app.GetConfig().Image.PetPath, petUpdate.Pet.ID.Hex(), RealFileName)
			}
			// Check if a file already exists at the FilePath
			if _, err = os.Stat(OldFilePath); err == nil {
				// File exists, attempt to remove it
				if err = os.Remove(OldFilePath); err != nil {
					log.Printf("Failed to remove file: %s, error: %v", OldFilePath, err) // More detailed logging
					ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed to Remove Existing File!", Error: err.Error()})
					return nil, err
				}
			} else if !os.IsNotExist(err) {
				// An error other than "not existing" occurred
				log.Printf("Error checking file: %s, error: %v", OldFilePath, err) // Log unexpected errors
				ctx.JSON(http.StatusInternalServerError, errors.ErrorWrapper{Message: "Error Checking Existing File!", Error: err.Error()})
				return nil, err
			}
		}
	}
	if files == nil || len(files) <= 0 {
		return petUpdate, nil
	}
	for _, File := range files {
		var FilePath string
		if petUpdate.Pet.ShelterId == nil {
			FilePath = GenerateImagePath(app.GetConfig().Image.PetPath, petUpdate.Pet.ID.Hex(), File.Filename)
		} else {
			FilePath = GenerateImagePath(app.GetConfig().Image.UserPath, app.GetConfig().Image.ShelterPath, petUpdate.Pet.ShelterId.Hex(), app.GetConfig().Image.PetPath, petUpdate.Pet.ID.Hex(), File.Filename)
		}
		// Create directories if they don't exist
		if err = os.MkdirAll(filepath.Dir(FilePath), 0755); err != nil {
			ctx.JSON(http.StatusInternalServerError, errors.ErrorWrapper{Message: "Failed To Create Directories ! ", Errors: err})
			return
		}
		// Save the uploaded file with the temporary path
		if err = ctx.SaveUploadedFile(File, FilePath); err != nil {
			ctx.JSON(http.StatusBadRequest, errors.ErrorWrapper{Message: "Failed to Upload Image !", Error: err.Error()})
			return
		}
		petUpdate.Pet.Image = append(petUpdate.Pet.Image, File.Filename)
	}
	return petUpdate, nil
}
