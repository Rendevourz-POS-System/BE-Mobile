package repository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	Shelter "main.go/domains/shelter/entities"
	"main.go/shared/collections"
)

type shelterRepository struct {
	database   *mongo.Database
	collection *mongo.Collection
}

func NewShelterRepository(database *mongo.Database) *shelterRepository {
	return &shelterRepository{database, database.Collection(collections.ShelterCollectionName)}
}

func (shelterRepo *shelterRepository) paginationShelter(search *Shelter.ShelterSearch) *options.FindOptions {
	findOptions := options.Find()
	orderBy := "CreatedAt" // Default sorting field "CreatedAt
	sortOrder := 1         // Ascending
	// Sorting
	if search.Sort == "Desc" {
		sortOrder = -1 // Descending
	}
	if search.OrderBy != "" {
		orderBy = search.OrderBy
	}
	findOptions.SetSort(bson.D{{Key: orderBy, Value: sortOrder}})
	// Pagination
	if search.Page > 0 && search.PageSize > 0 {
		skip := (search.Page - 1) * search.PageSize
		findOptions.SetSkip(int64(skip))
		findOptions.SetLimit(int64(search.PageSize))
	}
	return findOptions
}

func (shelterRepo *shelterRepository) filterShelter(search *Shelter.ShelterSearch) bson.D {
	var filter bson.D
	if search.Search != "" {
		regexFilter := bson.M{"$regex": primitive.Regex{
			Pattern: search.Search,
			Options: "i", // Case-insensitive search
		}}
		filter = append(filter, bson.E{
			Key: "$or",
			Value: bson.A{
				bson.M{"shelter_name": regexFilter},
				bson.M{"shelter_description": regexFilter},
			},
		})
	}
	// Filter for non-deleted (soft delete check) records.
	filter = append(filter, bson.E{
		Key: "$or",
		Value: bson.A{
			bson.M{"deleted_at": nil},                      // Matches if `deleted_at` is explicitly set to null
			bson.M{"deleted_at": bson.M{"$exists": false}}, // Matches if `deleted_at` field does not exist
		},
	})
	return filter
}

func (shelterRepo *shelterRepository) FindAllData(c context.Context, search *Shelter.ShelterSearch) (res []Shelter.Shelter, err error) {
	filter := shelterRepo.filterShelter(search)          // Filter
	findOptions := shelterRepo.paginationShelter(search) // Pagination
	data, err := shelterRepo.collection.Find(c, filter, findOptions)
	if err != nil {
		return nil, err
	}
	err = data.All(c, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (shelterRepo *shelterRepository) FindOneDataByUserId(c context.Context, Id *primitive.ObjectID) (res *Shelter.Shelter, err error) {
	if err = shelterRepo.collection.FindOne(c, bson.M{"user_id": Id}).Decode(&res); err != nil {
		return nil, errors.New("User Does Not Have Shelter ! ")
	}
	return res, nil
}

func (shelterRepo *shelterRepository) FindOneDataById(c context.Context, Id *primitive.ObjectID) (res *Shelter.Shelter, err error) {
	if err = shelterRepo.collection.FindOne(c, bson.M{"_id": Id}).Decode(&res); err != nil {
		return nil, errors.New("User Does Not Have Shelter ! ")
	}
	return res, nil
}

func (shelterRepo *shelterRepository) StoreData(c context.Context, shelter *Shelter.Shelter) (res *Shelter.Shelter, errs error) {
	var existingShelter Shelter.Shelter
	// Check if the shelter name already exists
	if err := shelterRepo.collection.FindOne(c, bson.M{"shelter_name": shelter.ShelterName}).Decode(&existingShelter); err == nil {
		if existingShelter.UserId != shelter.UserId {
			return nil, errors.New("Shelter Name Already Exist ! ")
		}
	}
	if errs = shelterRepo.collection.FindOneAndUpdate(c, bson.M{"user_id": shelter.UserId}, bson.M{"$set": shelter}).Err(); errs == nil {
		if errs = shelterRepo.collection.FindOne(c, bson.M{"user_id": shelter.UserId}).Decode(&res); errs == nil {
			return res, errors.New("Shelter already exist & Updated ! ")
		}
		return res, errors.New("Shelter already exist & Updated ! ")
	}
	insertResult, err := shelterRepo.collection.InsertOne(c, shelter)
	if err != nil {
		return nil, err
	}
	shelter.ID = insertResult.InsertedID.(primitive.ObjectID).Hex()
	if _, err = shelterRepo.database.Collection(collections.UserCollectionName).UpdateOne(c, bson.M{"_id": shelter.UserId}, bson.M{"$set": bson.M{"shelter_is_activated": true}}); err != nil {
		return nil, err
	}
	return shelter, nil
}
