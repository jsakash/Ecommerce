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

	AdminRoutes.POST("/img", controllers.AddImage)

	// GET Routes
	AdminRoutes.GET("/admin/list/user/:page", middleware.AdminAuth, controllers.ListAllUsers)
	AdminRoutes.GET("/admin/list/product", middleware.AdminAuth, controllers.ListAllProducts)
	AdminRoutes.GET("/admin/list/category", middleware.AdminAuth, controllers.CatogoryList)
	AdminRoutes.GET("/admin/list/coupons", middleware.AdminAuth, controllers.ListCoupons)
	AdminRoutes.GET("/admin/list/discounts", middleware.AdminAuth, controllers.ListDiscount)
	AdminRoutes.GET("/admin/list/orderinfo", middleware.AdminAuth, controllers.OrderInfo)
	//AdminRoutes.GET("/admin/list/discounts", middleware.AdminAuth, controllers.ListDiscount)
	AdminRoutes.GET("/admin/getproduct/:id", middleware.AdminAuth, controllers.FetchProduct)
	AdminRoutes.GET("/admin/razorpar/details", middleware.AdminAuth, controllers.RazorPayInfo)
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
