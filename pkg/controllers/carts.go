package controllers

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jsakash/ecommers/pkg/database"
	"github.com/jsakash/ecommers/pkg/models"
)

func AddToCart(c *gin.Context) {

	UsersID := c.GetUint("id")
	// var Subtotal int
	//var GrandTotal int
	var body struct {
		ProductsID uint
		Quantity   int
	}
	c.Bind(&body)

	var stock int
	database.DB.Raw("SELECT stock FROM products WHERE id = ?", body.ProductsID).Scan(&stock)
	if stock == 0 {
		c.JSON(400, gin.H{
			"message": "Product is Out Of Stock",
		})
		return
	}

	cart := models.Cart{
		UsersId:    UsersID,
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

	for i := 1; i <= body.Quantity; i++ {
		var discount models.Discount
		database.DB.Where("product_id = ?", body.ProductsID).Find(&discount)
		deductdiscount := discount.DiscountPercentage
		cartInfo := models.CartInfo{UsersId: UsersID, ProductsID: body.ProductsID, CartsID: cart.ID, Discount: deductdiscount}
		database.DB.Create(&cartInfo)
		// c.JSON(200, gin.H{
		// 	"new id":        i.Products_id,
		// 	"new":           i.Product_Price,
		// 	"newdisc":       i.Product_Price * deductdiscount / 100,
		// 	"selling Price": i.Product_Price - (i.Product_Price * deductdiscount / 100),
		// })
		// 	}
		// }
		//OrderInfo(DiscountedAmount, deductAm, GrandTotal)

	}

}

func CartList(c *gin.Context) {
	var Subtotal int
	UsersID := c.GetUint("id")
	var cart []struct {
		ID            int
		Products_id   int
		Product_Name  string
		Brand_Name    string
		Description   string
		Product_Price int
		Quantity      int
	}

	database.DB.Find(&cart)
	database.DB.Raw("SELECT carts.id,products_id,product_name,brand_name,product_price,quantity FROM carts INNER JOIN products on carts.products_id = products.id WHERE users_id=?", UsersID).Scan(&cart)
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

}
func RemoveFromCart(c *gin.Context) {

	id := c.Param("id")
	var cart models.Cart
	var cartinfo models.CartInfo
	//database.DB.Delete(&models.Cart{}, id)
	database.DB.Raw("DELETE FROM carts WHERE id=?", id).Scan(&cart)
	database.DB.Raw("DELETE FROM cart_infos WHERE carts_id=?", id).Scan(&cartinfo)

	c.JSON(200, gin.H{
		"message": "Deleted succesfully",
	})
}

func CartCheckoutDetails(c *gin.Context) {
	UsersID := c.GetUint("id")
	couponcode := c.Query("CouponCode")
	ApplyWallet := c.Query("ApplyWallet")
	var TotalAmpount int
	var TotalMrp int
	var checkoutinfo models.Checkoutinfo
	var balance int

	database.DB.Raw("SELECT balance FROM wallets WHERE users_id = ?", UsersID).Scan(&balance)
	c.JSON(200, gin.H{
		"Balance": balance,
	})

	var cartinfo []struct {
		Product_Name  string
		Brand_Name    string
		Product_Price int
		Quantity      int
		Discount      int
	}

	database.DB.Find(&cartinfo)
	database.DB.Raw("SELECT products_id,product_name,brand_name,product_price,discount FROM cart_infos INNER JOIN products on cart_infos.products_id = products.id WHERE users_id=?", UsersID).Scan(&cartinfo)
	c.JSON(200, gin.H{
		"Products": cartinfo,
	})
	var discountAmount int
	for _, i := range cartinfo {
		discountAmount = discountAmount + i.Product_Price*i.Discount/100
		c.JSON(200, gin.H{
			"Product Name":  i.Product_Name,
			"Product Price": i.Product_Price,
			"Discount":      i.Product_Price * i.Discount / 100,
			"selling Price": i.Product_Price - (i.Product_Price * i.Discount / 100),
		})

		TotalMrp = TotalMrp + i.Product_Price
	}
	TotalAmpount = TotalMrp - discountAmount
	if couponcode == "" {
		c.JSON(200, gin.H{
			"Total MRP": TotalMrp,
			"Discount":  discountAmount,
			//"Coupon Discount": STotalAmpount,
			"Total":         TotalAmpount,
			"Wallet Amount": balance,
		})
		//return
	}
	if ApplyWallet == "apply" {

		if balance > TotalAmpount {
			TotalAmpount = 0
			newBalance := balance - TotalAmpount
			database.DB.Model(&models.Wallet{}).Where("users_id = ?", UsersID).Update("balance", newBalance)
		} else if balance < TotalAmpount {
			TotalAmpount = TotalAmpount - balance
			newBalance := 0
			database.DB.Model(&models.Wallet{}).Where("users_id = ?", UsersID).Update("balance", newBalance)
		}
	}

	if couponcode == "" {
		database.DB.First(&checkoutinfo)
		checkoutinfo.UsersID = int(UsersID)
		checkoutinfo.Discount = discountAmount
		checkoutinfo.CouponDiscount = 0
		checkoutinfo.CouponCode = "Not Applied"
		checkoutinfo.TotalMrp = TotalMrp
		checkoutinfo.Total = TotalAmpount
		database.DB.Save(&checkoutinfo)

		c.JSON(200, gin.H{
			"messsage": "No COupon",
		})
		return
	}

	var coupon models.Coupon
	database.DB.Raw("SELECT coupon_percentage FROM coupons WHERE coupon_code = ?", couponcode).Find(&coupon)
	STotalAmpount := TotalMrp * coupon.CouponPercentage / 100
	TotalAmpount = TotalAmpount - STotalAmpount
	c.JSON(200, gin.H{
		"Total MRP":       TotalMrp,
		"Discount":        discountAmount,
		"Coupon Discount": STotalAmpount,
		"Total":           TotalAmpount,
		"Wallet amount":   balance,
	})

	database.DB.First(&checkoutinfo)
	checkoutinfo.UsersID = int(UsersID)
	checkoutinfo.Discount = discountAmount
	checkoutinfo.CouponDiscount = STotalAmpount
	checkoutinfo.CouponCode = couponcode
	checkoutinfo.TotalMrp = TotalMrp
	checkoutinfo.Total = TotalAmpount
	database.DB.Save(&checkoutinfo)
}

func CartCheckout(c *gin.Context) {
	userID := c.GetUint("id")
	PaymentMethod := c.Query("PaymentMethod")
	address := c.Query("AddressID")
	addressID, _ := strconv.Atoi(address)
	cod := "COD"
	var orderID string
	// Making random order id
	rand.Seed(time.Now().UnixNano())
	value := rand.Intn(999999-100000) + 100000
	id := strconv.Itoa(value)
	orderID = "OID" + id
	var checkoutinfo models.Checkoutinfo
	database.DB.First(&checkoutinfo)
	c.JSON(200, gin.H{
		"message": checkoutinfo,
	})
	discount := checkoutinfo.Discount
	couponDiscount := checkoutinfo.CouponDiscount
	couponCode := checkoutinfo.CouponCode
	total := checkoutinfo.Total

	c.JSON(200, gin.H{
		//"message":    checkoutinfo,
		"discount":   discount,
		"coupondisc": couponDiscount,
		"coupcod":    couponCode,
		"tot":        total,
	})

	if PaymentMethod == "COD" {
		orders := models.Orders{UsersID: userID, AddressID: uint(addressID), OrderID: orderID, Discount: discount, CouponDiscount: couponDiscount, CouponCode: couponCode, Payment_Method: cod, Total_Amount: total, Status: true}
		result := database.DB.Create(&orders)
		if result.Error != nil {
			c.JSON(400, gin.H{
				"message": "error",
			})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{
			"message": "Success",
		})
		placeOrder(userID, orderID)
	}

}

func placeOrder(uid uint, oid string) {
	usersID := uid
	orderId := oid
	status := "Confirmed"
	var cartinfo []struct {
		ProductsID    uint
		Product_Name  string
		Brand_Name    string
		Product_Price int
		Quantity      int
		Discount      int
	}
	var stock int
	var product models.Products
	database.DB.Find(&cartinfo)
	database.DB.Raw("SELECT products_id,product_name,brand_name,product_price,discount FROM cart_infos INNER JOIN products on cart_infos.products_id = products.id WHERE users_id=?", usersID).Scan(&cartinfo)
	// c.JSON(200, gin.H{
	// 	"Products": cartinfo,
	// })
	//var discountAmount int
	for _, i := range cartinfo {
		Pid := i.ProductsID
		Pname := i.Product_Name
		price := i.Product_Price
		Sprice := i.Product_Price - (i.Product_Price * i.Discount / 100)

		Disc := i.Product_Price * i.Discount / 100
		ordereditems := models.Ordereditems{UsersID: usersID, ProductsID: Pid, Order_ID: orderId, Product_Name: Pname, Price: price, SellingPrice: Sprice, Discount: Disc, Status: status}
		database.DB.Create(&ordereditems)
		database.DB.Raw("SELECT stock FROM products WHERE id = ?", Pid).Scan(&stock)
		database.DB.Raw("UPDATE products SET stock = ? WHERE id = ?", stock-1, Pid).Scan(&product)

	}

	var cart models.Cart
	database.DB.Raw("DELETE FROM carts WHERE users_id = ?", usersID).Scan(&cart)
	var cartinfodel models.CartInfo
	database.DB.Raw("DELETE FROM cart_infos WHERE users_id = ?", usersID).Scan(&cartinfodel)
}

func OrderedItems(c *gin.Context) {
	UsersID := c.GetUint("id")
	var items []models.Ordereditems
	database.DB.Where("users_id = ?", UsersID).Find(&items)

	// c.JSON(200, gin.H{
	// 	"message": items,
	// })

	for _, i := range items {
		c.JSON(200, gin.H{
			"id":            i.ID,
			"Product Price": i.Price,
			"ProductsID":    i.ProductsID,
			"Order_ID":      i.Order_ID,
			"Product_Name":  i.Product_Name,
			"SellingPrice":  i.SellingPrice,
			"Discount":      i.Discount,
			"Order Status":  i.Status,
		})
	}
}

func CancelOrder(c *gin.Context) {
	userID := c.GetUint("id")
	var items models.Ordereditems
	var updateStatus string = "Cancelled"
	id := c.Query("ProductID")

	database.DB.First(&items, id)
	database.DB.Model(&items).Where("id=?", id).Update("status", updateStatus)

	c.JSON(200, gin.H{
		"message": "Order Cancelled",
	})

	var price int
	database.DB.Raw("SELECT selling_price FROM ordereditems WHERE id = ?", id).Scan(&price)

	var balance int
	database.DB.Raw("SELECT balance FROM wallets WHERE users_id = ?", userID).Scan(&balance)

	newBalance := balance + price
	//var wallet models.Wallet
	database.DB.Model(&models.Wallet{}).Where("users_id = ?", userID).Update("balance", newBalance)
	//database.DB.Raw("UPDATE wallets SET balance =? WHERE users_id = ?", newBalance, userID)

	c.JSON(200, gin.H{
		"klm": newBalance,
	})
}

func WalletBalance(c *gin.Context) {

	userID := c.GetUint("id")
	var balance int

	database.DB.Raw("SELECT balance FROM wallets WHERE users_id = ?", userID).Scan(&balance)
	c.JSON(200, gin.H{
		"Balance": balance,
	})

}
