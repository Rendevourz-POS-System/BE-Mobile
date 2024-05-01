package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"main.go/configs/app"
	"main.go/configs/database"
	Master "main.go/domains/master/handlers/http"
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
	petTypeCollection := db.Database(dbName).Collection("pet_types")
	petTypeIndexModel := mongo.IndexModel{
		Keys:    bson.M{"type": 1}, // Ensure unique index on "type"
		Options: options.Index().SetUnique(true),
	}
	_, err := petTypeCollection.Indexes().CreateOne(context.TODO(), petTypeIndexModel)
	if err != nil {
		log.Fatalf("Failed to create unique index on pet types: %v", err)
	}
}

func RegisterTrustedProxies(router *gin.Engine) {
	app.TrustedProxies(router)
	return
}

func RegisterMiddlewares(router *gin.Engine) {
	router.Use(middlewares.NewCors(router))
	// Another Middlewares Here ...
}

func RegisterMiddlewares2(router *gin.Engine) {
	router.Use(middlewares.JwtAuthMiddleware(app.GetConfig().AccessToken.AccessTokenSecret))
	// Another Middlewares Here ...
}
func RegisterRoutes(router *gin.Engine) {
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Hello World!",
		})
	})
	UserHttp.NewUserHttp(router)
	UserHttp.NewUserTokenHttp(router)
	Shelter.NewShelterHttp(router)
	Shelter.NewPetHttp(router)
	Shelter.NewShelterFavoriteHttp(router)
	Master.NewPetTypeHttp(router)
}
