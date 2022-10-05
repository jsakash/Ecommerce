package controllers

import (
	"fmt"
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
	var body struct {
		ProductsID uint
		Quantity   int
	}
	c.Bind(&body)

	var stock int
	database.DB.Raw("SELECT stock FROM products WHERE id = ?", body.ProductsID).Scan(&stock)
	if stock == 0 {
		c.JSON(400, gin.H{
			"status":  false,
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
		"status":  true,
		"message": "Added To Cart",
	})

	for i := 1; i <= body.Quantity; i++ {
		var discount models.Discount
		database.DB.Where("product_id = ?", body.ProductsID).Find(&discount)
		deductdiscount := discount.DiscountPercentage
		cartInfo := models.CartInfo{UsersId: UsersID, ProductsID: body.ProductsID, Discount: deductdiscount}
		database.DB.Create(&cartInfo)
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
		"status":   true,
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
	/// checkoutinfo models.Checkoutinfo
	var balance int

	var checkoutinfo models.Checkoutinfo
	database.DB.First(&checkoutinfo, "users_id = ?", UsersID)

	if checkoutinfo.ID == 0 {
		checkInfo := models.Checkoutinfo{UsersID: int(UsersID)}
		database.DB.Create(&checkInfo)
	}

	// Making random order id
	rand.Seed(time.Now().UnixNano())
	value := rand.Intn(9999999999-1000000000) + 1000000000
	id := strconv.Itoa(value)
	orderID := "OID" + id

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
		var checkInfo models.Checkoutinfo
		database.DB.Raw("UPDATE checkoutinfos SET order_id = ?,discount = ?,coupon_discount = ?,coupon_code = ?,total_mrp = ?,total = ? WHERE users_id = ?", orderID, discountAmount, 0, "Not Applied", TotalMrp, TotalAmpount, UsersID).Scan(&checkInfo)
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

	if ApplyWallet == "apply" {

		if balance > TotalAmpount {
			newBalance := balance - TotalAmpount
			database.DB.Model(&models.Wallet{}).Where("users_id = ?", UsersID).Update("balance", newBalance)
			WalletHistory := models.Wallethistory{UsersID: UsersID, Debit: TotalAmpount, Credit: 0}
			database.DB.Create(&WalletHistory)
			TotalAmpount = 0
		} else if balance < TotalAmpount {
			TotalAmpount = TotalAmpount - balance
			newBalance := 0
			database.DB.Model(&models.Wallet{}).Where("users_id = ?", UsersID).Update("balance", newBalance)
			WalletHistory := models.Wallethistory{UsersID: UsersID, Debit: balance, Credit: 0}
			database.DB.Create(&WalletHistory)
		}
	}

	var checkInfo models.Checkoutinfo
	database.DB.Raw("UPDATE checkoutinfos SET order_id = ?,discount = ?,coupon_discount = ?,coupon_code = ?,total_mrp = ?,total = ? WHERE users_id = ?", orderID, discountAmount, STotalAmpount, couponcode, TotalMrp, TotalAmpount, UsersID).Scan(&checkInfo)
	c.JSON(200, gin.H{
		"messsage":    "Coupon Applied",
		"TotalAmount": TotalAmpount,
	})

}

func CartCheckout(c *gin.Context) {
	userID := c.GetUint("id")
	PaymentMethod := c.Query("PaymentMethod")
	address := c.Query("AddressID")
	addressID, _ := strconv.Atoi(address)
	cod := "COD"
	RazorPay := "RAZORPAY"
	PayStatus := "PENDING"
	OrderStatus := "PENDING"
	var orderID string

	// Fetching checkout informations
	var checkoutinfo models.Checkoutinfo
	database.DB.First(&checkoutinfo, "users_id = ?", userID)
	c.JSON(200, gin.H{
		"message": checkoutinfo,
	})
	orderID = checkoutinfo.OrderID
	discount := checkoutinfo.Discount
	couponDiscount := checkoutinfo.CouponDiscount
	couponCode := checkoutinfo.CouponCode
	total := checkoutinfo.Total

	c.JSON(200, gin.H{
		//"message":    checkoutinfo,
		"discount":   discount,
		"coupondisc": couponDiscount,
		"coupcod":    couponCode,
		"total":      total,
	})

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

	database.DB.Raw("SELECT products_id,product_name,brand_name,product_price,discount FROM cart_infos INNER JOIN products on cart_infos.products_id = products.id WHERE users_id=?", userID).Scan(&cartinfo)

	for _, i := range cartinfo {
		Pid := i.ProductsID
		Pname := i.Product_Name
		Sprice := i.Product_Price - (i.Product_Price * i.Discount / 100)
		fmt.Println(Sprice)
		Cprice := total
		Disc := i.Product_Price * i.Discount / 100
		ordereditems := models.Ordereditems{UsersID: userID, ProductsID: Pid, Order_ID: orderID, Product_Name: Pname, Price: Cprice, CouponDiscount: couponDiscount, Discount: Disc, OrderStatus: "CONFIMRED", PaymentStatus: OrderStatus, Payment_Method: PaymentMethod}
		database.DB.Create(&ordereditems)
		database.DB.Raw("SELECT stock FROM products WHERE id = ?", Pid).Scan(&stock)
		database.DB.Raw("UPDATE products SET stock = ? WHERE id = ?", stock-1, Pid).Scan(&product)
	}

	if PaymentMethod == cod {
		//OrderStatus = "CONFIRMED"
		orders := models.Orders{UsersID: userID, AddressID: uint(addressID), OrderID: orderID, Discount: discount, CouponDiscount: couponDiscount, CouponCode: couponCode, Payment_Method: cod, Total_Amount: total, OrderStatus: OrderStatus, PaymentStatus: PayStatus}
		result := database.DB.Create(&orders)
		if result.Error != nil {
			c.JSON(400, gin.H{
				"message": "error",
			})
			return
		}
		var cart models.Cart
		database.DB.Raw("DELETE FROM carts WHERE users_id = ?", userID).Scan(&cart)
		var cartinfodel models.CartInfo
		database.DB.Raw("DELETE FROM cart_infos WHERE users_id = ?", userID).Scan(&cartinfodel)
		var checkinfo models.Checkoutinfo
		database.DB.Raw("DELETE FROM checkoutinfos WHERE users_id = ?", userID).Scan(&checkinfo)
		c.JSON(http.StatusAccepted, gin.H{
			"status":  true,
			"message": "Order Placed",
		})
		//placeOrder(userID, orderID)
	}
	if PaymentMethod == RazorPay {
		orders := models.Orders{UsersID: userID, AddressID: uint(addressID), OrderID: orderID, Discount: discount, CouponDiscount: couponDiscount, CouponCode: couponCode, Payment_Method: RazorPay, Total_Amount: total, OrderStatus: OrderStatus, PaymentStatus: PayStatus}
		database.DB.Create(&orders)
	}

}

func placeOrder(uid uint, oid string) {
	usersID := uid
	orderId := oid
	status := "CONFIRMED"
	var order models.Orders
	database.DB.Raw("UPDATE orders SET order_status = ? ,payment_status = ? WHERE order_id = ?", status, status, orderId).Scan(&order)
	var orderitem models.Ordereditems
	database.DB.Raw("UPDATE ordereditems SET order_status = ?,payment_status = ? WHERE order_id = ?", status, status, orderId).Scan(&orderitem)
	var cart models.Cart
	database.DB.Raw("DELETE FROM carts WHERE users_id = ?", usersID).Scan(&cart)
	var cartinfodel models.CartInfo
	database.DB.Raw("DELETE FROM cart_infos WHERE users_id = ?", usersID).Scan(&cartinfodel)
	var checkinfo models.Checkoutinfo
	database.DB.Raw("DELETE FROM checkoutinfos WHERE users_id = ?", usersID).Scan(&checkinfo)
}

func WalletBalance(c *gin.Context) {

	userID := c.GetUint("id")
	var balance int
	database.DB.Raw("SELECT balance FROM wallets WHERE users_id = ?", userID).Scan(&balance)
	c.JSON(200, gin.H{
		"status":  true,
		"Balance": balance,
		"UserID":  userID,
	})
}
