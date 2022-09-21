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

func BuyNow(c *gin.Context) {
	UsersID := c.GetUint("id")
	ProductId := c.Query("PrductId")
	ApplyWallet := c.Query("ApplyWallet")
	couponcode := c.Query("CouponCode")
	var balance int
	var discountAmount int
	var TotalMrp int
	var TotalAmount int
	var product models.Products
	database.DB.First(&product, "id = ?", ProductId)

	var checkoutinfo models.Checkoutinfo
	database.DB.First(&checkoutinfo, "users_id = ?", UsersID)

	if checkoutinfo.ID == 0 {
		checkInfo := models.Checkoutinfo{UsersID: int(UsersID)}
		database.DB.Create(&checkInfo)
	}

	// Making random order id
	rand.Seed(time.Now().UnixNano())
	value := rand.Intn(999999-100000) + 100000
	id := strconv.Itoa(value)
	orderID := "OID" + id

	database.DB.Raw("SELECT balance FROM wallets WHERE users_id = ?", UsersID).Scan(&balance)
	c.JSON(200, gin.H{
		"Wallet Balance": balance,
	})

	var discount models.Discount
	database.DB.Where("product_id = ?", ProductId).Find(&discount)
	discountAmount = discountAmount + product.Product_Price*discount.DiscountPercentage/100
	c.JSON(200, gin.H{
		"Product Name":  product.Product_Name,
		"Product Price": product.Product_Price,
		"Discount":      discountAmount,
		"selling Price": product.Product_Price - discountAmount,
	})
	TotalAmount = TotalAmount + product.Product_Price - discountAmount
	TotalMrp = TotalMrp + product.Product_Price
	c.JSON(200, gin.H{
		"mrp":   TotalMrp,
		"oid":   orderID,
		"total": TotalAmount,
	})

	if ApplyWallet == "apply" {

		if balance > TotalAmount {
			TotalAmount = 0
			newBalance := balance - TotalAmount
			database.DB.Model(&models.Wallet{}).Where("users_id = ?", UsersID).Update("balance", newBalance)
		} else if balance < TotalAmount {
			TotalAmount = TotalAmount - balance
			newBalance := 0
			database.DB.Model(&models.Wallet{}).Where("users_id = ?", UsersID).Update("balance", newBalance)
		}
	}

	if couponcode == "" {
		var checkInfo models.Checkoutinfo
		database.DB.Raw("UPDATE checkoutinfos SET order_id = ?,discount = ?,coupon_discount = ?,coupon_code = ?,total_mrp = ?,total = ? WHERE users_id = ?", orderID, discountAmount, 0, "Not Applied", TotalMrp, TotalAmount, UsersID).Scan(&checkInfo)
		c.JSON(200, gin.H{
			"messsage": "No COupon",
		})
		return
	}

	var coupon models.Coupon
	database.DB.Raw("SELECT coupon_percentage FROM coupons WHERE coupon_code = ?", couponcode).Find(&coupon)
	STotalAmpount := TotalMrp * coupon.CouponPercentage / 100
	TotalAmount = TotalAmount - STotalAmpount
	c.JSON(200, gin.H{
		"Total MRP":       TotalMrp,
		"Discount":        discountAmount,
		"Coupon Discount": STotalAmpount,
		"Total":           TotalAmount,
		"Wallet amount":   balance,
	})

	//database.DB.Save(&checkoutinfo)
	var checkInfo models.Checkoutinfo
	database.DB.Raw("UPDATE checkoutinfos SET order_id = ?,discount = ?,coupon_discount = ?,coupon_code = ?,total_mrp = ?,total = ? WHERE users_id = ?", orderID, discountAmount, STotalAmpount, couponcode, TotalMrp, TotalAmount, UsersID).Scan(&checkInfo)
	c.JSON(200, gin.H{
		"messsage": "Coupon Applied",
	})

}

func BuyNowCheckout(c *gin.Context) {

	userID := c.GetUint("id")
	PaymentMethod := c.Query("PaymentMethod")
	address := c.Query("AddressID")
	addressID, _ := strconv.Atoi(address)
	cod := "COD"
	RazorPay := "RAZORPAY"
	OrderStatus := "PENDING"
	var orderID string

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
		"tot":        total,
		"orderID":    orderID,
	})

	if address == "" {
		c.JSON(400, gin.H{
			"message": "Please Select Address",
		})
		return
	}

	var shipAddress struct {
		Name         string
		Phone_Number int
		Pincode      int
		House_Adress string
		Area         string
		Landmark     string
		City         string
	}
	database.DB.Raw("SELECT name,phone_number,pincode,house_adress,area,landmark,city FROM addresses WHERE id = ?", address).Scan(&shipAddress)

	c.JSON(200, gin.H{
		"Shipping Adress": shipAddress,
	})

	if PaymentMethod == "" {
		c.JSON(400, gin.H{
			"message": "Please Select Payment Method",
		})
		return
	}

	if PaymentMethod == cod {
		orders := models.Orders{UsersID: userID, AddressID: uint(addressID), OrderID: orderID, Discount: discount, CouponDiscount: couponDiscount, CouponCode: couponCode, Payment_Method: cod, Total_Amount: total, Status: OrderStatus}
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
	if PaymentMethod == RazorPay {
		orders := models.Orders{UsersID: userID, AddressID: uint(addressID), OrderID: orderID, Discount: discount, CouponDiscount: couponDiscount, CouponCode: couponCode, Payment_Method: RazorPay, Total_Amount: total, Status: OrderStatus}
		database.DB.Create(&orders)
	}

}
