package repository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	Request "main.go/domains/request/entities"
	"main.go/domains/request/presistence"
	Pet "main.go/domains/shelter/entities"
	"main.go/shared/collections"
	"main.go/shared/helpers"
	"strings"
)

type requestRepo struct {
	database   *mongo.Database
	collection *mongo.Collection
}

func NewRequestRepository(database *mongo.Database) *requestRepo {
	return &requestRepo{database: database, collection: database.Collection(collections.RequestName)}
}

func (r *requestRepo) StoreOneRequest(ctx context.Context, req *Request.Request) (res *Request.Request, err error) {
	if strings.ToLower(req.Type) == "adoption" {
		var findPet *Pet.Pet
		errs := r.database.Collection(collections.PetCollectionName).FindOne(ctx, bson.M{"_id": req.PetId}).Decode(&findPet)
		if errs != nil {
			if err == mongo.ErrNoDocuments {
				return nil, errors.New("Pet not found !")
			}
			return nil, err
		}
		if *findPet.IsAdopted == true {
			return nil, errors.New("Pet Already Adopted !")
		}
		if *findPet.ReadyToAdopt == false {
			return nil, errors.New("Pet Not Ready For Adopt !")
		}
	}
	data, errInsert := r.collection.InsertOne(ctx, req)
	if errInsert != nil {
		return nil, errInsert
	}
	if err = r.collection.FindOne(ctx, bson.M{"_id": data.InsertedID}).Decode(&res); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *requestRepo) filterRequest(search *Request.SearchRequestPayload) bson.D {
	var filter bson.D
	if search.RequestId != nil {
		filter = append(filter, bson.E{
			Key:   "status",
			Value: search.RequestId,
		})
	}
	if search.Search != nil {
		regexFilter := bson.M{"$regex": primitive.Regex{
			Pattern: *search.Search,
			Options: "i", // Case-insensitive search
		}}
		filter = append(filter, bson.E{
			Key: "$or",
			Value: bson.A{
				bson.M{"reason": regexFilter},
				bson.M{"type": regexFilter},
			},
		})
	}
	if search.Status != nil {
		filter = append(filter, bson.E{
			Key:   "status",
			Value: helpers.RegexCaseInsensitivePattern(*search.Status),
		})
	}
	if search.ShelterId != nil {
		filter = append(filter, bson.E{
			Key:   "shelter_id",
			Value: *search.ShelterId,
		})
	}
	if search.UserId != nil {
		filter = append(filter, bson.E{
			Key:   "user_id",
			Value: *search.UserId,
		})
	}
	if search.Type != nil && len(*search.Type) > 0 {
		filter = append(filter, bson.E{
			Key:   "type",
			Value: bson.M{"$in": search.Type},
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

func (r *requestRepo) createPaginationOptions(search *Request.SearchRequestPayload) *options.FindOptions {
	findOptions := options.Find()
	page := 1
	pageSize := 100
	if search.PageSize > 0 {
		pageSize = search.PageSize
	}
	if search.Page > 0 {
		page = search.Page
	}
	findOptions.SetSkip(int64((page - 1) * pageSize))
	findOptions.SetLimit(int64(pageSize))
	return findOptions
}

func (r *requestRepo) FindAllRequest(ctx context.Context, req *Request.SearchRequestPayload) (res []Request.Request, err error) {
	filter := r.filterRequest(req)
	paginationOptions := r.createPaginationOptions(req)
	cur, err := r.collection.Find(ctx, filter, paginationOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	if errs := cur.All(ctx, &res); errs != nil {
		return nil, errs
	}
	return res, nil
}

func (r *requestRepo) FindOneRequestByData(ctx context.Context, data *bson.M) (res *Request.Request, err error) {
	if err = r.collection.FindOne(ctx, data).Decode(&res); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("Request Not Found ! ")
		}
		return nil, err
	}
	return res, nil
}

func (r *requestRepo) PutStatusRequest(ctx context.Context, req *Request.UpdateRescueAndSurrenderRequestStatus) (res *Request.UpdateRescueAndSurrenderRequestStatusResponse, err []string) {
	var (
		request *Request.Request
		pet     *Pet.Pet
	)
	filter := bson.M{"_id": req.RequestId} // Adjust the filter as per your requirements
	// Define the update operation to update only the `reason` field
	update := bson.M{
		"$set": bson.M{
			"status": presistence.Status(req.Status),
			"reason": req.Reason,
		},
	}
	// Perform the find and update operation
	errs := r.collection.FindOneAndUpdate(ctx, filter, update).Decode(&request)
	if errs != nil {
		if errs == mongo.ErrNoDocuments {
			return nil, []string{"Request Not Found!"}
		}
		return nil, []string{errs.Error()}
	}
	filterPet := bson.M{"_id": request.PetId}
	updatePet := bson.M{
		"$set": bson.M{
			"shelter_id": request.ShelterId,
		},
	}
	errUpdatePet := r.database.Collection(collections.PetCollectionName).FindOneAndUpdate(ctx, filterPet, updatePet).Decode(&pet)
	if errUpdatePet != nil {
		if errUpdatePet == mongo.ErrNoDocuments {
			return nil, []string{"Request Not Found!"}
		}
		return nil, []string{errUpdatePet.Error()}
	}
	res.Request = request
	res.Pet = pet
	return res, nil
}
