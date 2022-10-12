package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jsakash/ecommers/middleware"
	"github.com/jsakash/ecommers/pkg/controllers"
)

func UserRoutes(UserRoutes *gin.Engine) {
	//POST Routes
	UserRoutes.POST("/user/signup", controllers.Signup)
	UserRoutes.POST("/user/login", controllers.Login)
	UserRoutes.POST("/user/login/otp", controllers.OtpLog)
	UserRoutes.POST("/user/addaddress", middleware.UserAuth, controllers.AddAddress)
	UserRoutes.POST("/user/addtocart", middleware.UserAuth, controllers.AddToCart)
	UserRoutes.POST("/user/addwishlist", middleware.UserAuth, controllers.AddToWishlist)
	UserRoutes.POST("/user/order/cancel", middleware.UserAuth, controllers.CancelOrder)
	UserRoutes.POST("/user/cartcheckout", middleware.UserAuth, controllers.CartCheckout)

	// GET Routes
	UserRoutes.GET("/user/login/otp/validate", controllers.ValidateOtp)
	UserRoutes.GET("/user/selectaddress", middleware.UserAuth, controllers.SelectAddress)
	UserRoutes.GET("/user/productlist", middleware.UserAuth, controllers.ListAllProducts)
	UserRoutes.GET("/user/cartlist", middleware.UserAuth, controllers.CartList)
	UserRoutes.GET("/user/wishlist", middleware.UserAuth, controllers.Wishlist)
	UserRoutes.GET("/user/validate", middleware.UserAuth, controllers.Validate)
	UserRoutes.GET("/user/cartcheckoutdetails", middleware.UserAuth, controllers.CartCheckoutDetails)
	UserRoutes.GET("/user/product/search", middleware.UserAuth, controllers.SearchProduct)
	UserRoutes.GET("/user/wallet/history", middleware.UserAuth, controllers.WalletInfo)

	UserRoutes.GET("/user/oders", middleware.UserAuth, controllers.OrderedItems)
	UserRoutes.GET("/user/wallet/balance", middleware.UserAuth, controllers.WalletBalance)
	UserRoutes.GET("/user/payment/razorpay/:id", controllers.RazorPay)
	UserRoutes.GET("/payment-success", controllers.RPSuccess)
	UserRoutes.GET("/success", controllers.SuccesPage)

	UserRoutes.GET("/failed", func(c *gin.Context) {
		c.HTML(http.StatusOK, "failed.html", gin.H{
			"title": "Main website",
		})
	})

	// PUT Router
	UserRoutes.PUT("user/profile/edit/address", middleware.UserAuth, controllers.EditAddress)
	//PATCH Routes
	UserRoutes.PATCH("/user/changepassword", middleware.UserAuth, controllers.ChangePasswors)
	//DELETE Routes
	UserRoutes.DELETE("/user/removefromcart/:id", middleware.UserAuth, controllers.RemoveFromCart)
	UserRoutes.LoadHTMLGlob("./templates/*")

}
