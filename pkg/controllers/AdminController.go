package controllers

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jsakash/ecommers/pkg/database"
	"github.com/jsakash/ecommers/pkg/models"
)

func CreateAdmin(c *gin.Context) {

	admin := models.Admin{Email: "akashjs@gmail.com", Password: "akashjs"}
	result := database.DB.Create(&admin)

	if result.Error != nil {
		c.Status(400)
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"message": "Success",
	})
}

func AdminLogin(c *gin.Context) {
	// Get email and password of f req body
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	// Look Up the request
	var admin models.Admin
	database.DB.First(&admin, "email = ?", body.Email)

	if admin.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Email or Password",
		})
		return
	}
	// Compare sent in passwrod with saved password
	var adminPassword models.Admin
	database.DB.Find(&adminPassword)
	if adminPassword.Password != body.Password {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect Password",
		})
	}
	// Generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": admin.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// SIgn and get the complete encoded token as a string using the secret key
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}
	// Sent it back
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("AdminAuthorization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{

		"token": tokenString,
	})
}
func AdminValidate(c *gin.Context) {
	admin, _ := c.Get("admin")

	c.JSON(http.StatusOK, gin.H{
		"message": admin,
	})
}
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

func AddTax(c *gin.Context) {
	taxCategory := c.Query("taxcategory")
	tax := c.Query("Tax")
	taxPer, _ := strconv.Atoi(tax)
	Tax := models.Tax{Category: taxCategory, Tax: taxPer}
	result := database.DB.Create(&Tax)

	if result.Error != nil {
		c.Status(400)
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"message": "Success",
	})
}

func UpdateTax(c *gin.Context) {

	category := c.Query("category")
	newTax := c.Query("newtax")
	tax, _ := strconv.Atoi(newTax)
	var Tax models.Tax
	database.DB.Model(&Tax).Where("category = ?", category).Update("tax", tax)

	c.JSON(200, gin.H{
		"message": "Tax Updated",
	})
}
