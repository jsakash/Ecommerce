package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jsakash/ecommers/pkg/database"
	"github.com/jsakash/ecommers/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

var body struct {
	First_Name string
	Last_Name  string
	Email      string
	Password   string
	Status     bool
}

func Signup(c *gin.Context) {

	//Get the name/email/password off req body

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	// Hash the password

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
	}

	// Create the user
	user := models.Users{First_Name: body.First_Name, Last_Name: body.Last_Name, Email: body.Email, Password: string(hash), Status: true}

	var checkMail []models.Users
	database.DB.Find(&checkMail)

	// Checking username existence

	for _, i := range checkMail {
		if i.Email == user.Email {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Username Already Exist",
			})

			return
		}
	}

	if user.First_Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Name is required",
		})

		return

	}

	if user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email is required",
		})

		return

	}

	if user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password is required",
		})

		return

	}

	result := database.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})

		return
	}

	// Respond

	c.JSON(http.StatusOK, gin.H{})

}

func Login(c *gin.Context) {

	// Get Email and Password off req body
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body ",
		})

		return
	}

	// Look up request user
	var user models.Users
	database.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})

		return
	}

	// Compare sent in pass with saved user pass hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})

		return
	}

	if !user.Status {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "You are blocked ",
		})
		return
	}

	// Generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
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
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": user,
		"token":   tokenString,
	})

}

func Validate(c *gin.Context) {

	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})

}
