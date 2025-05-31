package main

import (
	"Diplomski/database"
	"Diplomski/routes"
	"Diplomski/utils"
	"fmt"


	"github.com/gin-gonic/gin"
)


func main() {
	err := database.ConnectDB() //database connection
    if err != nil {
			fmt.Println("Could not connect to the database.")
			return
		}
	server := gin.Default() 

	server.Use(utils.SetupCORS()) // CORS middleware

	routes.RegisterRoutes(server) 

	go utils.StartScheduler()

	server.Run(":8080")
	

}

