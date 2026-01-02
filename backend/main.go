package main

import (
	"Diplomski/database"
	"Diplomski/routes"
	"Diplomski/utils"
	"fmt"


	"github.com/gin-gonic/gin"
)


func main() {
	err := database.ConnectDB()
if err != nil {
    fmt.Println("Database connection error:", err)
    return
}
	server := gin.Default() 

	server.Use(utils.SetupCORS()) // CORS middleware

	routes.RegisterRoutes(server) 

	go utils.StartScheduler()

	server.Run(":8080")
	

}

