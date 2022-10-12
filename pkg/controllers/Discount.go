package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jsakash/ecommers/pkg/database"
	"github.com/jsakash/ecommers/pkg/models"
)

func AddDiscount(c *gin.Context) {
	// Get Info off req body
	discountName := c.PostForm("discountName")
	DPercentage := c.PostForm("discountPercentage")
	discountPercentage, _ := strconv.Atoi(DPercentage)
	PId := c.Query("productId")
	productId, _ := strconv.Atoi(PId)

	discount := models.Discount{DiscountName: discountName, DiscountPercentage: discountPercentage, ProductId: productId}
	var checkDisc []models.Discount
	database.DB.Find(&checkDisc)

	// Checking username existence
	for _, i := range checkDisc {
		if i.DiscountName == discountName {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Discount Name Already Exist",
			})
			return
		}
	}
	result := database.DB.Create(&discount)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Error Creating Coupon",
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  true,
		"message": "Discount Created",
	})

}

func DeleteDiscount(c *gin.Context) {
	var discount models.Discount
	discountName := c.Query("discountName")
	database.DB.Where("coupon_name = ?", discountName).Delete(&discount)
	//database.DB.Raw("DELETE FROM coupons WHERE coupon_name=?", couponName).Scan(&coupon)
	c.JSON(200, gin.H{
		"status":  true,
		"message": "Deleted succesfully",
	})
}
func ListDiscount(c *gin.Context) {

	var discount []models.Discount
	result := database.DB.Find(&discount)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"status":  false,
			"message": "No discount found",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": true,
		"data":   discount,
	})

}
