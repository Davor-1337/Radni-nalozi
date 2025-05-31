package utils

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func SetupCORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"}, 
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, 
		AllowHeaders:     []string{"Content-Type", "Authorization"}, 
		ExposeHeaders:    []string{"Authorization"}, 
		AllowCredentials: true,
		MaxAge:           3 * time.Hour,
	})
}