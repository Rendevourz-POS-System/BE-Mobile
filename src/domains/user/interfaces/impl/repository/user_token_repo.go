package repository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	JwtEmailClaims "main.go/src/domains/user/entities"
	"main.go/src/shared/collections"
	"main.go/src/shared/helpers"
	"time"
)

type userTokenRepo struct {
	database   *mongo.Database
	collection *mongo.Collection
}

func NewUserTokenRepository(database *mongo.Database) *userTokenRepo {
	return &userTokenRepo{database, database.Collection(collections.UserTokenCollectionName)}
}

func (repo *userTokenRepo) StoreToken(token string) error {
	return nil
}

func (repo *userTokenRepo) FindOneUserTokenByNonce(ctx context.Context, claims *JwtEmailClaims.JwtEmailClaims) (*primitive.ObjectID, error) {
	userToken := &JwtEmailClaims.UserToken{}
	if err := repo.collection.FindOne(ctx, bson.M{"_id": helpers.ParseStringToObjectId(claims.ID), "Otp": claims.Otp, "IsUsed": false}).Decode(&userToken); err != nil {
		return nil, err
	}
	// Check if the token is expired
	if userToken.ExpiredAt.Before(time.Now()) {
		return nil, errors.New("nonce is expired")
	}
	// Update the IsUsed field to true if the token is not expired
	update := bson.M{
		"$set": bson.M{
			"IsUsed": true,
		},
	}
	filter := bson.M{
		"_id":   claims.ID,
		"Token": claims.Nonce,
	}
	// Perform the update
	_, err := repo.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return &userToken.UserId, nil
}
func (repo *userTokenRepo) FindValidTokenByUserId(ctx context.Context, userId *primitive.ObjectID, Otp *int) (*primitive.ObjectID, error) {
	userToken := &JwtEmailClaims.UserToken{}
	if err := repo.collection.FindOne(ctx, bson.M{"UserId": userId, "Otp": *Otp, "IsUsed": false}).Decode(&userToken); err != nil {
		return nil, err
	}
	// Check if the token is expired
	if userToken.ExpiredAt.Before(time.Now()) {
		return nil, errors.New("nonce is expired")
	}
	// Update the IsUsed field to true if the token is not expired
	update := bson.M{
		"$set": bson.M{
			"IsUsed": true,
		},
	}
	filter := bson.M{
		"UserId": userId,
		"Otp":    Otp,
	}
	// Perform the update
	_, err := repo.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	// Perform the Delete operation
	deleteFilter := bson.M{
		"UserId": userId,
		"IsUsed": false,
	}
	_, err = repo.collection.DeleteMany(ctx, deleteFilter)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("Failed To Delete Res Data ! ")
		}
	}
	return &userToken.UserId, nil
}
