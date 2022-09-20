package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jsakash/ecommers/middleware"
	"github.com/jsakash/ecommers/pkg/controllers"
)

func AdminRoutes(AdminRoutes *gin.Engine) {

	// POST Routes
	AdminRoutes.POST("/admin/create", controllers.CreateAdmin)
	AdminRoutes.POST("/admin/login", controllers.AdminLogin)
	AdminRoutes.POST("/admin/add/product", middleware.AdminAuth, controllers.AddProduct)
	AdminRoutes.POST("/admin/add/category", middleware.AdminAuth, controllers.AddCategory)
	AdminRoutes.POST("/admin/add/color", middleware.AdminAuth, controllers.AddColor)
	AdminRoutes.POST("/admin/add/size", middleware.AdminAuth, controllers.AddSize)
	AdminRoutes.POST("/admin/add/coupon", middleware.AdminAuth, controllers.AddCoupon)
	AdminRoutes.POST("/admin/add/discount", middleware.AdminAuth, controllers.AddDiscount)
	// GET Routes
	AdminRoutes.GET("/admin/list/user", middleware.AdminAuth, controllers.ListAllUsers)
	AdminRoutes.GET("/admin/list/product", middleware.AdminAuth, controllers.ListAllProducts)
	AdminRoutes.GET("/admin/list/category", middleware.AdminAuth, controllers.CatogoryList)
	AdminRoutes.GET("/admin/list/coupons", middleware.AdminAuth, controllers.ListCoupons)
	AdminRoutes.GET("/admin/list/discounts", middleware.AdminAuth, controllers.ListDiscount)
	//AdminRoutes.GET("/admin/list/discounts", middleware.AdminAuth, controllers.ListDiscount)
	AdminRoutes.GET("/admin/getproduct/:id", middleware.AdminAuth, controllers.FetchProduct)
	// PATCH Routes
	AdminRoutes.PATCH("/admin/blockuser/:id", middleware.AdminAuth, controllers.BlockUser)
	AdminRoutes.PATCH("/admin/unblockuser/:id", middleware.AdminAuth, controllers.UnBlockUser)
	// PUT Routes
	AdminRoutes.PUT("/admin/updateproduct/:id", middleware.AdminAuth, controllers.UpdateProduct)
	//DELETE Routes
	AdminRoutes.DELETE("/admin/delete/coupon", middleware.AdminAuth, controllers.DeleteCoupon)
	AdminRoutes.DELETE("/admin/delete/discount", middleware.AdminAuth, controllers.DeleteDiscount)
	AdminRoutes.DELETE("/admin/delete/product/:id", middleware.AdminAuth, controllers.DeleteProduct)
	AdminRoutes.DELETE("/admin/delete/user/:id", middleware.AdminAuth, controllers.DeleteUser)

}

func UserRoutes(UserRoutes *gin.Engine) {
	//POST Routes
	UserRoutes.POST("/user/signup", controllers.Signup)
	UserRoutes.POST("/user/login", controllers.Login)
	UserRoutes.POST("/user/login/otp", controllers.OtpLog)
	UserRoutes.POST("/user/addaddress", middleware.UserAuth, controllers.AddAddress)
	UserRoutes.POST("/user/addtocart", middleware.UserAuth, controllers.AddToCart)
	UserRoutes.POST("/user/addwishlist", middleware.UserAuth, controllers.AddToWishlist)
	UserRoutes.POST("/user/order/cancel", middleware.UserAuth, controllers.CancelOrder)

	// GET Routes
	UserRoutes.GET("/user/login/otp/validate", controllers.ValidateOtp)
	UserRoutes.GET("/user/selectaddress", middleware.UserAuth, controllers.SelectAddress)
	UserRoutes.GET("/user/productlist", middleware.UserAuth, controllers.ListAllProducts)
	UserRoutes.GET("/user/cartlist", middleware.UserAuth, controllers.CartList)
	UserRoutes.GET("/user/wishlist", middleware.UserAuth, controllers.Wishlist)
	UserRoutes.GET("/user/validate", middleware.UserAuth, controllers.Validate)
	UserRoutes.GET("/user/cartcheckoutdetails", middleware.UserAuth, controllers.CartCheckoutDetails)
	UserRoutes.GET("/user/product/search", middleware.UserAuth, controllers.SearchProduct)
	UserRoutes.GET("/user/cartcheckout", middleware.UserAuth, controllers.CartCheckout)
	UserRoutes.GET("/user/oders", middleware.UserAuth, controllers.OrderedItems)
	UserRoutes.GET("/user/wallet/balance", controllers.WalletBalance)
	UserRoutes.GET("/user/payment/razorpay", controllers.RazorPay)
	UserRoutes.GET("/success", controllers.Success)
	// PUT Router
	UserRoutes.PUT("user/profile/edit", middleware.UserAuth, controllers.EditProfile)
	//PATCH Routes
	UserRoutes.PATCH("/user/changepassword", middleware.UserAuth, controllers.ChangePasswors)
	//DELETE Routes
	UserRoutes.DELETE("/user/removefromcart/:id", middleware.UserAuth, controllers.RemoveFromCart)
	UserRoutes.LoadHTMLGlob("./templates/*")
	//loading gohtml files from templates directory

	//UserRoutes.DELETE("/delete", middleware.UserAuth, controllers.DeleteOtp)
	//UserRoutes.GET("/user/total", middleware.UserAuth, controllers.TotalPrice)

}
