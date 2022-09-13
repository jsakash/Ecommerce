package models

import (
	"gorm.io/gorm"
)

type Products struct {
	gorm.Model
	Product_Name  string
	Brand_Name    string
	Product_Price int
	Description   string
	CategoryID    uint
	ColorsID      uint
	SizeID        uint
	Stock         int
	Cart Cart
	Wishlist Wishlist
	Orders Orders
	
}

type Category struct {
	gorm.Model
	Category string
	Products Products
}

type Colors struct {
	gorm.Model
	Color    string
	Products Products
}

type Size struct {
	gorm.Model
	Size     int
	Products Products
}
