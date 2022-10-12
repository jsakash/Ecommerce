package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jsakash/ecommers/initializers"
	"github.com/jsakash/ecommers/pkg/database"
	"github.com/jsakash/ecommers/pkg/routes"
)

func init() {
	initializers.LoadEnvVariables() // Initializing .env file
	database.ConnectToDb()          // Initializing connection with database
	gin.SetMode(gin.ReleaseMode)    // Setting Gin to Release Mod

}

func main() {

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	routes.AdminRoutes(router)
	router.Run(":8080")                  // Port Declaration to serve the routes
	router.LoadHTMLGlob("./templates/*") // To load the html files
}
