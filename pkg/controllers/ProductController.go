package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jsakash/ecommers/pkg/database"
	"github.com/jsakash/ecommers/pkg/models"
)

func AddProduct(c *gin.Context) {

	// Get data off req body

	var body struct {
		Product_Name  string
		Brand_Name    string
		Product_Price int
		Description   string
		CategoryID    uint
		ColorsID      uint
		SizeID        uint
		Stock         int
	}

	c.Bind(&body)

	// Create

	products := models.Products{
		Product_Name:  body.Product_Name,
		Brand_Name:    body.Brand_Name,
		Product_Price: body.Product_Price,
		Description:   body.Description,
		CategoryID:    body.CategoryID,
		ColorsID:      body.ColorsID,
		SizeID:        body.SizeID,
		Stock:         body.Stock,
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

	c.JSON(200, gin.H{
		"message": "Product added",
	})

}

func ListAllProducts(c *gin.Context) {

	// var products []models.Products

	// database.DB.Find(&products)

	// for _, i := range products {
	// 	c.JSON(200, gin.H{
	// 		"Product Name": i.Product_Name,
	// 		"Brand Name":   i.Brand_Name,
	// 		"Description":  i.Description,
	// 		"id":           i.ID,
	// 	})

	// }

	//db := database.GetDb()
	var products []struct {
		Product_Name  string
		Brand_Name    string
		Description   string
		Product_Price int
		Category      string
		Color         string
		Size          int
		Stock         int
	}

	database.DB.Find(&products)

	database.DB.Raw("SELECT product_name,brand_name,description,product_price,category,color,size,stock FROM products INNER JOIN categories on products.category_id = categories.id INNER JOIN colors ON products.colors_id = colors.id INNER JOIN sizes ON products.size_id = sizes.id").Scan(&products)
	c.JSON(200, gin.H{
		"Products": products,
	})

}

func UpdateProduct(c *gin.Context) {

	// Get id off the url
	id := c.Param("id")

	// Get data off req body

	var body struct {
		Product_Name  string
		Brand_Name    string
		Product_Price int
		Description   string
		CategoryID    uint
		ColorsID      uint
		SizeID        uint
		Stock         int
	}

	c.Bind(&body)

	// Find the post we are updating
	var products []models.Products
	database.DB.First(&products, id)

	// Update it
	database.DB.Model(&products).Updates(models.Products{
		Product_Name:  body.Product_Name,
		Brand_Name:    body.Brand_Name,
		CategoryID:    body.CategoryID,
		Product_Price: body.Product_Price,
		Description:   body.Description,
		ColorsID:      body.ColorsID,
		SizeID:        body.SizeID,
		Stock:         body.Stock,
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

func AddCategory(c *gin.Context) {

	// Get data off req body
	var body struct {
		Category string
	}

	c.Bind(&body)

	// Create
	catogory := models.Category{
		Category: body.Category,
	}

	result := database.DB.Create(&catogory)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Error",
		})
	}
	// Return it
	c.JSON(200, gin.H{
		"message": "Category added",
	})

}

func CatogoryList(c *gin.Context) {

	var catogory []models.Category

	database.DB.Find(&catogory)
	for _, i := range catogory {
		c.JSON(200, gin.H{
			"catogory": i.Category,
			"id":       i.ID,
		})

	}
}
func AddColor(c *gin.Context) {

	// Get data off req body
	var body struct {
		Color string
		//ProductsID Products
	}
	c.Bind(&body)

	// Create
	color := models.Colors{
		Color: body.Color,
	}
	result := database.DB.Create(&color)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Error",
		})
	}
	// Return it
	c.JSON(200, gin.H{
		"message": "Color Added Successfully",
	})

}

func AddSize(c *gin.Context) {

	// Get data off req body
	var body struct {
		Size int
	}

	c.Bind(&body)

	// Create
	size := models.Size{
		Size: body.Size,
	}
	result := database.DB.Create(&size)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Error",
		})
	}
	// Return it
	c.JSON(200, gin.H{
		"message": "Size added",
	})
}
