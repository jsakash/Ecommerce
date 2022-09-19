package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jsakash/ecommers/pkg/database"
	"github.com/jsakash/ecommers/pkg/models"
)

func ListAllUsers(c *gin.Context) {

	var users []models.Users
	database.DB.Find(&users)
	for _, i := range users {
		c.JSON(200, gin.H{
			"Name":  i.First_Name,
			"Email": i.Email,
			"id":    i.ID,
		})
	}
}

func BlockUser(c *gin.Context) {

	var user models.Users
	var updateStatus bool = false
	id := c.Param("id")

	database.DB.First(&user, id)
	database.DB.Model(&user).Where("id=?", id).Update("status", updateStatus)

	c.JSON(200, gin.H{
		"message": "User Blocked",
	})
}

func UnBlockUser(c *gin.Context) {

	var user models.Users
	var updateStatus bool = true
	id := c.Param("id")

	database.DB.First(&user, id)
	database.DB.Model(&user).Where("id=?", id).Update("status", updateStatus)

	c.JSON(200, gin.H{
		"message": "User UnBlocked",
	})

}

func DeleteUser(c *gin.Context) {

	id := c.Param("id")
	database.DB.Delete(&models.Users{}, id)

	c.JSON(200, gin.H{
		"message": "Deleted succesfully",
	})

}
