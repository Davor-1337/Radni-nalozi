package routes

import (
	"Diplomski/models"
	"Diplomski/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)


func signUp(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse requested data."})
		return
	}

	err = user.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save user."})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "User created successfully."})
}

func SignupRequest(context *gin.Context) {
	var zahtjev models.Zahtjev
	err := context.ShouldBindJSON(&zahtjev)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse requested data."})
		return
	}

	
	validRoles := map[string]bool{"admin": true, "serviser": true, "klijent": true}
	if !validRoles[zahtjev.Role] {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Non existing role."})
		return
	}

	
	hashedPassword, err := utils.HashPassword(zahtjev.Password)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not hash password."})
		return
	}
	zahtjev.Password = hashedPassword

	
	err = zahtjev.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save request."})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Request succesfully sent."})
}

func login(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse requested data."})
		return
	}

	err = user.ValidateCredentials()
	

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	token, err := utils.GenerateToken(user.User_ID, user.Username, user.Role)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not authenticate user."})
	}

	context.JSON(http.StatusOK, gin.H{"message": "Login successful.", "token": token, "username": user.Username})
}

func updatePassword(context *gin.Context)  {
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
	
		userIDFloat, ok := claims["User_ID"].(float64) 
    if !ok {
        context.JSON(http.StatusForbidden, gin.H{"message": "Access forbidden: invalid user ID."})
        context.Abort()
        return
    }
    userId := int64(userIDFloat) 
		user, err := models.GetUserByIdForPassword(userId)
		
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("Could not parse requested data: %v", err),
	})
		return 
	}

	type PasswordChangeRequest struct {
		NewPassword string `json:"new_password"`
		OldPassword string `json:"old_password"`
	}
	var req PasswordChangeRequest

	err = context.ShouldBindJSON(&req)
	
	if err != nil {
		fmt.Println("BindJSON error:", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse requested data.2"})
		return
	}
	passwordIsValid := utils.CheckPasswordHash(req.OldPassword, user.Password)

	if !passwordIsValid {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid old password."})
		return
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not hash password."})
		return
	}

	 err = user.ChangePassword(userId, hashedPassword)

	 if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not update password."})
		return 
	 }
		
	 context.JSON(http.StatusOK, gin.H{"message": "Password updated successfully." })

		
	}



