package repository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	Shelter "main.go/src/domains/shelter/entities"
	"main.go/src/shared/collections"
	"main.go/src/shared/helpers"
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

func (shelterRepo *shelterRepository) createLocationPipeline(pipeline mongo.Pipeline, search *Shelter.ShelterSearch) mongo.Pipeline {
	// Lookup to fetch the corresponding shelter
	pipeline = append(pipeline, bson.D{{
		"$lookup", bson.M{
			"from":         "shelter_locations",
			"localField":   "shelter_location",
			"foreignField": "_id",
			"as":           "location_details",
		}},
	})
	// Unwind the result to simplify processing (consider handling missing shelters)
	pipeline = append(pipeline, bson.D{{"$unwind", bson.M{
		"path":                       "$location_details",
		"preserveNullAndEmptyArrays": true, // Keeps pets even if the shelter is missing
	}}})
	if search != nil {
		if search.ShelterLocationName != "" {
			regexPattern := helpers.RegexCaseInsensitivePattern(search.ShelterLocationName)
			pipeline = append(pipeline, bson.D{{"$match", bson.M{"location_details.location_name": regexPattern}}})
		}
	}
	return pipeline
}

func (shelterRepo *shelterRepository) createPetsPipeline(pipeline mongo.Pipeline, search *Shelter.ShelterSearch) mongo.Pipeline {
	pipeline = append(pipeline, bson.D{{
		"$lookup", bson.M{
			"from":         "pet_types",
			"localField":   "pet_type_accepted",
			"foreignField": "_id",
			"as":           "pet_type_details",
		}},
	})
	// Unwind the result to simplify processing (consider handling missing shelters)
	pipeline = append(pipeline, bson.D{{"$unwind", bson.M{
		"path":                       "$pet_type_details",
		"preserveNullAndEmptyArrays": true, // Keeps pets even if the shelter is missing
	}}})

	if search.PetType != "" {
		regexPattern := helpers.RegexCaseInsensitivePattern(search.PetType)
		pipeline = append(pipeline, bson.D{{"$match", bson.M{"pet_type_details.type": regexPattern}}})
	}
	return pipeline
}

func (shelterRepo *shelterRepository) createPipeline(filter bson.D, search *Shelter.ShelterSearch) mongo.Pipeline {
	pipeline := mongo.Pipeline{
		{{"$match", filter}}, // Apply search filters
	}
	if search != nil {
		pipeline = shelterRepo.createLocationPipeline(pipeline, search)
	} else {
		pipeline = shelterRepo.createLocationPipeline(pipeline, nil)
	}

	//pipeline = shelterRepo.createPetsPipeline(pipeline, search)

	// Here we adjust fields
	// Here we add extra fields and assume all other pet fields should be included
	pipeline = append(pipeline, bson.D{{"$addFields", bson.M{
		"shelter_location_name": "$location_details.location_name",
		//"pet_type_accepted_name": bson.M{"$match": "$pet_type_details.type"},
	}}})
	return pipeline
}

func (shelterRepo *shelterRepository) createPaginationPipeline(pipeline mongo.Pipeline, search *Shelter.ShelterSearch) mongo.Pipeline {
	// Pagination can be added here if required
	if search.Page > 0 && search.PageSize > 0 {
		skip := (search.Page - 1) * search.PageSize
		pipeline = append(pipeline, bson.D{{"$skip", skip}}, bson.D{{"$limit", search.PageSize}})
	}
	return pipeline
}

func (shelterRepo *shelterRepository) getAllFavoriteShelters(c context.Context, userID *primitive.ObjectID) (shelterIDs []primitive.ObjectID, err error) {
	cursor, errs := shelterRepo.database.Collection(collections.ShelterFavoriteName).Find(c, bson.M{"user_id": userID})
	if errs != nil {
		if errs == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errs
	}
	for cursor.Next(c) {
		var favorite Shelter.ShelterFavorite
		if err = cursor.Decode(&favorite); err != nil {
			return nil, err
		}
		shelterIDs = append(shelterIDs, favorite.ShelterId)
	}
	return shelterIDs, nil
}

func (shelterRepo *shelterRepository) FindAllData(c context.Context, search *Shelter.ShelterSearch) (res []Shelter.ShelterResponsePayload, err error) {
	filter := shelterRepo.filterShelter(search) // Filter
	if search.UserId != primitive.NilObjectID {
		favoriteShelterIDs, errs := shelterRepo.getAllFavoriteShelters(c, &search.UserId)
		if errs != nil {
			return nil, errs
		}
		if len(favoriteShelterIDs) > 0 {
			filter = append(filter, bson.E{
				Key:   "_id",
				Value: bson.M{"$in": favoriteShelterIDs},
			})
		} else {
			return []Shelter.ShelterResponsePayload{}, nil
		}
	}
	pipeline := shelterRepo.createPipeline(filter, search)
	//findOptions := shelterRepo.paginationShelter(search) // Pagination
	pipeline = shelterRepo.createPaginationPipeline(pipeline, search) // Create pagination pipeline
	data, err := shelterRepo.collection.Aggregate(c, pipeline)
	if err != nil {
		return nil, err
	}
	err = data.All(c, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (shelterRepo *shelterRepository) FindOneDataByUserId(c context.Context, Id *primitive.ObjectID) (res *Shelter.ShelterResponsePayload, err error) {
	pipeline := shelterRepo.createPipeline(bson.D{{"user_id", Id}}, nil)
	data, err := shelterRepo.collection.Aggregate(c, pipeline)
	if err != nil {
		return nil, err
	}
	defer data.Close(c)
	if data.Next(c) {
		err = data.Decode(&res)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("User Does Not Have Shelter!")
	}
	return res, nil
}

func (shelterRepo *shelterRepository) FindOneDataById(c context.Context, Id *primitive.ObjectID) (res *Shelter.ShelterResponsePayload, err error) {
	pipeline := shelterRepo.createPipeline(bson.D{{"_id", Id}}, nil)
	data, err := shelterRepo.collection.Aggregate(c, pipeline)
	if err != nil {
		return nil, err
	}
	defer data.Close(c)
	if data.Next(c) {
		err = data.Decode(&res)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("User Does Not Have Shelter!")
	}
	return res, nil
}

func (shelterRepo *shelterRepository) FindOneDataByIdForRequest(c context.Context, Id *primitive.ObjectID) (res *Shelter.Shelter, err error) {
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
	shelter.ID = insertResult.InsertedID.(primitive.ObjectID)
	if _, err = shelterRepo.database.Collection(collections.UserCollectionName).UpdateOne(c, bson.M{"_id": shelter.UserId}, bson.M{"$set": bson.M{"shelter_is_activated": true}}); err != nil {
		return nil, err
	}
	return shelter, nil
}

func (shelterRepo *shelterRepository) UpdateOneShelter(ctx context.Context, shelter *Shelter.Shelter) (res *Shelter.Shelter, err error) {
	filter := bson.D{{Key: "_id", Value: shelter.ID}}
	update := bson.D{{Key: "$set", Value: shelter}}
	// Perform the update operation
	result, err := shelterRepo.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	if result.MatchedCount == 0 {
		return nil, mongo.ErrNoDocuments
	}
	// Optionally, you can retrieve the updated document
	err = shelterRepo.collection.FindOne(ctx, filter).Decode(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
