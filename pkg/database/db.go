package database

import (
	"os"

	"github.com/jsakash/ecommers/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDb() {
	var err error
	dsn := os.Getenv("DB")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed TO connect to DB")
	}

	DB.AutoMigrate(&models.Users{})
	DB.AutoMigrate(&models.Admin{})
	DB.AutoMigrate(&models.Products{})
}
