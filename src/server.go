package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go/example"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"main.go/configs/app"
	"main.go/configs/database"
	Master "main.go/domains/master/handlers/http"
	"main.go/domains/payment/interfaces/impl/repository"
	"main.go/domains/payment/interfaces/impl/usecase"
	Request "main.go/domains/request/handlers/http"
	Shelter "main.go/domains/shelter/handlers/http"
	UserHttp "main.go/domains/user/handlers/http"
	"main.go/middlewares"
	"net/http"
)

func Migrate(db *mongo.Client, dbName string) {
	err := database.Migrate(db, dbName)
	if err != nil {
		panic("Error migrating database : " + err.Error())
	}
}

func SetupDatabaseIndexes(db *mongo.Client, dbName string) {
	// Setting up indexes for the PetType collection
	petTypeIndexModel := mongo.IndexModel{
		Keys:    bson.M{"type": 1}, // Ensure unique index on "type"
		Options: options.Index().SetUnique(true),
	}
	shelterLocationIndexModel := mongo.IndexModel{
		Keys:    bson.M{"location_name": 1}, // Ensure unique index on "type"
		Options: options.Index().SetUnique(true),
	}
	userNikIndexModel := mongo.IndexModel{
		Keys:    bson.M{"nik": 1}, // Ensure unique index on "type"{
		Options: options.Index().SetUnique(true),
	}
	userEmailIndexModel := mongo.IndexModel{
		Keys:    bson.M{"email": 1}, // Ensure unique index on "type"{
		Options: options.Index().SetUnique(true),
	}
	// Create an index using the CreateOne() method
	_, err := database.User.Indexes().CreateOne(context.TODO(), userNikIndexModel)
	if err != nil {
		log.Fatal(err)
	}
	_, err = database.User.Indexes().CreateOne(context.TODO(), userEmailIndexModel)
	if err != nil {
		log.Fatal(err)
	}
	_, err = database.PetType.Indexes().CreateOne(context.TODO(), petTypeIndexModel)
	if err != nil {
		log.Fatalf("Failed to create unique index on pet types: %v", err)
	}
	_, err = database.ShelterLocation.Indexes().CreateOne(context.TODO(), shelterLocationIndexModel)
	if err != nil {
		log.Fatalf("Failed to create unique index on shelter location : %v", err)
	}
}

//func WatchPetTypeChanges(db *mongo.Database, cache *SomeCacheType) {
//	pipeline := mongo.Pipeline{
//		{{"$match", bson.D{{"operationType", "insert"}}}},
//	}
//	options := options.ChangeStream().SetFullDocument(options.UpdateLookup)
//	changeStream, err := db.Collection("pet_types").Watch(context.Background(), pipeline, options)
//	if err != nil {
//		log.Fatalf("Failed to watch changes: %v", err)
//	}
//	defer changeStream.Close(context.Background())
//
//	for changeStream.Next(context.Background()) {
//		var change bson.M
//		if err := changeStream.Decode(&change); err != nil {
//			log.Printf("Could not decode change: %v", err)
//			continue
//		}
//		updatedDoc := change["fullDocument"].(bson.M)
//		cache.Update(updatedDoc["_id"].(primitive.ObjectID), updatedDoc)
//	}
//}

//func EnsureValidPetTypes(db *mongo.Client, dbName string) error {
//	pipeline := mongo.Pipeline{
//		{{"$lookup", bson.D{
//			{"from", "pet_types"},
//			{"localField", "pet_type_accepted"},
//			{"foreignField", "_id"},
//			{"as", "matched_pet_types"},
//		}}},
//		{{"$match", bson.D{
//			{"matched_pet_types", bson.D{{"$not", bson.D{{"$size", 0}}}}},
//		}}},
//		{{"$project", bson.D{
//			{"_id", 1},
//			{"ShelterName", 1},
//			{"matched_pet_types.Type", 1}, // Optional: project types to view which types are matched
//		}}},
//	}
//
//	cursor, err := db.Database(dbName).Collection(collections.ShelterCollectionName).Aggregate(context.TODO(), pipeline)
//	if err != nil {
//		return fmt.Errorf("failed to execute aggregation: %v", err)
//	}
//	var invalidShelters []bson.M
//	if err = cursor.All(context.Background(), &invalidShelters); err != nil {
//		return fmt.Errorf("failed to decode results: %v", err)
//	}
//
//	// Optionally handle invalid shelters, such as logging them or taking corrective action
//	for _, shelter := range invalidShelters {
//		fmt.Printf("Invalid shelter ID: %v\n", shelter["_id"])
//	}
//
//	return nil
//}

func RegisterTrustedProxies(router *gin.Engine) {
	app.TrustedProxies(router)
	return
}

func RegisterMiddlewares(router *gin.Engine) {
	router.Use(middlewares.NewCors(router))
	// Another Middlewares Here ...
}

func RegisterMiddlewares2(router *gin.Engine) {
	router.Use(middlewares.JwtAuthMiddleware(app.GetConfig().AccessToken.AccessTokenSecret, ""))
	// Another Middlewares Here ...
}
func RegisterRoutes(router *gin.Engine) {
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Hello World!",
		})
	})
	fmt.Println("Example Core Param --> ", example.CoreParam())
	fmt.Println("Example Client Radndom --> ", example.Random())
	midtransUsecase := usecase.NewMidtransUsecase(repository.NewMidtrans())
	userTokenHttp := UserHttp.NewUserTokenHttp(router)
	userHandler := UserHttp.NewUserHttp(router, userTokenHttp)
	shelterHttp := Shelter.NewShelterHttp(router)
	Shelter.NewPetHttp(router, shelterHttp)
	Shelter.NewShelterFavoriteHttp(router)
	Shelter.NewPetFavoriteHttp(router)
	Master.NewPetTypeHttp(router)
	Master.NewShelterLocationHttp(router)
	donationHandlers := Request.NewDonationShelterHttp(router)
	adoptionHandlers := Request.NewAdoptionShelterHttp(router)
	Request.NewRequestHttp(router, midtransUsecase, donationHandlers, adoptionHandlers, userHandler, shelterHttp)

}
