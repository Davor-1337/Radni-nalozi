package routes

import (
	"Diplomski/database"
	"Diplomski/models"
	"Diplomski/utils"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	
)

func getWorkOrders(context *gin.Context) {

	workOrders, err := models.GetAllWorkOrders()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch work orders."})
		return
	}
	context.JSON(http.StatusOK, workOrders)
}

func getWorkOrder(context *gin.Context) {
   workOrderId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse work order id."})
		return
	}

    workOrder, err := models.GetWorkOrderByID(workOrderId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch work order."})
		return
	}

	context.JSON(http.StatusOK, workOrder)
}

func createWorkOrder(context *gin.Context) {
  var newWorkOrder models.RadniNalog
  err := context.ShouldBindJSON(&newWorkOrder)

  if err != nil {
      context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse requested data."})
      return
  }

  err = newWorkOrder.InsertWorkOrder()
  if err != nil {
      context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create work order."})
      return
  }

  context.JSON(http.StatusCreated, gin.H{"message": "Work order created!", "Work order": newWorkOrder})
}

func updateWorkOrder(context *gin.Context) {

  workOrderId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse work order id."})
		return
	}

	var updatedWorkOrder models.RadniNalog
	err = context.ShouldBindJSON(&updatedWorkOrder) 

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse requested data."})
		return
	}

	updatedWorkOrder.Nalog_ID = workOrderId
	err = updatedWorkOrder.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update work order."})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Work order updated successfully!"})
}

func deleteWorkOrder(context *gin.Context){

	workOrderID, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse work order id."})
		return
	}

	workOrder, err := models.GetWorkOrderByID(workOrderID)

	if err != nil{
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not fetch the work order."})
		return
	}


	err = workOrder.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete the work order."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Delete completed successfully!"})
}

func finishWorkOrder(context *gin.Context) {

	nalogId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not parse work order id."})
		return
	}

	serviserIdStr := context.Query("Serviser_ID")
	if serviserIdStr == "" {
    context.JSON(http.StatusBadRequest, gin.H{"error": "Serviser_ID is required"})
    return
}

	serviserId, err := strconv.ParseInt(serviserIdStr, 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not parse tehnician id."})
		return
	}

	err = models.Finish(nalogId, serviserId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to finish work order."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Work order finished successfully."})
}

func assignWorkOrderToTehnician(context *gin.Context) {

	var task models.EvidencijaSati
	err := context.ShouldBindJSON(&task)

	err = task.AssignTask()
	 if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not assign work order to tehnician."})
		return
	 }

	 context.JSON(http.StatusOK, gin.H{"message": "Work order successfully assigned to tehnician."})
}

func inputMaterial(context *gin.Context) {
	nalogId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not parse work order id."})
		return
	}

	var materialAdded models.UtroseniMaterijal
	err = context.ShouldBindJSON(&materialAdded)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not parse requested data."})
		return
	}

	err = materialAdded.InputMaterial(nalogId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not add material"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Material successfully added."})

}

func addWorkingHours(context *gin.Context) {
	nalogId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not parse work order id."})
		return
	}

	var addedHours models.EvidencijaSati
	err = context.ShouldBindJSON(&addedHours)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not bind requested data."})
	}

	isAssigned, err := models.CheckTaskAssignment(nalogId, addedHours.Serviser_ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message":"Failed to validate assignment"})
		return
	}

	if !isAssigned {
		context.JSON(http.StatusForbidden, gin.H{"message": "Tehnician is not assigned to this work order."})
		return
	}

	err = addedHours.InsertHours(nalogId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not insert working hours."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Hours successfully added."})
}

func getWorkingHours(context *gin.Context){
	nalogId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not parse work order id."})
		return
	}

	workingHours, err := models.GetHours(nalogId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not parse requested data."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"Working hours": workingHours.BrojRadnihSati})
}

func updateWorkOrderStatus(context *gin.Context) {
    var input struct {
        WorkOrderID int    `json:"work_order_id" binding:"required"`
        Status      string `json:"status" binding:"required"`
    }

    if err := context.ShouldBindJSON(&input); err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input data"})
        return
    }

   
    query := `UPDATE RadniNalog SET Status = @Status WHERE Nalog_ID = @Nalog_ID`
    _, err := database.DB.Exec(query,
        sql.Named("Status", input.Status),
        sql.Named("Nalog_ID", input.WorkOrderID),
    )
    if err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating work order status"})
        return
    }

   
    var toEmail, clientName string
    query = `SELECT k.Email, k.Naziv 
             FROM Klijenti k 
             JOIN RadniNalog r ON k.Klijent_ID = r.Klijent_ID 
             WHERE r.Nalog_ID = @Nalog_ID`
    err = database.DB.QueryRow(query, sql.Named("Nalog_ID", input.WorkOrderID)).Scan(&toEmail, &clientName)
    if err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching client information"})
        return
    }

    
    go func() {
        err := utils.SendStatusChangeEmail(toEmail, clientName, fmt.Sprintf("%d", input.WorkOrderID), input.Status)
        if err != nil {
            fmt.Printf("Failed to send email notification: %v\n", err)
        }
    }()

    context.JSON(http.StatusOK, gin.H{
        "message":      "Work order status updated and notification sent",
        "workOrderID":  input.WorkOrderID,
        "newStatus":    input.Status,
    })
}

func getWorkOrdersByStatus(context *gin.Context) {
    
    var input struct {
    Status string `json:"status" binding:"required"`
}

if err := context.ShouldBindJSON(&input); err != nil {
    context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input data"})
    return
}

    query := `SELECT Nalog_ID, OpisProblema, Prioritet, DatumOtvaranja, Status FROM RadniNalog WHERE Status = @Status`
    rows, err := database.DB.Query(query, sql.Named("Status", input.Status))
    if err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching work orders"})
        return
    }
    defer rows.Close()


    var workOrders []struct {
        WorkOrderID int    `json:"work_order_id"`
        OpisProblema        string `json:"opisProblema"`
        Prioritet        string `json:"prioritet"`
        DatumOtvarnja string `json:"datumOtvaranja"`
        Status      string `json:"status"`
    }

    for rows.Next() {
        var workOrder struct {
            WorkOrderID int    `json:"work_order_id"`
        		OpisProblema        string `json:"opisProblema"`
						Prioritet        string `json:"prioritet"`
						DatumOtvarnja string `json:"datumOtvaranja"`
						Status      string `json:"status"`
        }
        if err := rows.Scan(&workOrder.WorkOrderID, &workOrder.OpisProblema, &workOrder.Prioritet, &workOrder.DatumOtvarnja, &workOrder.Status); err != nil {
            context.JSON(http.StatusInternalServerError, gin.H{"message": "Error scanning work orders"})
            return
        }
        workOrders = append(workOrders, workOrder)
    }

    
    context.JSON(http.StatusOK, gin.H{
        "status":     input.Status,
        "workOrders": workOrders,
    })
}

func filterWorkOrders(context *gin.Context) {
   
    var filters struct {
        DateFrom  string `json:"date_from"` 
        DateTo    string `json:"date_to"`      
        Priority  string `json:"priority"`   
        Status    string `json:"status"`
    }

    
    if err := context.ShouldBindJSON(&filters); err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input data", "error": err.Error()})
        return
    }

    
    query := `SELECT Nalog_ID, OpisProblema, Prioritet, DatumOtvaranja, Status FROM RadniNalog WHERE 1=1`



args := []interface{}{}
if filters.DateFrom != "" {
		query += " AND DatumOtvaranja >= @DateFrom"
    args = append(args, sql.Named("DateFrom", filters.DateFrom))	
}
if filters.DateTo != "" {
	query += " AND DatumOtvaranja <= @DateTo"
    args = append(args, sql.Named("DateTo", filters.DateTo))
}
if filters.Priority != "" {
	query += " AND Prioritet = @Priority"
    args = append(args, sql.Named("Priority", filters.Priority))
}
if filters.Status != "" {
	query += " AND Status = @Status"
    args = append(args, sql.Named("Status", filters.Status))
}


    rows, err := database.DB.Query(query, args...)
    if err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{"message": "Error querying work orders", "error": err.Error()})
        return
    }
    defer rows.Close()


    var workOrders []struct {
        ID       int    `json:"id"`
        Opis    string `json:"opis"`
        Prioritet  string `json:"prioritet"`
        Datum string `json:"datum"`
        Status   string `json:"status"`
    }

    for rows.Next() {
        var workOrder struct {
            ID       int    `json:"id"`
            Opis    string `json:"opis"`
            Prioritet  string `json:"prioritet"`
            Datum string `json:"datum"`
            Status   string `json:"status"`
        }
        if err := rows.Scan(&workOrder.ID, &workOrder.Opis, &workOrder.Prioritet, &workOrder.Datum, &workOrder.Status); err != nil {
            context.JSON(http.StatusInternalServerError, gin.H{"message": "Error scanning work orders", "error": err.Error()})
            return
        }
        workOrders = append(workOrders, workOrder)
    }

    
    context.JSON(http.StatusOK, gin.H{
        "work_orders": workOrders,
    })
}

func filterArchivedWorkOrders(context *gin.Context) {
   
    var filters struct {
        DateFrom  string `json:"date_from"` 
        DateTo    string `json:"date_to"`     
        Priority  string `json:"priority"`   
        Status    string `json:"status"`
    }

    
    if err := context.ShouldBindJSON(&filters); err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input data", "error": err.Error()})
        return
    }

    
    query := `SELECT Nalog_ID, OpisProblema, Prioritet, DatumOtvaranja, Status FROM ArhiviraniRadniNalozi WHERE 1=1`



args := []interface{}{}
if filters.DateFrom != "" {
		query += " AND DatumOtvaranja >= @DateFrom"
    args = append(args, sql.Named("DateFrom", filters.DateFrom))	
}
if filters.DateTo != "" {
	query += " AND DatumOtvaranja <= @DateTo"
    args = append(args, sql.Named("DateTo", filters.DateTo))
}
if filters.Priority != "" {
	query += " AND Prioritet = @Priority"
    args = append(args, sql.Named("Priority", filters.Priority))
}
if filters.Status != "" {
	query += " AND Status = @Status"
    args = append(args, sql.Named("Status", filters.Status))
}


    rows, err := database.DB.Query(query, args...)
    if err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{"message": "Error querying work orders", "error": err.Error()})
        return
    }
    defer rows.Close()


    var ArchivedWorkOrders []struct {
        ID       int    `json:"id"`
        Opis    string `json:"opis"`
        Prioritet  string `json:"prioritet"`
        Datum string `json:"datum"`
        Status   string `json:"status"`
    }

    for rows.Next() {
        var ArchivedWorkOrder struct {
            ID       int    `json:"id"`
            Opis    string `json:"opis"`
            Prioritet  string `json:"prioritet"`
            Datum string `json:"datum"`
            Status   string `json:"status"`
        }
        if err := rows.Scan(&ArchivedWorkOrder.ID, &ArchivedWorkOrder.Opis, &ArchivedWorkOrder.Prioritet, &ArchivedWorkOrder.Datum, &ArchivedWorkOrder.Status); err != nil {
            context.JSON(http.StatusInternalServerError, gin.H{"message": "Error scanning work orders", "error": err.Error()})
            return
        }
        ArchivedWorkOrders = append(ArchivedWorkOrders, ArchivedWorkOrder)
    }

    
    context.JSON(http.StatusOK, gin.H{
        "work_orders": ArchivedWorkOrders,
    })
}

func getArchivedWorkOrders(context *gin.Context) {

	workOrders, err := models.GetAllArchivedWorkOrders()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch archived work orders."})
		return
	}
	context.JSON(http.StatusOK, workOrders)
}

func archiveWorkOrder(context *gin.Context) {
    var nalog struct {
	ID int64 `json:"nalog_id"`
}

if err := context.ShouldBindJSON(&nalog); err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input data", "error": err.Error()})
        return
    }

err := models.Archive(nalog.ID)
if err != nil {
    context.JSON(http.StatusBadRequest, gin.H{"message": "Could not archive this work order."})
    return
}

context.JSON(http.StatusOK, gin.H{"message": "Work order archived successfully."})
}