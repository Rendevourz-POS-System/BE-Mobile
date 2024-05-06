package repository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	User "main.go/domains/user/entities"
	"main.go/shared/collections"
	"main.go/shared/helpers"
)

type userRepository struct {
	database   *mongo.Database
	collection *mongo.Collection
}

func NewUserRepository(database *mongo.Database) *userRepository {
	return &userRepository{database, database.Collection(collections.UserCollectionName)}
}

func (userRepo *userRepository) FindAll(c context.Context) (res []User.User, err error) {
	data, err := userRepo.collection.Find(c, bson.D{})
	if err != nil {
		return nil, err
	}
	err = data.All(c, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (userRepo *userRepository) StoreOne(c context.Context, user *User.User) (*User.User, bool, error) {
	var existingUser User.User
	err := userRepo.collection.FindOne(c, bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		// A user with this email already exists, so return the existing user
		if !existingUser.Verified {
			return nil, false, errors.New("this account already exists and is not active yet ! ")
		}
		return &existingUser, true, nil
	} else if err != mongo.ErrNoDocuments {
		// An actual error occurred while trying to find the user, other than "no documents found"
		return nil, false, err
	}

	// Proceed with insertion if no existing user was found
	insertResult, err := userRepo.collection.InsertOne(c, user)
	if err != nil {
		return nil, false, err // Return the error encountered during insertion
	}

	// Ensure the InsertedID is an ObjectID to use it in the FindOne query
	objectID, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		// The type assertion failed, return a descriptive error
		return nil, false, errors.New("inserted ID is not of type ObjectID")
	}

	// Fetch the newly inserted document to return a complete user object, including its new _id
	var newUser User.User
	if err = userRepo.collection.FindOne(c, bson.M{"_id": objectID}).Decode(&newUser); err != nil {
		return nil, false, err // Return any error encountered during fetching
	}

	return &newUser, false, nil
}

func (userRepo *userRepository) FindByEmail(c context.Context, email string) (*User.User, error) {
	var user User.User
	err := userRepo.collection.FindOne(c, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found ! ")
		}
		return nil, err
	}
	return &user, nil
}

func (userRepo *userRepository) GenerateAndStoreToken(c context.Context, userId primitive.ObjectID, email string) (string, error) {
	minute := 30
	userToken := &User.UserToken{
		UserId:    userId,
		Token:     helpers.GenerateRandomString(32),
		IsUsed:    false,
		CreatedAt: helpers.GetCurrentTime(nil),
		ExpiredAt: helpers.GetCurrentTime(&minute),
		DeletedAt: nil,
	}
	data, err := userRepo.database.Collection(collections.UserTokenCollectionName).InsertOne(c, userToken)
	if err != nil {
		return "", err
	}
	// Fetch the newly inserted document to return a complete user object, including its new _id
	var newUserToken *User.UserToken
	if err = userRepo.database.Collection(collections.UserTokenCollectionName).FindOne(c, bson.M{"_id": data.InsertedID}).Decode(&newUserToken); err != nil {
		return "", err // Return any error encountered during fetching
	}
	//fmt.Printf("UserData : %v\n", newUserToken)
	secretCode, errs := helpers.GenerateJwtTokenForVerificationEmail(newUserToken.Id.Hex(), email, userToken.Token)
	if errs != nil {
		return "", errs
	}
	//fmt.Printf("SecretCode : %s", secretCode)
	//test, _ := helpers.ClaimsJwtTokenForVerificationEmail(secretCode)
	//fmt.Printf("CodeData : %v\n", test)
	return secretCode, nil
}

func (userRepo *userRepository) FindUserById(c context.Context, userId string) (res *User.User, errs error) {
	var user User.User
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}
	err = userRepo.collection.FindOne(c, bson.M{"_id": objectId}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (userRepo *userRepository) PutUser(ctx context.Context, user *User.User) (res *User.User, err error) {
	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}
	// Set the options to return the updated document
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	if err = userRepo.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&res); err != nil {
		return nil, err
	}
	return res, nil
}
