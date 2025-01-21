package routes

import (
	"Diplomski/database"
	"Diplomski/models"
	"database/sql"
	"net/http"
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
)

func getInvoices(context *gin.Context){
	invoices, err := models.GetAllInvoices()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch invoices."})
	}
	context.JSON(http.StatusOK, invoices)
}

func createInvoice(context *gin.Context){
	var newInvoice models.Faktura
	err := context.ShouldBindJSON(&newInvoice)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not parse requested data."})
		return
	}
	err = newInvoice.InsertInvoice()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create invoice."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Invoice created!", "invoice": newInvoice})
}

func getInvoice(context *gin.Context) {
	invoiceId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message":"Could not parse invoice id."})
		return
	}

	invoice, err := models.GetInvoiceByID(invoiceId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch invoice."})
	}

	context.JSON(http.StatusOK, invoice)
}

func updateInvoice(context *gin.Context) {
	invoiceId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not parse invoice id."})
		return
	}

	var updatedInvoice models.Faktura
	err = context.ShouldBindJSON(&updatedInvoice)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message":"Could not parse requested data."})
		return
	}

	updatedInvoice.Faktura_ID = invoiceId
	err = updatedInvoice.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update invoice."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Invoice updated successfully!"})
}

func deleteInvoice(context *gin.Context) {
	invoiceId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not parse invoice id."})
		return
	}

	invoice, err := models.GetInvoiceByID(invoiceId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch invoice."})
		return
	}

	err = invoice.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete invoice."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Invoice deleted successfully!"})
}

func generateInvoice(context *gin.Context) {
    workOrderId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not parse work order id."})
		return
	}

    // Dohvati sate rada
    var totalHours float64
		query := `SELECT SUM(BrojRadnihSati) 
				FROM EvidencijaSati 
        WHERE Nalog_ID = @Nalog_ID`
    row := database.DB.QueryRow(query, sql.Named("Nalog_ID", workOrderId))
  
    err = row.Scan(&totalHours)

    // Dohvati utrošeni materijal
    type Material struct {
        Name     string
        Quantity int
        Price    float64
    }

    var materials []Material
		query = `
        SELECT m.NazivMaterijala, um.KolicinaUtrosena, m.Cijena 
        FROM UtroseniMaterijal um
        JOIN Materijal m ON um.Materijal_ID = m.Materijal_ID
        WHERE um.Nalog_ID = @Nalog_ID`

    rows, err := database.DB.Query(query, sql.Named("Nalog_ID", workOrderId))
    var materialsUsed string
    err = rows.Scan(&materialsUsed)
    
    var totalMaterialCost float64
    for rows.Next() {
        var material Material
        if err := rows.Scan(&material.Name, &material.Quantity, &material.Price); err != nil {
            context.JSON(http.StatusInternalServerError, gin.H{"message": "Error reading materials"})
            return
        }
        materials = append(materials, material)
        totalMaterialCost += material.Price * float64(material.Quantity)
    }

    
    laborCost := totalHours * 20 
    totalCost := laborCost + totalMaterialCost

    //Generisanje automatske fakture bez unosa u bazu
    context.JSON(http.StatusOK, gin.H{
        "work_order_id": workOrderId,
        "total_hours":   totalHours,
        "labor_cost":    laborCost,
        "materials":     materials,
        "total_cost":    totalCost,
    })
		
}

func finalizeInvoice(context *gin.Context) {
     workOrderId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not parse work order id."})
		return
	}

  
    var exists bool
    query :=  `SELECT CASE 
    WHEN EXISTS (SELECT 1 FROM RadniNalog WHERE Nalog_ID = @Nalog_ID) 
    THEN 1 
    ELSE 0 
END AS ExistsCheck`
    err = database.DB.QueryRow(query, sql.Named("Nalog_ID", workOrderId)).Scan(&exists)

    if err != nil || !exists {
        context.JSON(http.StatusNotFound, gin.H{"message": "Work order not found"})
        return
    }

  
    var totalPrice float64
    query = `SELECT 
        (SELECT SUM(BrojRadnihSati * @Satnica) FROM EvidencijaSati WHERE Nalog_ID = @Nalog_ID) +
        (SELECT SUM(m.Cijena * um.KolicinaUtrosena) 
         FROM UtroseniMaterijal um 
         JOIN Materijal m ON m.Materijal_ID = um.Materijal_ID
         WHERE um.Nalog_ID = @Nalog_ID) AS TotalPrice`
    satnica := 20
    row := database.DB.QueryRow(query, 
    sql.Named("Satnica", satnica),
    sql.Named("Nalog_ID", workOrderId),)


    err = row.Scan(&totalPrice)
    if err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{"message": "Error calculating total price"})
        return
    }

   
    // Unos fakturu u bazu
    query = `INSERT INTO Faktura (Nalog_ID, Iznos, DatumFakture) VALUES (@Nalog_ID, @Iznos, GETDATE())`
    
    stmt, err := database.DB.Prepare(query)
    if err != nil {
        return 
    }

    defer stmt.Close()

    _, err = stmt.Exec(
        sql.Named("Nalog_ID", workOrderId),
        sql.Named("Iznos", totalPrice))

    if err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{"message": "Error finalizing invoice"})
        return
    }

    datum := time.Now().Format("2006-01-02 15:04:05")

    context.JSON(http.StatusOK, gin.H{
        "message":      "Invoice finalized successfully",
        "workOrderID":  workOrderId,
        "Iznos":     totalPrice,
        "DatumFakture":   datum,
    })
}
