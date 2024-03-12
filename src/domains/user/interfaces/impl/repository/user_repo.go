package repository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	User "main.go/domains/user/entities"
	"main.go/domains/user/presistence"
)

type userRepository struct {
	database   *mongo.Database
	collection *mongo.Collection
}

func NewUserRepository(database *mongo.Database) *userRepository {
	return &userRepository{database, database.Collection(presistence.CollectionName)}
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

func (userRepo *userRepository) StoreOne(c context.Context, user *User.User) (*User.User, error) {
	var existingUser User.User
	err := userRepo.collection.FindOne(c, bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		// A user with this email already exists, so return the existing user
		return &existingUser, nil
	} else if err != mongo.ErrNoDocuments {
		// An actual error occurred while trying to find the user, other than "no documents found"
		return nil, err
	}

	// Proceed with insertion if no existing user was found
	insertResult, err := userRepo.collection.InsertOne(c, user)
	if err != nil {
		return nil, err // Return the error encountered during insertion
	}

	// Ensure the InsertedID is an ObjectID to use it in the FindOne query
	objectID, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		// The type assertion failed, return a descriptive error
		return nil, errors.New("inserted ID is not of type ObjectID")
	}

	// Fetch the newly inserted document to return a complete user object, including its new _id
	var newUser User.User
	if err := userRepo.collection.FindOne(c, bson.M{"_id": objectID}).Decode(&newUser); err != nil {
		return nil, err // Return any error encountered during fetching
	}

	return &newUser, nil
}
