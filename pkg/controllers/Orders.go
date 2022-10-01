package controllers

import (
	"math/rand"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func AddImage(c *gin.Context) {

	SubPicPath4, _ := c.FormFile("roompicpath4")
	extension := filepath.Ext(SubPicPath4.Filename)
	SubPic4 := uuid.New().String() + extension
	c.SaveUploadedFile(SubPicPath4, "./public/"+SubPic4)

	var addroom models.Addimage

	addroom.Cover = SubPic4
	database.DB.Select("cover").Create(&addroom)

}
