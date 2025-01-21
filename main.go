package main

import (
	"Diplomski/database"
	"Diplomski/routes"
	"fmt"
	"github.com/gin-gonic/gin"
	"Diplomski/utils"

)

func main() {
	err := database.ConnectDB()
    if err != nil {
			fmt.Println("Could not connect to the database.")
		}
	server := gin.Default() 

	routes.RegisterRoutes(server)

	server.Run(":8080")
	
	utils.StartScheduler()
}

