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

type Coupon struct {
	gorm.Model
	CouponName       string
	CouponCode       string
	CouponPercentage int
}

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
	Wallet       Wallet
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

type CartInfo struct {
	gorm.Model
	UsersId      uint
	ProductsID   uint
	Price        int
	Discount     int
	SellingPrice int
}

type Wishlist struct {
	gorm.Model
	UsersID    uint
	ProductsID uint
}

type Orders struct {
	gorm.Model
	UsersID        uint
	AddressID      uint
	OrderID        string
	Discount       int
	CouponDiscount int
	CouponCode     string
	Payment_Method string
	Total_Amount   int
	PaymentStatus  string
	OrderStatus    string
}

type Ordereditems struct {
	gorm.Model
	UsersID        uint
	ProductsID     uint
	Order_ID       string
	Product_Name   string
	Price          int
	CouponDiscount int
	Discount       int
	PaymentStatus  string
	OrderStatus    string
	Payment_Method string
	//AmountPaid     int
}

type Otp struct {
	gorm.Model
	Mobile string
	Otp    string
}

type RazorPay struct {
	gorm.Model
	UserId          string
	RazorPaymentId  string
	RazorPayOrderID string
	Signature       string
	OrderId         string
	AmountPaid      string
}

type Addimage struct {
	gorm.Model
	Cover string
}

type Wallethistory struct {
	gorm.Model
	UsersID uint
	Credit  int
	Debit   int
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
	OrderID        string
	Discount       int
	CouponDiscount int
	CouponCode     string
	TotalMrp       int
	Total          int
}
