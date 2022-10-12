package controllers

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jsakash/ecommers/pkg/database"
	"github.com/jsakash/ecommers/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}

func Signup(c *gin.Context) {

	var body struct {
		First_Name   string
		Last_Name    string
		Email        string
		Password     string
		Phone_Number string
		Status       bool
	}

	//Get the name/email/password off req body
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Hash the given password
	hashPass := HashPassword(body.Password)

	// Create the user
	user := models.Users{First_Name: body.First_Name, Last_Name: body.Last_Name, Email: body.Email, Password: hashPass, Phone_Number: body.Phone_Number, Status: true}
	var checkMail []models.Users
	database.DB.Find(&checkMail)

	// Checking username existence
	for _, i := range checkMail {
		if i.Email == user.Email {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Username Already Exist",
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
	c.JSON(http.StatusOK, gin.H{
		"message": "Account Created",
	})

	var users models.Users
	database.DB.First(&users, "email = ?", body.Email)

	wallet := models.Wallet{UsersID: users.ID, Balance: 0}
	database.DB.Create(&wallet)

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
			"status":  false,
			"message": "Invalid email or password",
		})

		return
	}

	// Compare sent in pass with saved user pass hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid password",
		})
		return
	}

	if !user.Status {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "You are blocked",
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
	c.SetCookie("UserAuthorization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{

		"status":  true,
		"message": "OK",
		"data":    tokenString,
	})

}

func ChangePasswors(c *gin.Context) {

	var body struct {
		Password    string
		newPassword string
	}
	userID := c.GetUint("id")
	password := body.Password
	newPassword := body.newPassword

	// Look up request user
	var user models.Users
	database.DB.First(&user, "id = ?", userID)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid user password",
		})

		return
	}

	// Compare sent in pass with saved user pass hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid password",
		})
		return
	}

	newHashPass := HashPassword(newPassword)
	database.DB.Model(&user).Where("id=?", userID).Update("password", newHashPass)
	// Return Response
	c.JSON(200, gin.H{
		"status":  true,
		"message": "Password changed successfully",
	})

}

// to test middleware
func Validate(c *gin.Context) {
	//user := c.GetInt("id")
	check, _ := c.Get("user")

	id := c.GetUint("id")

	c.JSON(http.StatusOK, gin.H{
		"message": id,
		"user":    check,
	})

}

func AddAddress(c *gin.Context) {

	id := c.GetUint("id") // Get ID from token
	var body struct {
		Name         string
		Phone_number int
		Pincode      int
		House_Adress string
		Area         string
		Landmark     string
		City         string
	}

	c.Bind(&body)

	// Create
	address := models.Address{
		UsersID:      id,
		Name:         body.Name,
		Phone_number: body.Phone_number,
		Pincode:      body.Pincode,
		House_Adress: body.House_Adress,
		Area:         body.Area,
		Landmark:     body.Landmark,
		City:         body.City,
	}

	result := database.DB.Create(&address)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"staus":   false,
			"message": "Error",
		})
	}
	// Return it

	c.JSON(200, gin.H{
		"status":  true,
		"message": "Address added sucessfully",
	})

}

func SelectAddress(c *gin.Context) {

	// Get id off the loggedin user
	UsersID := c.GetUint("id")

	// Find the product
	var address []models.Address
	// database.DB.First(&address, UsersID)
	database.DB.Where("addresses.users_id = ?", UsersID).Find(&address)

	for _, i := range address {
		c.JSON(200, gin.H{
			"Name":          i.Name,
			"Phone Number":  i.Phone_number,
			"Pincode":       i.Pincode,
			"House Address": i.House_Adress,
			"Area":          i.Area,
			"Landmark":      i.Landmark,
			"City":          i.City,
			"id":            i.ID,
		})
	}

}

func EditAddress(c *gin.Context) {

	// Get id off the url
	id := c.Query("id")

	// Get data off req body
	var body struct {
		Name         string
		Phone_number int
		Pincode      int
		House_Adress string
		Area         string
		Landmark     string
		City         string
	}

	c.Bind(&body)

	// Find the post we are updating
	var Address []models.Address
	database.DB.First(&Address, id)

	// Update it
	database.DB.Model(&Address).Updates(models.Address{
		Name:         body.Name,
		Phone_number: body.Phone_number,
		Pincode:      body.Phone_number,
		House_Adress: body.House_Adress,
		Area:         body.Area,
		Landmark:     body.Landmark,
		City:         body.City,
	})
	// Response
	c.JSON(200, gin.H{
		"ststus":  true,
		"message": "Address updated",
	})

}
