package models

import "gorm.io/gorm"

type Users struct {
	gorm.Model
	First_Name   string
	Last_Name    string
	Email        string `gorm:"unique"`
	Password     string
	Phone_Number string
	Status       bool
	Address      Address
	Cart         Cart
	Wishlist     Wishlist
	Orders       Orders
}

type Address struct {
	gorm.Model
	UsersID      uint
	Name         string
	Phone_number int
	Pincode      int
	House_Adress string
	Area         string
	Landmark     string
	City         string
	Orders       Orders
}

type Cart struct {
	gorm.Model
	UsersId     uint
	ProductsId  uint
	Quantity    int
	Total_Price int
}

type Wishlist struct {
	gorm.Model
	UsersID    uint
	ProductsID uint
}

type Orders struct {
	gorm.Model
	UsersID        uint
	ProductsId     uint
	AddressID      uint
	Order_ID       string
	Payment_Method string
	Order_Status   string
	Total_Amount   int
}

type Ordereditems struct {
	gorm.Model
	UsersID      uint
	ProductsID   uint
	Order_ID     int
	Order_Status string
}

type Otp struct {
	gorm.Model
	Mobile string
	Otp    string
}
