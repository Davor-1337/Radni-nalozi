package middlewares

import (
	"Diplomski/models"
	"Diplomski/utils"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)


func AdminMiddleware(context *gin.Context) {
	
		token := context.Request.Header.Get("Authorization")
	
		if token == "" {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized."})
			context.Abort()
			return
		}

        if len(token) > 7 && token[:7] == "Bearer " {
        token = token[7:]
    }

		claims, err := utils.VerifyToken(token)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized."})
			context.Abort()
			return
		}

		role, ok := claims["Role"].(string)
		if !ok || role != "admin" {
			context.JSON(http.StatusForbidden, gin.H{"message": "Access forbidden: only admin can access this route."})
			context.Abort()
			return
		}


		
		context.Next()
	}
func AdminOrAssignedTehnicianMiddleware(context *gin.Context) {
    token := context.Request.Header.Get("Authorization")
    if token == "" {
        context.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized."})
        context.Abort()
        return
    }

    if len(token) > 7 && token[:7] == "Bearer " {
        token = token[7:]
    }

    claims, err := utils.VerifyToken(token)
    if err != nil {
        context.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token."})
        context.Abort()
        return
    }

    role, ok := claims["Role"].(string)
    if !ok {
        context.JSON(http.StatusForbidden, gin.H{"message": "Access forbidden: invalid role."})
        context.Abort()
        return
    }

    if role == "admin" {
        context.Next()
        return
    }

    if role == "serviser" {
        userIDFloat, ok := claims["User_ID"].(float64)
        if !ok {
            context.JSON(http.StatusForbidden, gin.H{"message": "Access forbidden: invalid user ID."})
            context.Abort()
            return
        }
        userID := int64(userIDFloat)

        
        idParam := context.Param("id")
        if idParam != "" {
            workOrderID, err := strconv.ParseInt(idParam, 10, 64)
            if err != nil {
                context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid work order ID."})
                context.Abort()
                return
            }

            isAssigned, err := models.IsTehnicianAssignedToWorkOrder(userID, workOrderID)
            if err != nil {
                context.JSON(http.StatusInternalServerError, gin.H{"message": "Error checking work order assignment."})
                context.Abort()
                return
            }

            if !isAssigned {
                context.JSON(http.StatusForbidden, gin.H{"message": "Access forbidden: you are not assigned to this work order."})
                context.Abort()
                return
            }
        }

        context.Set("userID", userID)
        context.Next()
        return
    }

    context.JSON(http.StatusForbidden, gin.H{"message": "Access forbidden."})
    context.Abort()
}


func AuthenticateAdminOrClient(context *gin.Context) {
    token := context.Request.Header.Get("Authorization")
    if token == "" {
        context.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized."})
        context.Abort()
        return
    }

    if len(token) > 7 && token[:7] == "Bearer " {
        token = token[7:]
    }

    claims, err := utils.VerifyToken(token)
    if err != nil {
        context.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token."})
        context.Abort()
        return
    }
    
    role, ok := claims["Role"].(string)

    if role == "admin" {
        context.Next()
        return
    }

    if !ok || role != "klijent" {
        context.JSON(http.StatusForbidden, gin.H{"message": "Access forbidden: only clients can access this route."})
        context.Abort()
        return
    }

    
    userIDFloat, ok := claims["User_ID"].(float64) 
    if !ok {
        context.JSON(http.StatusForbidden, gin.H{"message": "Access forbidden: invalid user ID."})
        context.Abort()
        return
    }
    userID := int64(userIDFloat) 
    fmt.Println("Klijent ID:", userID)

    context.Set("userID", userID)
    context.Next()
}

func RoleBasedAccess(context *gin.Context) {

    token := context.Request.Header.Get("Authorization")
    if token == "" {
        context.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized."})
        context.Abort()
        return
    }


    if len(token) > 7 && token[:7] == "Bearer " {
        token = token[7:]
    }


    claims, err := utils.VerifyToken(token)
    if err != nil {
        context.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token."})
        context.Abort()
        return
    }


    role, ok := claims["Role"].(string)
    if !ok {
        context.JSON(http.StatusForbidden, gin.H{"message": "Access forbidden: invalid role."})
        context.Abort()
        return
    }


    userIDFloat, ok := claims["User_ID"].(float64)
    if !ok {
        context.JSON(http.StatusForbidden, gin.H{"message": "Access forbidden: invalid user ID."})
        context.Abort()
        return
    }
    userID := int64(userIDFloat)
    context.Set("UserID", userID)


    if strings.Contains(context.FullPath(), "/api/obavjestenja/") {
        context.Next()
        return
    }


    workOrderIDStr := context.Param("id")
    if workOrderIDStr == "" {
        context.Next()
        return
    }

    workOrderID, err := strconv.Atoi(workOrderIDStr)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid work order ID."})
        context.Abort()
        return
    }


    switch role {
    case "admin":
        context.Next()
        return
    case "serviser":
        belongs, err := models.IsTehnicianAssignedToWorkOrder(userID, int64(workOrderID))
        if err != nil || !belongs {
            context.JSON(http.StatusForbidden, gin.H{"message": "Access forbidden: work order not assigned to you."})
            context.Abort()
            return
        }
    case "klijent":
        belongs, err := models.IsWorkOrderOwnedByClient(userID, int64(workOrderID))
        if err != nil || !belongs {
            context.JSON(http.StatusForbidden, gin.H{"message": "Access forbidden: work order not associated with your account."})
            context.Abort()
            return
        }
    default:
        context.JSON(http.StatusForbidden, gin.H{"message": "Access forbidden: role not recognized."})
        context.Abort()
        return
    }

    context.Next()
}
 
