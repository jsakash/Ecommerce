package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jsakash/ecommers/pkg/database"
	"github.com/jsakash/ecommers/pkg/models"
)

// To create an admin with credentials stored in .env file
func CreateAdmin(c *gin.Context) {

	Email := os.Getenv("ADMIN_EMAIL")
	Password := os.Getenv("ADMIN_PASSWORD")
	admin := models.Admin{Email: Email, Password: Password}
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
			"status":  false,
			"message": "Failed to read body",
		})
		return
	}
	// Look Up the request
	var admin models.Admin
	database.DB.First(&admin, "email = ?", body.Email)

	if admin.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid Email or Password",
		})
		return
	}
	// Compare sent in passwrod with saved password
	var adminPassword models.Admin
	database.DB.Find(&adminPassword)
	if adminPassword.Password != body.Password {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Incorrect Password",
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
			"status":  false,
			"message": "Failed to create token",
		})
		return
	}
	// Sent it back
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("AdminAuthorization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "ok",
		"data":    tokenString,
	})
}
