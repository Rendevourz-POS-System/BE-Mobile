package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	appConfig "main.go/configs/app"
	_const "main.go/configs/const"
	"main.go/configs/database"
)

func init() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatal("Error loading local.env file")
	}
}

func main() {
	app := gin.Default()
	db := database.ConnectDatabase(_const.DB_SHELTER_APP)
	defer func() {
		err := database.CloseDatabase()
		if err != nil {
			log.Fatalf("Error closing database: %v", err)
		}
	}()
	// Migrate database
	Migrate(db, _const.DB_SHELTER_APP)
	// Create indexes
	SetupDatabaseIndexes(db, _const.DB_SHELTER_APP)
	// Register Trusted Proxy And Certificate
	RegisterTrustedProxies(app)
	// Register Middleware
	RegisterMiddlewares(app)
	// Register Routes
	RegisterRoutes(app)
	err := app.Run(fmt.Sprintf(":%d", appConfig.GetConfig().App.Port))
	if err != nil {
		log.Fatalf("Failed to run the app: %v", err)
	}
}
