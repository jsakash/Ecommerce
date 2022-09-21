package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jsakash/ecommers/pkg/database"
	"github.com/jsakash/ecommers/pkg/models"
	razorpay "github.com/razorpay/razorpay-go"
)

func RazorPay(c *gin.Context) {
	UsersID := c.Param("id")
	var data models.Checkoutinfo
	database.DB.First(&data, "users_id = ?", UsersID)
	total := data.Total
	var user models.Users
	database.DB.First(&user, "id = ?", UsersID)

	client := razorpay.NewClient("rzp_test_XNE5NbTY0RM4SH", "6Pv9wOHuWCNZrQMtSvY57tst")

	amountInSubUnits := total * 100
	datas := map[string]interface{}{
		"amount":   amountInSubUnits,
		"currency": "INR",
		"receipt":  "receipt_id",
	}

	body, err := client.Order.Create(datas, nil)

	if err != nil {

		c.JSON(400, gin.H{
			"mesg": "errror",
		})
	}

	value := body["id"]

	orderIDCreated := value.(string)

	c.HTML(200, "app.html", gin.H{

		"UserID":           data.UsersID,
		"OrderIdCreated":   orderIDCreated,
		"TotalPrice":       total,
		"Name":             user.First_Name,
		"Email":            user.Email,
		"Phone_Number":     user.Phone_Number,
		"OrderId":          data.OrderID,
		"AmountInSubUnits": amountInSubUnits,
	})

	if err != nil {
		//http.Response.Status
		c.JSON(200, gin.H{
			"msg": orderIDCreated,
		})
		return
	}
}

func RPSuccess(c *gin.Context) {

	id := c.Query("user_id")
	razorPaymentId := c.Query("payment_id")
	razorPayOrderID := c.Query("order_id")
	signature := c.Query("signature")
	orderId := c.Query("id")

	check := models.Check{
		UserId:          id,
		RazorPaymentId:  razorPaymentId,
		RazorPayOrderID: razorPayOrderID,
		Signature:       signature,
		OrderId:         orderId,
	}
	database.DB.Create(&check)
	c.JSON(200, gin.H{
		"status": true,
	})

	var checkinfo models.Checkoutinfo
	database.DB.Raw("DELETE FROM checkoutinfos WHERE users_id = ?", id).Scan(&checkinfo)
	uid, _ := strconv.ParseUint(id, 10, 64)
	//fmt.Printf("%d\n", val_uint)
	placeOrder(uint(uid), orderId)
}

func SuccesPage(c *gin.Context) {

	c.HTML(200, "success.html", nil)
}
