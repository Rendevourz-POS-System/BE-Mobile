package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	PetType "main.go/domains/master/entities"
	"main.go/shared/collections"
)

type petTypeRepo struct {
	database   *mongo.Database
	collection *mongo.Collection
}

func NewPetTypeRepo(database *mongo.Database) *petTypeRepo {
	return &petTypeRepo{database, database.Collection(collections.PetTypeName)}
}

func (r *petTypeRepo) FindAllPets(ctx context.Context) ([]PetType.PetType, error) {
	res := []PetType.PetType{}
	data, err := r.collection.Find(ctx, bson.M{})
	if err = data.All(ctx, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *petTypeRepo) StorePetType(ctx context.Context, req *PetType.PetType) (res *PetType.PetType, err error) {
	// Attempt to insert the new pet type
	insertResult, err := r.collection.InsertOne(ctx, req)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			// This error means a document with the same "type" already exists
			return nil, fmt.Errorf("a pet type with the same type already exists")
		}
		return nil, err // Handle other potential errors during insertion
	}

	// Fetch and return the inserted document using the InsertedID
	if oid, ok := insertResult.InsertedID.(primitive.ObjectID); ok {
		result := &PetType.PetType{} // Assuming PetType is a struct that needs to be initialized
		err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(result)
		if err != nil {
			return nil, err // Handle error if the document could not be found
		}
		return result, nil
	} else {
		return nil, fmt.Errorf("failed to convert the inserted ID to ObjectID")
	}
}
