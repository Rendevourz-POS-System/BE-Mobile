package database

import (
	"fmt"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"main.go/configs/app"
	_const "main.go/configs/const"
	"net/url"
	"time"
)

var (
	gormConfig = &gorm.Config{
		Logger:               logger.Default.LogMode(logger.Info),
		FullSaveAssociations: true,
		//NamingStrategy: schema.NamingStrategy{
		//	TablePrefix: conf.Schema, // Set the schema name here
		//},
	}
	db *gorm.DB
)

func createDNS(conf app.Database, index int) string {
	query := url.Values{}
	query.Add("database", conf.Database[index])
	safeDNS := url.UserPassword(conf.Username[index], conf.Password[index]).String()
	log.Println("Connecting to database len: ", safeDNS, conf.Host[index], conf.Port[index])
	return fmt.Sprintf("sqlserver://%s@%s:%d?%s",
		safeDNS, conf.Host[index], conf.Port[index], query.Encode(),
	)
}

func ConnectDatabase(data string) *gorm.DB {
	conf := app.GetConfig().Database // Ensure this method correctly fetches your configuration
	if data == _const.DB_SHELTER_APP {
		if db == nil {
			DNS := createDNS(conf, 0)
			var err error
			//db, err = gorm.Open(sqlserver.Open(DNS), gormConfig)
			db, err = gorm.Open(sqlserver.Open(DNS), gormConfig)
			if err != nil {

				log.Fatalf("Error connecting to database : %v", err)
			}
			sqlDB, _ := db.DB()
			// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
			sqlDB.SetMaxIdleConns(10)
			// SetMaxOpenConns sets the maximum number of open connections to the database.
			sqlDB.SetMaxOpenConns(100)
			// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
			sqlDB.SetConnMaxLifetime(time.Hour)
			return db
		}
		return db
	}
	return db
}
