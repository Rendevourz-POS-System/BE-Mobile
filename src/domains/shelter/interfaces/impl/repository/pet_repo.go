package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	Pet "main.go/domains/shelter/entities"
	"main.go/shared/collections"
	"main.go/shared/helpers"
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
	if search.ShelterId != "" {
		filter = append(filter, bson.E{
			Key:   "shelter_id",
			Value: helpers.ParseStringToObjectId(search.ShelterId),
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
	if search.Gender != "" {
		filter = append(filter, bson.E{
			Key:   "pet_gender",
			Value: helpers.RegexCaseInsensitivePattern(search.Gender),
		})
	}
	if search.Type != "" {
		filter = append(filter, bson.E{
			Key:   "pet_type",
			Value: helpers.RegexCaseInsensitivePattern(search.Type),
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

func (r *petRepo) createPipeline(filter *bson.D, search *Pet.PetSearch) mongo.Pipeline {
	pipeline := mongo.Pipeline{
		{{"$match", filter}}, // Apply search filters
	}
	// Create Shelter pipeline (Aggregated)
	pipeline = r.createShelterPipeline(pipeline, search)
	// Create Location pipeline (Aggregated)
	pipeline = r.createLocationPipeline(pipeline, search)

	// Adjusting the $project stage
	// Here we add extra fields and assume all other pet fields should be included
	pipeline = append(pipeline, bson.D{{"$addFields", bson.M{
		"shelter_name":          "$shelter.shelter_name",   // Adds shelter_name field
		"shelter_location":      "$location.location_name", // Adds location_name field
		"shelter_location_name": "$location.location_name", // Adds shelter_location field
	}}})
	return pipeline
}

func (r *petRepo) createShelterPipeline(pipeline mongo.Pipeline, search *Pet.PetSearch) mongo.Pipeline {
	// Lookup to fetch the corresponding shelter
	pipeline = append(pipeline, bson.D{{
		"$lookup", bson.M{
			"from":         "shelters",
			"localField":   "shelter_id",
			"foreignField": "_id",
			"as":           "shelter",
		}},
	})
	// Unwind the result to simplify processing (consider handling missing shelters)
	pipeline = append(pipeline, bson.D{{"$unwind", bson.M{
		"path":                       "$shelter",
		"preserveNullAndEmptyArrays": true, // Keeps pets even if the shelter is missing
	}}})

	// Add location filter if specified and make it case insensitive
	if search.ShelterName != "" {
		regexPattern := helpers.RegexCaseInsensitivePattern(search.ShelterName)
		pipeline = append(pipeline, bson.D{{"$match", bson.M{"shelter.shelter_name": regexPattern}}})
	}
	return pipeline
}

func (r *petRepo) createLocationPipeline(pipeline mongo.Pipeline, search *Pet.PetSearch) mongo.Pipeline {
	// Additional lookup to fetch the location from the shelter
	pipeline = append(pipeline, bson.D{{
		"$lookup", bson.M{
			"from":         "shelter_locations",
			"localField":   "shelter.shelter_location",
			"foreignField": "_id",
			"as":           "location",
		}},
	})

	// Unwind the location (consider handling missing locations)
	pipeline = append(pipeline, bson.D{{"$unwind", bson.M{
		"path":                       "$location",
		"preserveNullAndEmptyArrays": true,
	}}})

	// Add location filter if specified and make it case insensitive
	if search.Location != "" {
		regexPattern := helpers.RegexCaseInsensitivePattern(search.Location)
		pipeline = append(pipeline, bson.D{{"$match", bson.M{"location.location_name": regexPattern}}})
	}
	return pipeline
}

func (r *petRepo) createPaginationPipeline(pipeline mongo.Pipeline, search *Pet.PetSearch) mongo.Pipeline {
	// Pagination can be added here if required
	if search.Page > 0 && search.PageSize > 0 {
		skip := (search.Page - 1) * search.PageSize
		pipeline = append(pipeline, bson.D{{"$skip", skip}}, bson.D{{"$limit", search.PageSize}})
	}
	return pipeline
}

func (r *petRepo) FindAllPets(ctx context.Context, search *Pet.PetSearch) (res []Pet.PetResponsePayload, err error) {
	filter := r.filterPets(search)                          // Filter
	pipeline := r.createPipeline(filter, search)            // Create pipeline
	pipeline = r.createPaginationPipeline(pipeline, search) // Create pagination pipeline
	data, errs := r.collection.Aggregate(ctx, pipeline)
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

func (r *petRepo) UpdatePet(ctx context.Context, pet *Pet.Pet) (res *Pet.Pet, errs error) {
	filter := bson.D{{Key: "_id", Value: pet.ID}}
	update := bson.D{{Key: "$set", Value: pet}}
	// Perform the update operation
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, mongo.ErrNoDocuments
	}
	// Optionally, you can retrieve the updated document
	err = r.collection.FindOne(ctx, filter).Decode(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
