package controllers

import (
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jsakash/ecommers/pkg/database"
	"github.com/jsakash/ecommers/pkg/models"
)

func AddProduct(c *gin.Context) {

	// Get data off req body

	// var body struct {
	// 	Product_Name  string
	// 	Brand_Name    string
	// 	Product_Price int
	// 	Description   string
	// 	CategoryID    uint
	// 	ColorsID      uint
	// 	SizeID        uint
	// 	Stock         int
	// 	Cover         string
	// 	SubPic1       string
	// 	SubPic2       string
	// 	SubPic3       string
	// }

	// c.Bind(&body)

	// Create
	ProdName := c.PostForm("ProductName")
	BrandName := c.PostForm("BrandName")
	Pprice := c.PostForm("ProductPrice")
	ProductPrice, _ := strconv.Atoi(Pprice)
	Description := c.PostForm("Description")
	CatID := c.PostForm("CategoryID")
	CategoryID, _ := strconv.ParseUint(CatID, 10, 64)
	ColID := c.PostForm("ColorID")
	ColorID, _ := strconv.ParseUint(ColID, 10, 64)
	sizID := c.PostForm("SizeID")
	SizeID, _ := strconv.ParseUint(sizID, 10, 64)
	stk := c.PostForm("Stock")
	Stock, _ := strconv.Atoi(stk)

	CoverPicPath, _ := c.FormFile("cover")
	extension := filepath.Ext(CoverPicPath.Filename)
	CoverPic := uuid.New().String() + extension
	c.SaveUploadedFile(CoverPicPath, "./public/"+CoverPic)

	SubPicPath1, _ := c.FormFile("subpic1")
	extension = filepath.Ext(SubPicPath1.Filename)
	SubPic1 := uuid.New().String() + extension
	c.SaveUploadedFile(SubPicPath1, "./public/"+SubPic1)

	SubPicPath2, _ := c.FormFile("subpic2")
	extension = filepath.Ext(SubPicPath2.Filename)
	SubPic2 := uuid.New().String() + extension
	c.SaveUploadedFile(SubPicPath2, "./public/"+SubPic2)

	products := models.Products{
		Product_Name:  ProdName,
		Brand_Name:    BrandName,
		Product_Price: ProductPrice,
		Description:   Description,
		CategoryID:    uint(CategoryID),
		ColorsID:      uint(ColorID),
		SizeID:        uint(SizeID),
		Stock:         Stock,
		Cover:         CoverPic,
		SubPic1:       SubPic1,
		SubPic2:       SubPic2,
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

	color := c.Query("color")
	category := c.Query("category")
	size := c.Query("size")
	//filter := c.Query("filters")
	limit := 3
	pageStr := c.Query("page")
	page, _ := strconv.Atoi(pageStr)
	offset := (page - 1) * limit

	var products []struct {
		ID            uint
		Product_Name  string
		Brand_Name    string
		Description   string
		Product_Price int
		Category      string
		Color         string
		Size          int
		Stock         int
		Cover         string
		SubPic1       string
		SubPic2       string
	}

	if category == "" && color == "" && size == "" {

		database.DB.Raw("SELECT products.id,product_name,brand_name,description,product_price,category,color,size,stock,cover,sub_pic1,sub_pic2 FROM products INNER JOIN categories on products.category_id = categories.id INNER JOIN colors ON products.colors_id = colors.id INNER JOIN sizes ON products.size_id = sizes.id ORDER BY products.id Limit(?) OFFSET(?)", limit, offset).Scan(&products)
		c.JSON(200, gin.H{
			"Products": products,
		})
	} else if color == "" && size == "" && category != "" {
		database.DB.Raw("SELECT products.id,product_name,brand_name,description,product_price,category,color,size,stock,cover,sub_pic1,sub_pic2 FROM products INNER JOIN categories on products.category_id = categories.id INNER JOIN colors ON products.colors_id = colors.id INNER JOIN sizes ON products.size_id = sizes.id WHERE category = ? ORDER BY products.id Limit(?) OFFSET(?)", category, limit, offset).Scan(&products)
		c.JSON(200, gin.H{
			"Products": products,
		})
	} else if category == "" && size == "" && color != "" {
		database.DB.Raw("SELECT products.id,product_name,brand_name,description,product_price,category,color,size,stock,cover,sub_pic1,sub_pic2 FROM products INNER JOIN categories on products.category_id = categories.id INNER JOIN colors ON products.colors_id = colors.id INNER JOIN sizes ON products.size_id = sizes.id WHERE color = ? ORDER BY products.id Limit(?) OFFSET(?)", color, limit, offset).Scan(&products)
		c.JSON(200, gin.H{
			"Products": products,
		})
	} else if category == "" && color == "" && size != "" {
		database.DB.Raw("SELECT products.id,product_name,brand_name,description,product_price,category,color,size,stock,cover,sub_pic1,sub_pic2 FROM products INNER JOIN categories on products.category_id = categories.id INNER JOIN colors ON products.colors_id = colors.id INNER JOIN sizes ON products.size_id = sizes.id WHERE size = ? ORDER BY products.id Limit(?) OFFSET(?)", size, limit, offset).Scan(&products)
		c.JSON(200, gin.H{
			"Products": products,
		})
	} else if category == "" && color != "" && size != "" {
		database.DB.Raw("SELECT products.id,product_name,brand_name,description,product_price,category,color,size,stock,cover,sub_pic1,sub_pic2 FROM products INNER JOIN categories on products.category_id = categories.id INNER JOIN colors ON products.colors_id = colors.id INNER JOIN sizes ON products.size_id = sizes.id WHERE size = ? AND color = ? ORDER BY products.id Limit(?) OFFSET(?)", size, color, limit, offset).Scan(&products)
		c.JSON(200, gin.H{
			"Products": products,
		})
	} else if size == "" && color != "" && category != "" {
		database.DB.Raw("SELECT products.id,product_name,brand_name,description,product_price,category,color,size,stock,cover,sub_pic1,sub_pic2 FROM products INNER JOIN categories on products.category_id = categories.id INNER JOIN colors ON products.colors_id = colors.id INNER JOIN sizes ON products.size_id = sizes.id WHERE category = ? AND color = ? ORDER BY products.id Limit(?) OFFSET(?)", category, color, limit, offset).Scan(&products)
		c.JSON(200, gin.H{
			"Products": products,
		})
	} else if color == "" && size != "" && category != "" {
		database.DB.Raw("SELECT products.id,product_name,brand_name,description,product_price,category,color,size,stock,cover,sub_pic1,sub_pic2 FROM products INNER JOIN categories on products.category_id = categories.id INNER JOIN colors ON products.colors_id = colors.id INNER JOIN sizes ON products.size_id = sizes.id WHERE category = ? AND size = ? ORDER BY products.id Limit(?) OFFSET(?)", category, size, limit, offset).Scan(&products)
		c.JSON(200, gin.H{
			"Products": products,
		})
	}

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

func SearchProduct(c *gin.Context) {

	productName := c.Query("product")

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

	database.DB.Raw("SELECT product_name,brand_name,description,product_price,category,color,size,stock FROM products INNER JOIN categories on products.category_id = categories.id INNER JOIN colors ON products.colors_id = colors.id INNER JOIN sizes ON products.size_id = sizes.id WHERE product_name = ?", productName).Scan(&products)
	c.JSON(200, gin.H{
		"Products": products,
	})

}
