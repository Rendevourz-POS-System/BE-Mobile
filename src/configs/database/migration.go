package database

import (
	"gorm.io/gorm"
	User "main.go/domains/user/entities"
)

func Migrate(db *gorm.DB) error {
	return db.Debug().AutoMigrate(
		&User.User{})
}
