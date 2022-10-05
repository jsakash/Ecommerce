package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jsakash/ecommers/pkg/database"
	"github.com/jsakash/ecommers/pkg/models"
)

func AddCoupon(c *gin.Context) {

	couponName := c.Query("couponName")
	couponCode := c.Query("couponCode")
	Percentage := c.Query("couponPercentage")
	couponPercentage, _ := strconv.Atoi(Percentage)

	coupon := models.Coupon{CouponName: couponName, CouponCode: couponCode, CouponPercentage: couponPercentage}

	var checkCoup []models.Coupon
	database.DB.Find(&checkCoup)

	// Checking username existence
	for _, i := range checkCoup {
		if i.CouponName == couponName {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Coupon Name Already Exist",
			})
			return
		}
	}

	result := database.DB.Create(&coupon)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Error Creating Coupon",
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  true,
		"message": "Coupon Crearted",
	})

}

func DeleteCoupon(c *gin.Context) {
	var coupon models.Coupon
	couponName := c.Query("couponName")
	database.DB.Where("coupon_name = ?", couponName).Delete(&coupon)
	//database.DB.Raw("DELETE FROM coupons WHERE coupon_name=?", couponName).Scan(&coupon)
	c.JSON(200, gin.H{
		"status":  true,
		"message": "Deleted succesfully",
	})

}

func ListCoupons(c *gin.Context) {
	var coupons []models.Coupon
	result := database.DB.Find(&coupons)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"ststus":  false,
			"message": "No coupon found",
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  true,
		"message": coupons,
	})

}
