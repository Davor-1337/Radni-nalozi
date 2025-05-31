package routes

import (
	"Diplomski/models"
    "Diplomski/utils"
	"fmt"
	"net/http"
	"strconv"
	"database/sql"
	"Diplomski/database"
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

func getOrdersForClient(context *gin.Context) {
	clientId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse client id."})
		return
	}

	orders, err := models.GetAllOrdersForClient(clientId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch orders for client."})
		return
	}
	fmt.Println(orders)
	context.JSON(http.StatusOK, orders)
}

func getOrdersForClientFromToken(context *gin.Context) {
   
    token := context.GetHeader("Authorization")
    if len(token) > 7 && token[:7] == "Bearer " {
        token = token[7:] 
    }

    
    claims, err := utils.VerifyToken(token)
    if err != nil {
        context.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token."})
        return
    }

    
    clientIDFloat, ok := claims["User_ID"].(float64)
    if !ok {
        context.JSON(http.StatusForbidden, gin.H{"message": "Invalid client ID."})
        return
    }
    clientID := int64(clientIDFloat)
    fmt.Println("Fetching orders for Client ID from token:", clientID)

    
    orders, err := models.GetAllOrdersForClient(clientID)
    if err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch orders for client."})
        return
    }

    
    context.JSON(http.StatusOK, orders)
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


	client, err := models.GetClientByID(clientId)

	if err != nil{
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not fetch the client."})
		return
	}

	err = client.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete the client ."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Delete completed successfully!"})
}

func filterClients(c *gin.Context) {
   
    var payload struct {
        Search string `json:"search"`
    }
    if err := c.ShouldBindJSON(&payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message":"Invalid input","error":err.Error()})
        return
    }

    
    type Klijent struct {
        Klijent_ID   int64  `json:"Klijent_ID"`
        Naziv        string `json:"Naziv"`
        KontaktOsoba string `json:"KontaktOsoba"`
        Email        string `json:"Email"`
        Tel          string `json:"Tel"`
        Adresa       string `json:"Adresa"`
        User_ID      int64  `json:"User_ID"`
    }

    
    query := `
      SELECT Klijent_ID, Naziv, KontaktOsoba, Email, Tel, Adresa, User_ID
        FROM Klijenti
       WHERE 1=1`
    args := []interface{}{}

    
    if payload.Search != "" {
        query += `
          AND (
            LOWER(Naziv)        LIKE LOWER(@Search) OR
            LOWER(KontaktOsoba) LIKE LOWER(@Search) OR
            LOWER(Email)        LIKE LOWER(@Search) OR
            LOWER(Tel)          LIKE LOWER(@Search) OR
            LOWER(Adresa)       LIKE LOWER(@Search)
          )`
        args = append(args, sql.Named("Search", "%"+payload.Search+"%"))
    }

    
    rows, err := database.DB.Query(query, args...)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message":"Query error","error":err.Error()})
        return
    }
    defer rows.Close()

    
    var list []Klijent
    for rows.Next() {
        var k Klijent
        if err := rows.Scan(
            &k.Klijent_ID,
            &k.Naziv,
            &k.KontaktOsoba,
            &k.Email,
            &k.Tel,
            &k.Adresa,
            &k.User_ID,
        ); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"message":"Scan error","error":err.Error()})
            return
        }
        list = append(list, k)
    }

    
    c.JSON(http.StatusOK, gin.H{"clients": list})
}