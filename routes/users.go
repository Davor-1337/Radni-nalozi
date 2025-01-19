package routes

import (
	"Diplomski/models"
	"Diplomski/utils"
	"log"
	"net/http"

	// "strconv"
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

func login(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse requested data."})
		return
	}

	err = user.ValidateCredentials()
	log.Println("Fetched role:", user.Role)

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	token, err := utils.GenerateToken(user.User_ID, user.Username, user.Role)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not authenticate user."})
	}

	context.JSON(http.StatusOK, gin.H{"message": "Login successful.", "token": token})
}