package routes

import (
	"github.com/gin-gonic/gin"
	//"github.com/go-chi/chi/middleware"
	"github.com/jsakash/ecommers/middleware"
	"github.com/jsakash/ecommers/pkg/controllers"
)

func UserRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.POST("/user/signup", controllers.Signup)
	incomingRoutes.POST("/user/login", controllers.Login)
	incomingRoutes.GET("/user/validate", middleware.RequireAuth, controllers.Validate)
	incomingRoutes.POST("/admin/create", controllers.CreateAdmin)
	incomingRoutes.POST("/admin/login", controllers.AdminLogin)
	incomingRoutes.GET("/admin/validate", middleware.RequireAuth, controllers.AdminValidate)
	incomingRoutes.GET("/admin/userlist", controllers.ListAllUsers)
	incomingRoutes.PATCH("/admin/blockuser/:id", controllers.BlockUser)
	incomingRoutes.PATCH("/admin/unblockuser/:id", controllers.UnBlockUser)
	incomingRoutes.DELETE("/admin/deleteuser/:id", controllers.DeleteUser)
	incomingRoutes.POST("/admin/addproduct", controllers.AddProduct)
	incomingRoutes.GET("/admin/productlist", controllers.ListAllProducts)
	incomingRoutes.PUT("/admin/updateproduct/:id", controllers.UpdateProduct)
	incomingRoutes.GET("/admin/getproduct/:id", controllers.FetchProduct)
	incomingRoutes.DELETE("/admin/deleteproduct/:id", controllers.DeleteProduct)

}
