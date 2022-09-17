package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jsakash/ecommers/pkg/database"
	"github.com/jsakash/ecommers/pkg/models"
)

var stot int
var GrandTotal int

func AddToCart(c *gin.Context) {

	id := c.GetUint("id")
	var body struct {
		ProductsID uint
		Quantity   int
	}
	c.Bind(&body)

	cart := models.Cart{
		UsersId:    id,
		ProductsId: body.ProductsID,
		Quantity:   body.Quantity,
	}

	result := database.DB.Create(&cart)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Error",
		})

		return
	}
	// Return it
	c.JSON(200, gin.H{
		"message": "Added To Cart",
	})
}

func CartList(c *gin.Context) {
	var Subtotal int
	UsersID := c.GetUint("id")
	confirmaton := c.Query("confirm")
	check := "checkout"
	var cart []struct {
		Products_id   int
		Product_Name  string
		Brand_Name    string
		Description   string
		Product_Price int
		Quantity      int
	}

	database.DB.Find(&cart)
	database.DB.Raw("SELECT products_id,product_name,brand_name,product_price,quantity FROM carts INNER JOIN products on carts.products_id = products.id WHERE users_id=?", UsersID).Scan(&cart)
	c.JSON(200, gin.H{
		"Products": cart,
	})

	for _, i := range cart {
		sum := i.Product_Price * i.Quantity
		Subtotal = Subtotal + sum
	}
	c.JSON(200, gin.H{
		"Subtotal": Subtotal,
	})

	if confirmaton == check {
		Checkout(Subtotal)
		c.JSON(200, gin.H{
			"mdg": stot,
		})
	} else {
		return
	}
}
func RemoveFromCart(c *gin.Context) {

	id := c.Param("id")
	database.DB.Delete(&models.Cart{}, id)

	c.JSON(200, gin.H{
		"message": "Deleted succesfully",
	})
}

func Checkout(tot int) {

	stot = tot

}

func CartCheckoutDetails(c *gin.Context) {
	UsersID := c.GetUint("id")

	var cart []struct {
		Products_id   int
		Product_Name  string
		Brand_Name    string
		Product_Price int
		Quantity      int
	}

	database.DB.Find(&cart)
	database.DB.Raw("SELECT products_id,product_name,brand_name,product_price,quantity FROM carts INNER JOIN products on carts.products_id = products.id WHERE users_id=?", UsersID).Scan(&cart)
	c.JSON(200, gin.H{
		"Products": cart,
	})

	var gtax models.Tax
	//database.DB.Where("category = ?", "general_tax").First(&tax)
	database.DB.Where("category =?", "general_tax").Find(&gtax)
	Generaltax := gtax.Tax

	// c.JSON(200, gin.H{
	// 	"msg": TaxReduction,
	// })
	var otax models.Tax
	database.DB.Where("category =?", "other_taxes").Find(&otax)
	Othertax := otax.Tax
	GeneralTaxAddition := stot * Generaltax / 100
	OtherTaxAddition := stot * Othertax / 100
	GrandTotal = stot + GeneralTaxAddition + OtherTaxAddition
	c.JSON(200, gin.H{
		"General Tax":  GeneralTaxAddition,
		"Other Tax":    OtherTaxAddition,
		"Total Amount": GrandTotal,
	})

	for _, i := range cart {
		c.JSON(200, gin.H{
			"Products": i.Products_id,
		})
	}

}
