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
	dsn := os.Getenv("DB_SOURCE")
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
		&models.CartInfo{},
		&models.Cart{},
		&models.Wallet{},
		&models.Checkoutinfo{},
		&models.RazorPay{},
		&models.Addimage{},
		&models.Wallethistory{},
	)
}
