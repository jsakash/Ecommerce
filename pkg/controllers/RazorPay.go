package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jsakash/ecommers/pkg/database"
	"github.com/jsakash/ecommers/pkg/models"
	razorpay "github.com/razorpay/razorpay-go"
)

func RazorPay(c *gin.Context) {

	var data models.Checkoutinfo
	database.DB.First(&data)
	total := data.Total
	var user models.Users
	database.DB.First(&user, "id = ?", data.UsersID)

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
		c.JSON(400, gin.H{
			"msg": orderIDCreated,
		})
		return
	}
}

func Success(c *gin.Context) {
	// id := c.Query("user_id")
	// razorPaymentId := c.Query("payment_id")
	// razorPayOrderID := c.Query("order_id")
	// signature := c.Query("signature")
	// orderId := c.Query("id")

	c.HTML(200, "success.html", gin.H{
		// "id":      id,
		// "rpid":    razorPaymentId,
		// "rpoid":   razorPayOrderID,
		// "sign":    signature,
		// "orderid": orderId,
	})

}

//.Execute(c.Writer, pageVariables)
