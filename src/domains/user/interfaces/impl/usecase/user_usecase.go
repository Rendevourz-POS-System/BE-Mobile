package usecase

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"main.go/configs/app"
	User "main.go/domains/user/entities"
	"main.go/domains/user/interfaces"
	"main.go/domains/user/mail/controller"
	"main.go/shared/helpers"
	"main.go/shared/helpers/image_helpers"
	"os"
)

type userUsecase struct {
	userRepo interfaces.UserRepository
}

func NewUserUsecase(userRepo interfaces.UserRepository) *userUsecase {
	return &userUsecase{userRepo}
}

func (u *userUsecase) GetAllData(ctx context.Context) ([]User.User, error) {
	data, err := u.userRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (u *userUsecase) RegisterUser(ctx context.Context, user *User.User) (res *User.User, errs []string) {
	var err error
	validate := helpers.NewValidator()
	if err = validate.Struct(user); err != nil {
		errs = helpers.CustomError(err)
		return nil, errs
	}
	data := u.setDefaultUserData(user)
	resData, checkUserData, err2 := u.userRepo.StoreOne(ctx, data)
	if err2 != nil {
		errs = append(errs, err2.Error())
		return nil, errs
	}
	if resData != nil && !checkUserData {
		secretCode, Otp, err := u.userRepo.GenerateAndStoreToken(ctx, resData.ID, resData.Email)
		if err != nil {
			errs = append(errs, err.Error())
			return nil, errs
		}
		_, SendEmailVerification := u.SendEmailVerification(ctx, user, secretCode, Otp)
		if SendEmailVerification != nil {
			errs = append(errs, SendEmailVerification.Error())
			return nil, errs
		}
		return resData, nil
	}
	return resData, nil
}

func (u *userUsecase) setDefaultUserData(user *User.User) *User.User {
	StaffSatus := helpers.CheckStaffStatus(user.Role)
	return &User.User{
		Nik:                user.Nik,
		Email:              user.Email,
		Username:           user.Username,
		PostalCode:         user.PostalCode,
		Province:           user.Province,
		PhoneNumber:        user.PhoneNumber,
		Password:           helpers.HashPassword(user.Password),
		Address:            user.Address,
		City:               user.City,
		District:           user.District,
		Verified:           false,
		ShelterIsActivated: false,
		State:              user.State,
		Image:              "",
		ImageBase64:        "",
		StaffStatus:        StaffSatus,
		Role:               helpers.GetRole(user.Role),
		CreatedAt:          helpers.GetCurrentTime(nil),
	}
}

func (u *userUsecase) SendEmailVerification(ctx context.Context, data *User.User, secretCode string, Otp *int) (res *User.User, err error) {
	// send email verification
	ok := controller.SendEmail(&User.MailSend{
		To:      data.Email,
		Subject: "Email Verification",
		Content: helpers.ParsePointerIntToString(Otp),
		Cc:      "",
		Bcc:     "",
		Attach:  app.GetConfig().Email.Attachments,
	})
	if ok != nil {
		return nil, fmt.Errorf("failed to send email verification ! ")
	}
	return res, nil
}

func (u *userUsecase) LoginUser(ctx context.Context, userReq *User.LoginPayload) (*User.LoginResponse, error) {
	var res = &User.LoginResponse{}
	user, err := u.userRepo.FindByEmail(ctx, userReq.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user has not register yet ! ")
	}
	ok := helpers.ComparePassword(user.Password, userReq.Password)
	if !ok {
		return nil, errors.New("password or email doesn't match ! ")
	}
	if !user.Verified {
		res.User = *user
		return res, errors.New("user is not active ! ")
	}
	token, err := helpers.GenerateToken(user)
	if err != nil {
		return nil, err
	}
	res.Token = token
	res.Username = user.Username
	return res, nil
}

func (u *userUsecase) GetUserByUserId(ctx context.Context, id string) (*User.User, error) {
	user, err := u.userRepo.FindUserById(ctx, id)
	if user.Image != "" && len(user.Image) > 0 {
		file, err2 := os.ReadFile(image_helpers.GenerateImagePath(
			app.GetConfig().Image.UserPath, app.GetConfig().Image.ProfilePath, user.ID.Hex(), user.Image))
		if err2 != nil {
			return nil, err2
		}
		user.ImageBase64 = base64.StdEncoding.EncodeToString(file) // Convert to Base64
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) GetUserByUserIdForUpdate(ctx context.Context, id string, newImage *string) (*User.User, error) {
	user, err := u.userRepo.FindUserById(ctx, id)
	if user.Image != "" && len(user.Image) > 0 {
		if newImage != nil {
			user.Image = *newImage
		}
		file, err2 := os.ReadFile(image_helpers.GenerateImagePath(
			app.GetConfig().Image.UserPath, app.GetConfig().Image.ProfilePath, id, user.Image))
		if err2 != nil {
			return nil, err2
		}
		user.ImageBase64 = base64.StdEncoding.EncodeToString(file) // Convert to Base64
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) UpdateUserData(ctx context.Context, req *User.UpdateProfilePayload) (res *User.User, errs []string) {
	var err error
	var newImage *string = nil
	validate := helpers.NewValidator()
	if err = validate.Struct(req); err != nil {
		errs = helpers.CustomError(err)
		return nil, errs
	}
	if req.Image != "" && len(req.Image) > 0 {
		newImage = &req.Image
	}
	userDB, err := u.GetUserByUserIdForUpdate(ctx, req.ID.Hex(), newImage)
	if err != nil {
		errs = append(errs, err.Error())
		return nil, errs
	}
	data, err := u.userRepo.PutUser(ctx, u.updateFindUser(req, userDB))
	if err != nil {
		errs = append(errs, err.Error())
		return nil, errs
	}
	return data, nil
}

func (u *userUsecase) updateFindUser(user *User.UpdateProfilePayload, userDB *User.User) *User.User {
	var userImage = userDB.Image
	if user.Image != "" {
		userImage = user.Image
	}
	return &User.User{
		ID:                 user.ID,
		Nik:                user.Nik,
		PhoneNumber:        user.PhoneNumber,
		Address:            user.Address,
		State:              user.State,
		City:               user.City,
		Province:           user.Province,
		District:           user.District,
		PostalCode:         user.PostalCode,
		Email:              user.Email,
		Username:           user.Username,
		Password:           userDB.Password,
		StaffStatus:        userDB.StaffStatus,
		ShelterIsActivated: userDB.ShelterIsActivated,
		Role:               userDB.Role,
		Image:              userImage,
		Verified:           userDB.Verified,
		CreatedAt:          userDB.CreatedAt,
		UpdatedAt:          helpers.GetCurrentTime(nil),
		DeletedAt:          nil,
	}
}

func (u *userUsecase) UpdatePassword(ctx context.Context, req *User.UpdatePasswordPayload) error {
	validate := helpers.NewValidator()
	if err := validate.Struct(req); err != nil {
		errs := helpers.CustomError(err)
		return errors.New(errs[0])
	}
	findUser, err2 := u.userRepo.FindUserById(ctx, req.Id.Hex())
	if err2 != nil {
		return err2
	}
	if !helpers.ComparePassword(findUser.Password, req.Password) {
		return errors.New(fmt.Sprintf("Password Or Email Doest Not Match ! "))
	}
	req.NewPassword = helpers.HashPassword(req.NewPassword)
	if err := u.userRepo.PutUserPassword(ctx, req); err != nil {
		return errors.New(fmt.Sprintf("Failed To Update User Password ! "))
	}
	return nil
}

func (u *userUsecase) VerifyEmailVerification(ctx context.Context, req *User.EmailVerifiedPayload, userTokenUsecase interfaces.UserTokenUsecase) (res *User.User, err []string) {
	fmt.Println("Id --> ", req.UserId.Hex())
	validate := helpers.NewValidator()
	if errs := validate.Struct(req); errs != nil {
		err = helpers.CustomError(errs)
		return nil, err
	}
	//claims, errs := helpers.ClaimsJwtTokenForVerificationEmail(req.Token)
	//if errs != nil {
	//	err = append(err, errs.Error())
	//	return nil, err
	//}
	userId, errFindToken := userTokenUsecase.FindValidTokenByUserId(ctx, req.UserId, req.Otp)
	if errFindToken != nil {
		err = append(err, errFindToken.Error())
		return nil, err
	}
	if *userId != *req.UserId {
		err = append(err, errors.New("Invalid User Id ! ").Error())
		return nil, err
	}
	findUser, errFindUser := u.userRepo.FindUserById(ctx, req.UserId.Hex())
	if errFindUser != nil {
		err = append(err, errFindUser.Error())
		return nil, err
	}
	if findUser.Verified {
		err = append(err, "Email Verified Already ! ")
		return findUser, err
	}
	res, errVerified := u.userRepo.VerifiedUserEmail(ctx, req.UserId)
	if errVerified != nil {
		err = append(err, errVerified.Error())
		return nil, err
	}
	return res, nil
}
