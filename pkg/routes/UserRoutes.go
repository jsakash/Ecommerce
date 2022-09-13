package routes

import (
	"github.com/gin-gonic/gin"
	//"github.com/go-chi/chi/middleware"
	"github.com/jsakash/ecommers/middleware"
	"github.com/jsakash/ecommers/pkg/controllers"
)

func UserRoutes(UserRoutes *gin.Engine) {

	UserRoutes.POST("/user/signup", controllers.Signup)
	UserRoutes.POST("/user/login", controllers.Login)
	UserRoutes.POST("/user/login/otp", controllers.OtpLog)
	UserRoutes.POST("/user/addaddress", middleware.UserAuth, controllers.AddAddress)
	UserRoutes.GET("/user/login/otp/validate", controllers.ValidateOtp)
	UserRoutes.GET("/user/selectaddress", middleware.UserAuth, controllers.SelectAddress)
	UserRoutes.POST("/user/addtocart", middleware.UserAuth, controllers.AddToCart)
	UserRoutes.GET("/user/productlist", middleware.UserAuth, controllers.ListAllProducts)
	UserRoutes.GET("/user/cartlist", middleware.UserAuth, controllers.CartList)
	UserRoutes.DELETE("/user/removefromcart/:id", middleware.UserAuth, controllers.RemoveFromCart)
	UserRoutes.POST("/user/addwishlist", middleware.UserAuth, controllers.AddToWishlist)
	UserRoutes.GET("/user/wishlist", middleware.UserAuth, controllers.Wishlist)
	UserRoutes.GET("/user/validate", middleware.UserAuth, controllers.Validate)
	//UserRoutes.DELETE("/delete", middleware.UserAuth, controllers.DeleteOtp)
	//UserRoutes.GET("/user/total", middleware.UserAuth, controllers.TotalPrice)
}
