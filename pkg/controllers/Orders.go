package controllers

import (
	"math/rand"

	"github.com/gin-gonic/gin"
	"github.com/jsakash/ecommers/pkg/database"
	"github.com/jsakash/ecommers/pkg/models"
)

var letters = []rune("123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

func OrderIdGeneration(value int) string {
	b := make([]rune, value)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func OrderInfo(c *gin.Context) {

	var Orders []models.Orders
	database.DB.Find(&Orders)
	for _, i := range Orders {
		c.JSON(200, gin.H{
			"status":         true,
			"UserId":         i.UsersID,
			"OrderID":        i.OrderID,
			"Discount":       i.Discount,
			"CouponDiscount": i.CouponDiscount,
			"CouponCode":     i.CouponCode,
			"PaymentMethod":  i.Payment_Method,
			"TotalAmount":    i.Total_Amount,
		})
	}

}

func OrderedItems(c *gin.Context) {
	UsersID := c.GetUint("id")
	var items []models.Ordereditems
	database.DB.Where("users_id = ?", UsersID).Find(&items)

	for _, i := range items {

		c.JSON(200, gin.H{
			"status":          true,
			"id":              i.ID,
			"Amount_Paid":     i.Price,
			"ProductsID":      i.ProductsID,
			"Order_ID":        i.Order_ID,
			"Product_Name":    i.Product_Name,
			"Discount":        i.Discount,
			"Coupon_Discount": i.CouponDiscount,
			"Order Status":    i.OrderStatus,
		})
	}
}

func CancelOrder(c *gin.Context) {
	userID := c.GetUint("id")
	var items models.Ordereditems
	var updateStatus string = "CANCELLED"
	id := c.Query("ProductID")

	database.DB.First(&items, id)
	if items.OrderStatus == updateStatus {
		c.JSON(400, gin.H{
			"status":  false,
			"message": "Order already Cancelled",
		})
		return
	}
	database.DB.Model(&items).Where("id=?", id).Update("order_status", updateStatus)

	var price int
	database.DB.Raw("SELECT price FROM ordereditems WHERE id = ?", id).Scan(&price)

	var balance int
	database.DB.Raw("SELECT balance FROM wallets WHERE users_id = ?", userID).Scan(&balance)
	newBalance := balance + price

	if items.Payment_Method == "COD" {
		c.JSON(200, gin.H{
			"status":  true,
			"message": "Order Cancelled",
		})
		return
	}

	WalletHistory := models.Wallethistory{UsersID: userID, Debit: 0, Credit: price}
	database.DB.Create(&WalletHistory)

	var totalAmount int
	database.DB.Raw("SELECT total_amaount FROM orders WHERE users_id = ?", userID).Scan(&totalAmount)
	Ntotal := totalAmount - balance
	//var wallet models.Wallet
	database.DB.Model(&models.Wallet{}).Where("users_id = ?", userID).Update("balance", newBalance)
	database.DB.Model(&models.Orders{}).Where("users_id = ?", userID).Update("total_amount", Ntotal)
	//database.DB.Raw("UPDATE wallets SET balance =? WHERE users_id = ?", newBalance, userID)
	c.JSON(200, gin.H{
		"status":  true,
		"message": "Order Cancelled",
	})
}
