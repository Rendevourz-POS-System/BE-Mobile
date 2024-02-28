package src

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"main.go/configs/app"
	"main.go/configs/database"
	"main.go/middlewares"
	"net/http"
)

func Migrate(db *gorm.DB) {
	err := database.Migrate(db)
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

func RegisterRoutes(router *gin.Engine) {
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Hello World!",
		})
	})
}
