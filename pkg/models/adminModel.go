package models

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	Email    string
	Password string
}

type Tax struct {
	gorm.Model
	Category string
	Tax      int
}

type Discount struct {
	gorm.Model
	DiscountName       string
	DiscountPercentage int
	ProductId          int
}

type Coupon struct {
	gorm.Model
	CouponName       string
	CouponCode       string
	CouponPercentage int
}
