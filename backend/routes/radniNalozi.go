package routes

import (
	"Diplomski/database"
	"Diplomski/models"
	"Diplomski/utils"
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"log"
	"strings"
	"net/http"
	"strconv"
	"time"
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

func getPendingWorkOrders(context *gin.Context){
    workOrders, err := models.GetPendingOrders()
    if err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch work orders."})
        return
    }
    context.JSON(http.StatusOK, workOrders)
}

func getWorkOrderStatusHandler(context *gin.Context) {
	completed, inProgress, err := models.GetWorkOrderStatusCount()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch work order status counts"})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"completed":  completed,
		"inProgress": inProgress,
	})
}

func getTotalWorkOrders(context *gin.Context) {
	total, err := models.GetTotalWorkOrderCount()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Greška pri dohvatanju broja radnih naloga."})
		return
	}
	context.JSON(http.StatusOK, gin.H{"ukupno": total})
}

  

func get4WorkOrders(context *gin.Context) {

	workOrders, err := models.Get4WorkOrders()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch work orders."})
		return
	}
	context.JSON(http.StatusOK, workOrders)
}

func getActiveWorkOrders(context *gin.Context) {

	workOrders, err := models.GetActiveWorkOrders()
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

func getWorkOrderStats(context *gin.Context) {
    
    stats, err := models.GetWorkOrderStats()
    if err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{
            "message": "Could not fetch work order statistics.",
            "error":   err.Error(),
        })
        return
    }

    context.JSON(http.StatusOK, stats)
}

func createWorkOrder(context *gin.Context) {
    
    bodyBytes, _ := io.ReadAll(context.Request.Body)
    

    
    context.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

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

func handleWorkOrderStatus(c *gin.Context) {
    
    type statusRequest struct {
        Akcija string `json:"akcija"`
    }

    
    idParam := c.Param("id")
    nalogID, err := strconv.ParseInt(idParam, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid work order ID."})
        return
    }

    
    var req statusRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})  
        return
    }

    
    var noviStatus string
    switch req.Akcija {
    case "Odobri", "Otvoren":
        noviStatus = "Otvoren"
        
    case "Odbij", "Odbijen":
        noviStatus = "Odbijen"
        
    default:
        
        c.JSON(http.StatusBadRequest, gin.H{"message": "Unknown action."})
        return
    }

    
    workOrder := &models.RadniNalog{ Nalog_ID: nalogID, Status: noviStatus }
    if err := workOrder.UpdateStatus(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Status update failed"})
        return
    }

    
    var klijentID int64
    err = database.DB.QueryRow(
        "SELECT Klijent_ID FROM RadniNalog WHERE Nalog_ID = @NalogID",
        sql.Named("NalogID", nalogID),
    ).Scan(&klijentID)
    if err != nil {
        log.Printf("❌ Ne mogu dohvatiti Klijent_ID za Nalog_ID %d: %v", nalogID, err)
    } else {
        log.Printf("✅ Dohvaćen Klijent_ID: %d", klijentID)

        
        notif := &models.Notification{
            UserID:  klijentID,
            Type:    "promena_statusa",
            Message: fmt.Sprintf("Vaš radni nalog #%d je ddd", nalogID),
        }
        if err := notif.Create(); err != nil {
            log.Printf("❌ Greška pri kreiranju notifikacije: %v", err)
            
        } else {
            log.Printf("✅ Notifikacija uspješno kreirana za korisnika %d", klijentID)
        }
    }

    
    c.JSON(http.StatusOK, gin.H{
        "message":     "Status radnog naloga je ažuriran",
        "novi_status": noviStatus,
    })
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
        context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid work order id."})
        return
    }

    raw, exists := context.Get("userID")
    if !exists {
        context.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
        return
    }

    serviserId, ok := raw.(int64)
    if !ok {
        if f, okf := raw.(float64); okf {
            serviserId = int64(f)
        } else {
            context.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid user ID type"})
            return
        }
    }

    
    if err := models.Finish(nalogId, serviserId); err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to finish work order."})
        return
    }

    
    var clientID int64
    err = database.DB.QueryRow(
        "SELECT Klijent_ID FROM RadniNalog WHERE Nalog_ID = @Nalog_ID",
        sql.Named("Nalog_ID", nalogId),
    ).Scan(&clientID)
    if err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to find client for work order."})
        return
    }

    
    notif := &models.Notification{
        UserID:  clientID,
        Type:    "Promjena statusa",
        Message: fmt.Sprintf("Vaš radni nalog #%d je završen.", nalogId),
    }

    go func() {
        if err := notif.Create(); err != nil {
            log.Printf("Failed to create notification: %v", err)
        }
    }()

    context.JSON(http.StatusOK, gin.H{"message": "Work order finished successfully."})
}


func assignWorkOrderToTechnician(c *gin.Context) {
    
    
    var task models.EvidencijaSati
    if err := c.ShouldBindJSON(&task); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid payload"})
        return
    }

    
    if err := task.AssignTask(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not assign work order to technician."})
        return
    }
    

    
    techNotif := &models.Notification{
        UserID:  task.Serviser_ID,
        Type:    "Dodjela",
        Message: fmt.Sprintf("Dodijeljen vam je radni nalog #%d.", task.Nalog_ID),
    }
    go func() {
        if err := techNotif.Create(); err != nil {
            log.Printf("❌ Failed to create notification for technician: %v", err)
        }
    }()

    
    var klijentID int64
    
    err := database.DB.QueryRow(`
        SELECT r.Klijent_ID
          FROM RadniNalog r
         WHERE r.Nalog_ID = @NalogID`,
        sql.Named("NalogID", task.Nalog_ID),
    ).Scan(&klijentID)
    if err != nil {
        log.Printf("❌ Failed to fetch client for work order %d: %v", task.Nalog_ID, err)
    } else {
        clientNotif := &models.Notification{
            UserID:  klijentID,
            Type:    "PromenaStatusa",
            Message: fmt.Sprintf("Vaš radni nalog #%d je otvoren.", task.Nalog_ID),
        }
        go func() {
            if err := clientNotif.Create(); err != nil {
                log.Printf("❌ Failed to create notification for client: %v", err)
            }
        }()
    }

    
    c.JSON(http.StatusOK, gin.H{"message": "Work order successfully assigned to technician."})
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

func updateWorkOrderStatus(c *gin.Context) {
    nalogID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
    var payload struct {
        Status string `json:"Status" binding:"required"`
    }
    if err := c.ShouldBindJSON(&payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
        return
    }

    
    if _, err := database.DB.Exec(
        `UPDATE RadniNalog SET Status = @Status WHERE Nalog_ID = @NalogID`,
        sql.Named("Status", payload.Status),
        sql.Named("NalogID", nalogID),
    ); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating status"})
        return
    }

    
    if payload.Status == "Odbijen" || payload.Status == "Završen" {
        var klijentID int64
        var clientEmail, clientName string
        if err := database.DB.QueryRow(`
            SELECT k.Klijent_ID, k.Email, k.Naziv
              FROM Klijenti k
              JOIN RadniNalog r ON k.Klijent_ID = r.Klijent_ID
             WHERE r.Nalog_ID = @NalogID`,
            sql.Named("NalogID", nalogID),
        ).Scan(&klijentID, &clientEmail, &clientName); err == nil {
            notif := &models.Notification{
                UserID:  klijentID,
                Type:    "promena_statusa",
                Message: fmt.Sprintf("Vaš radni nalog #%d je %s.", nalogID, strings.ToLower(payload.Status)),
            }
            notif.Create()
            go utils.SendStatusChangeEmail(clientEmail, clientName,
                fmt.Sprintf("%d", nalogID), payload.Status)
        }
    }
    

    c.JSON(http.StatusOK, gin.H{"message": "Status updated"})
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

func filterWorkOrders(c *gin.Context) {
    var filters struct {
        Search string `json:"search"`
    }
    if err := c.ShouldBindJSON(&filters); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": err.Error()})
        return
    }

    type WorkOrder struct {
        Nalog_ID       int       `json:"Nalog_ID"`
        Klijent_ID     int       `json:"Klijent_ID"`
        OpisProblema   string    `json:"OpisProblema"`
        Prioritet      string    `json:"Prioritet"`
        DatumOtvaranja time.Time `json:"DatumOtvaranja"`
        Status         string    `json:"Status"`
        Lokacija       string    `json:"Lokacija"`
        BrojNaloga     string    `json:"BrojNaloga"`
    }

    
    query := `
      SELECT Nalog_ID, Klijent_ID, OpisProblema, Prioritet,
             DatumOtvaranja, Status, Lokacija, BrojNaloga
        FROM RadniNalog
       WHERE 1=1`
    args := []interface{}{}

    if filters.Search != "" {
        
        query += `
        AND (
             LOWER(Status)       LIKE LOWER(@Search) OR
             LOWER(Prioritet)    LIKE LOWER(@Search) OR
             LOWER(OpisProblema) LIKE LOWER(@Search) OR
             LOWER(Lokacija)     LIKE LOWER(@Search) OR
             LOWER(BrojNaloga)   LIKE LOWER(@Search) OR
             CONVERT(VARCHAR, DatumOtvaranja, 23) LIKE @Search
        )`
        args = append(args, sql.Named("Search", "%"+filters.Search+"%"))
    }

    rows, err := database.DB.Query(query, args...)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Error querying", "error": err.Error()})
        return
    }
    defer rows.Close()

    var result []WorkOrder
    for rows.Next() {
        var wo WorkOrder
        if err := rows.Scan(
            &wo.Nalog_ID, &wo.Klijent_ID, &wo.OpisProblema,
            &wo.Prioritet, &wo.DatumOtvaranja, &wo.Status,
            &wo.Lokacija, &wo.BrojNaloga,
        ); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"message": "Error scanning", "error": err.Error()})
            return
        }
        result = append(result, wo)
    }

    c.JSON(http.StatusOK, gin.H{"work_orders": result})
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