package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jsakash/ecommers/initializers"
	"github.com/jsakash/ecommers/pkg/database"
	"github.com/jsakash/ecommers/pkg/routes"
)

func init() {
	initializers.LoadEnvVariables()
	database.ConnectToDb()

}

func main() {

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	routes.AdminRoutes(router)
	router.Run(":8080")

}
