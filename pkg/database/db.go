package database

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
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

	DB.AutoMigrate(
		&models.Admin{},
		&models.Users{},
		&models.Address{},
		&models.Category{},
		&models.Colors{},
		&models.Size{},
		&models.Products{},
		&models.Cart{},
		&models.Wishlist{},
		&models.Orders{},
		&models.Ordereditems{},
		&models.Otp{},
		&models.Tax{},
		&models.Discount{},
		&models.Coupon{},
	)
}

//connecting database and returning database engine
func GetDb() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Unable to get data from env!!")
	}
	// dsn := os.Getenv("database_address")
	dsn := "host=localhost user=akashjs password=312002 dbname=ecommers port=5432 sslmode=disable"

	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database!")
	}
	return DB
}
