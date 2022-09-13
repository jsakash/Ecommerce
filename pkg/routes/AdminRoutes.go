package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jsakash/ecommers/middleware"
	"github.com/jsakash/ecommers/pkg/controllers"
)

func AdminRoutes(AdminRoutes *gin.Engine) {

	AdminRoutes.POST("/admin/create", controllers.CreateAdmin)
	AdminRoutes.POST("/admin/login", controllers.AdminLogin)
	AdminRoutes.GET("/admin/validate", middleware.AdminAuth, controllers.AdminValidate)
	AdminRoutes.GET("/admin/userlist", middleware.AdminAuth, controllers.ListAllUsers)
	AdminRoutes.PATCH("/admin/blockuser/:id", middleware.AdminAuth, controllers.BlockUser)
	AdminRoutes.PATCH("/admin/unblockuser/:id", middleware.AdminAuth, controllers.UnBlockUser)
	AdminRoutes.DELETE("/admin/deleteuser/:id", middleware.AdminAuth, controllers.DeleteUser)
	AdminRoutes.POST("/admin/addproduct", middleware.AdminAuth, controllers.AddProduct)
	AdminRoutes.GET("/admin/productlist", middleware.AdminAuth, controllers.ListAllProducts)
	AdminRoutes.PUT("/admin/updateproduct/:id", middleware.AdminAuth, controllers.UpdateProduct)
	AdminRoutes.GET("/admin/getproduct/:id", middleware.AdminAuth, controllers.FetchProduct)
	AdminRoutes.DELETE("/admin/deleteproduct/:id", middleware.AdminAuth, controllers.DeleteProduct)
	AdminRoutes.POST("/admin/addcategory", middleware.AdminAuth, controllers.AddCategory)
	AdminRoutes.GET("/admin/categorylist", middleware.AdminAuth, controllers.CatogoryList)
	AdminRoutes.POST("/admin/addcolor", middleware.AdminAuth, controllers.AddColor)
	AdminRoutes.POST("/admin/addsize", middleware.AdminAuth, controllers.AddSize)

}
