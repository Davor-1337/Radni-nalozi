package routes

import (
	"Diplomski/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getNotifications(c *gin.Context) {
	userID, exists := c.Get("UserID")
	if !exists {
		log.Println("UserID nije postavljen u kontekstu.")
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Not authenticated"})
			return
	}

	list, err := models.GetNotificationsForUser(userID.(int64))
	if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch notifications"})
			return
	}
	c.JSON(http.StatusOK, list)
}

func DeleteNotification(c *gin.Context) {

	userID, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
	}


	userIDInt64, ok := userID.(int64)
	if !ok {
			
			userIDFloat, ok := userID.(float64)
			if !ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UserID format"})
					return
			}
			userIDInt64 = int64(userIDFloat)
	}


	notificationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
			return
	}

	
	belongs, err := models.DoesNotificationBelongToUser(userIDInt64, notificationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while checking notification"})
			return
	}
	if !belongs {
		c.JSON(http.StatusForbidden, gin.H{"error": "Notification does not belong to user"})
			return
	}


	notification := models.Notification{ID: notificationID}
	if err := notification.Delete(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while deleting notification"})
			return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification deleted successfully"})
}