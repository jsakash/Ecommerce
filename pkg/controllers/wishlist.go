package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jsakash/ecommers/pkg/database"
	"github.com/jsakash/ecommers/pkg/models"
)

func AddToWishlist(c *gin.Context) {
	id := c.GetUint("id")
	var body struct {
		ProductsID uint
	}
	c.Bind(&body)

	wishlist := models.Wishlist{
		UsersID:    id,
		ProductsID: body.ProductsID,
	}

	result := database.DB.Create(&wishlist)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Error",
		})
	}
	// Return it
	c.JSON(200, gin.H{
		"message": "Added To Wishlist",
	})
}

func Wishlist(c *gin.Context) {

	UsersID := c.GetUint("id")
	var wishlist []models.Wishlist

	database.DB.Where("wishlists.users_id = ?", UsersID).Find(&wishlist)

	for _, i := range wishlist {
		c.JSON(200, gin.H{
			"UserId":    i.UsersID,
			"ProductId": i.ProductsID,
			"ID":        i.ID,
		})
	}
}
