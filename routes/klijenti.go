package routes

import (
	"Diplomski/models"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

func getClients(context *gin.Context) {
	clients, err := models.GetAllClients()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch clients."})
		return
	}
	context.JSON(http.StatusOK, clients)
}

func getClient(context *gin.Context) {
    clientId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse client id."})
		return
	}

    client, err := models.GetClientByID(clientId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch client."})
		return
	}

	context.JSON(http.StatusOK, client)
}

func createClient(context *gin.Context) {
    var newClient models.Klijent
    err := context.ShouldBindJSON(&newClient)

    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse requested data."})
        return
    }

    err = newClient.InsertClient()
    if err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create client."})
        return
    }

    context.JSON(http.StatusCreated, gin.H{"message": "Client created!", "client": newClient})
}

func updateClient(context *gin.Context) {
    	clientId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse client id."})
		return
	}


	// client, err := models.GetClientByID(clientId)

	
	// if event.UserID != userId {
	// 	context.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized to update event."})
	// 	return
	// }


	var updatedClient models.Klijent
	err = context.ShouldBindJSON(&updatedClient) 

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse requested data."})
		return
	}

	updatedClient.Klijent_ID = clientId 
	err = updatedClient.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update client."})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Client updated successfully!"})
}

func deleteClient(context *gin.Context){
clientId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse client id."})
		return
	}

	// userId := context.GetInt64("userId")
	client, err := models.GetClientByID(clientId)

	if err != nil{
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not fetch the client."})
		return
	}

	// if event.UserID != userId {
	// 	context.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized to update event."})
	// 	return
	// }

	err = client.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete the client ."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Delete completed successfully!"})
}