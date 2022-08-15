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
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{

		"token": tokenString,
	})

}

func AdminValidate(c *gin.Context) {
	user, _ := c.Get("admin")

	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})

}

func ListAllUsers(c *gin.Context) {

	var users []models.Users

	database.DB.Find(&users)

	c.JSON(200, gin.H{
		"users": users,
	})

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

func AddProduct(c *gin.Context) {

	// Get data off req body

	var body struct {
		Product_Name     string
		Brand_Name       string
		Product_Catogory string
		Description      string
		Color            string
		Size             int
		Stock            int
	}

	c.Bind(&body)

	// Create

	products := models.Products{
		Product_Name:     body.Product_Name,
		Brand_Name:       body.Brand_Name,
		Product_Catogory: body.Product_Catogory,
		Description:      body.Description,
		Color:            body.Color,
		Size:             body.Size,
		Stock:            body.Stock,
	}

	//var duplicate []models.Products
	//database.DB.Find(&duplicate)

	// for _, i := range duplicate {
	// 	if i. == body {
	// 		c.JSON(400, gin.H{
	// 			"message": "Product already exists",
	// 		})
	// 		return
	// 	}
	// }

	result := database.DB.Create(&products)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Error",
		})
	}
	// Return it

	c.Status(200)

}

func ListAllProducts(c *gin.Context) {

	var products []models.Products

	database.DB.Find(&products)

	c.JSON(200, gin.H{
		"users": products,
	})

}

func UpdateProduct(c *gin.Context) {

	// Get id off the url
	id := c.Param("id")

	// Get data off req body

	var body struct {
		Product_Name     string
		Brand_Name       string
		Product_Catogory string
		Description      string
		Color            string
		Size             int
		Stock            int
	}

	c.Bind(&body)

	// Find the post we are updating
	var products []models.Products
	database.DB.First(&products, id)

	// Update it
	database.DB.Model(&products).Updates(models.Products{
		Product_Name:     body.Product_Name,
		Brand_Name:       body.Brand_Name,
		Product_Catogory: body.Product_Catogory,
		Description:      body.Description,
		Color:            body.Color,
		Size:             body.Size,
		Stock:            body.Stock,
	})

	// Response
	c.JSON(200, gin.H{
		"message": "Product updated",
	})

}

func FetchProduct(c *gin.Context) {

	// Get id off the url
	id := c.Param("id")

	// Find the product
	var product models.Products
	database.DB.First(&product, id)

	c.JSON(200, gin.H{
		"product": product,
	})

}

func DeleteProduct(c *gin.Context) {

	id := c.Param("id")
	database.DB.Delete(&models.Products{}, id)

	c.JSON(200, gin.H{
		"message": "Deleted succesfully",
	})

}
