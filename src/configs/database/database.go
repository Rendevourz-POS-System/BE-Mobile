package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"main.go/configs/app"
	_const "main.go/configs/const"
	"main.go/shared/helpers"
	"net/url"
)

var (
	db *mongo.Client
)

func createDNS(conf app.Database, index int) string {
	safeDNS := url.UserPassword(conf.Username[index], conf.Password[index]).String()
	//log.Println("Connecting to database len: ", safeDNS, conf.Host[index], conf.Port[index])
	//return fmt.Sprintf("mongodb://%s@%s:%d/?tls=false&authMechanism=SCRAM-SHA-256",
	//	safeDNS, conf.Host[index], conf.Port[index],
	//)
	log.Println("Connecting to database len: ", safeDNS, conf.Host[index], conf.Port[index])
	return fmt.Sprintf("mongodb+srv://%s@%s/?retryWrites=true&w=majority&appName=shelter-apps-db",
		safeDNS, conf.Host[index],
	)
}

func ConnectDatabase(data string) *mongo.Client {
	conf := app.GetConfig().Database // Ensure this method correctly fetches your configuration
	if data == _const.DB_SHELTER_APP {
		if db == nil {
			DNS := createDNS(conf, 0)
			server := options.ServerAPI(options.ServerAPIVersion1)
			opts := options.Client().ApplyURI(DNS).SetServerAPIOptions(server)
			var err error
			db, err = mongo.Connect(context.Background(), opts)
			if err != nil {
				log.Fatalf("Error connecting to database: %v", err)
			}
			return db
		}
		return db
	}
	return db
}

func CloseDatabase() error {
	return db.Disconnect(context.Background())
}

func GetDatabase(data string) *mongo.Database {
	return db.Database(helpers.ParseDatabase(data))
}
