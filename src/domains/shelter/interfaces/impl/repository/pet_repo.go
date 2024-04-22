package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	Pet "main.go/domains/shelter/entities"
	"main.go/shared/collections"
)

type petRepo struct {
	database   *mongo.Database
	collection *mongo.Collection
}

func NewPetRepository(database *mongo.Database) *petRepo {
	return &petRepo{database, database.Collection(collections.PetCollectionName)}
}

func (r *petRepo) filterPets(search *Pet.PetSearch) *bson.D {
	filter := bson.D{}
	if search.Search != "" {
		regexFilter := bson.M{"$regex": primitive.Regex{
			Pattern: search.Search,
			Options: "i", // Case-insensitive search
		}}
		filter = append(filter, bson.E{
			Key: "$or",
			Value: bson.A{
				bson.M{"pet_name": regexFilter},
				bson.M{"pet_age": regexFilter},
				bson.M{"pet_description": regexFilter},
			},
		})
	}
	if search.ShelterId != primitive.NilObjectID {
		filter = append(filter, bson.E{
			Key:   "shelter_id",
			Value: search.ShelterId,
		})
	}
	if search.Location != "" {
		filter = append(filter, bson.E{
			Key:   "location",
			Value: search.Location,
		})
	}
	if search.AgeEnd <= 0 {
		search.AgeEnd = 100
	}
	if search.AgeStart > 0 {
		filter = append(filter, bson.E{
			Key: "pet_age",
			Value: bson.M{
				"$gte": search.AgeStart,
				"$lte": search.AgeEnd,
			},
		})
	}
	if search.Type != "" {
		filter = append(filter, bson.E{
			Key:   "pet_type",
			Value: search.Type,
		})
	}
	return &filter
}

func (r *petRepo) paginationPets(search *Pet.PetSearch) *options.FindOptions {
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

func (r *petRepo) FindAllPets(ctx context.Context, search *Pet.PetSearch) (res []Pet.Pet, err error) {
	filter := r.filterPets(search)          // Filter
	findOptions := r.paginationPets(search) // Pagination
	data, errs := r.collection.Find(ctx, *filter, findOptions)
	if errs != nil {
		return nil, errs
	}
	if err = data.All(nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *petRepo) StorePets(ctx context.Context, data *Pet.Pet) (res *Pet.Pet, err []string) {
	var errs error
	var insertedResult *mongo.InsertOneResult
	if insertedResult, errs = r.collection.InsertOne(ctx, data); errs != nil {
		err = append(err, errs.Error())
		return nil, err
	}
	if errs = r.collection.FindOne(ctx, bson.M{"_id": insertedResult.InsertedID}).Decode(&res); errs != nil {
		err = append(err, errs.Error())
		return nil, err
	}
	return res, nil
}
