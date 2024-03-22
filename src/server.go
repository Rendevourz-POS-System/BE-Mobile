package main

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"main.go/configs/app"
	"main.go/configs/database"
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
	RegisterMiddlewares(router)
}
