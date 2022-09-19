package models

import "gorm.io/gorm"

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
	Cart          Cart
	CartInfo      CartInfo
	Wishlist      Wishlist
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
type Discount struct {
	gorm.Model
	DiscountName       string
	DiscountPercentage int
	ProductId          int
}

type Wallet struct {
	UsersID uint
	Balance int
}

type Checkoutinfo struct {
	gorm.Model
	UsersID        int
	Discount       int
	CouponDiscount int
	CouponCode     string
	TotalMrp       int
	Total          int
}
