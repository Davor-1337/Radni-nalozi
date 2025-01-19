package routes

import (
	"Diplomski/models"
	"net/http"
	"strconv"

	// "strconv"
	"github.com/gin-gonic/gin"
)

func getTehnicians(context *gin.Context){
	tehnicians, err := models.GetAllTehnicians()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Coult not fetch tehnicians."})
		return
	}
	context.JSON(http.StatusOK, tehnicians)
}

func createTehnician(context *gin.Context){
 var newTehnician models.Serviser
 err := context.ShouldBindJSON(&newTehnician)

 if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse requested data."})
        return
    }


 err = newTehnician.InsertTehnician()
 if err != nil {
	context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create tehnician."})
        return
 }

 context.JSON(http.StatusCreated, gin.H{"message": "Tehnician created!", "tehnician": newTehnician})
}

func getTehnician(context *gin.Context) {
	tehnicianId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		 context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse tehnician id."})
		 return
	}

	tehnician, err := models.GetTehnicianByID(tehnicianId)

	if err != nil {
		 context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch tehnician"})
		 return
	}

	 context.JSON(http.StatusOK, tehnician)
	 
}

func updateTehnician(context *gin.Context) {
	tehnicianId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not parse tehnician id."})
		return
	}

	
	var updatedTehnician models.Serviser 
	err = context.ShouldBindJSON(&updatedTehnician)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not parse requested data."})
		return
	}

	updatedTehnician.Serviser_ID = tehnicianId
	err = updatedTehnician.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update tehnician."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Tehnician updated succsessfully!"})

}

func getWorkOrderByTehnician(context *gin.Context) {
	tehnicianId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not parse tehnician id."})
		return
	}

	workOrder, err := models.GetWorkOrderByTehnicianID(tehnicianId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch data."})
		return
	}

	context.JSON(http.StatusOK, workOrder)
}



// func deleteTehnician(context *gin.Context) { PROVJERITI
// 	tehnicianId, err := strconv.ParseInt(context.Param("id"), 10, 64)
// 	if err != nil {
// 		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not parse tehnician id."})
// 		return
// 	}

// }

func getTotalHours(context *gin.Context) {
	serviserId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid tehnician id."})
		return
	}

	totalHours, err := models.GetHoursForTehnician(serviserId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not parse requested data."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"serviserID": serviserId, "Total hours": totalHours})
}