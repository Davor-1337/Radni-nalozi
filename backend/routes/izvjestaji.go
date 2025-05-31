package routes

import (
	"Diplomski/models"
	"fmt"
	"net/http"
	
	"strconv"

	"github.com/gin-gonic/gin"
)

func getReport(context *gin.Context) {
		workOrderId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		 context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse work order id."})
		 return
	}

	report, err := models.GetWorkOrderReport(workOrderId)

	if err != nil {
		 context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch work order report."})
		 return
	}

	 context.JSON(http.StatusOK, report)
	 
}

func getReports(context *gin.Context) {
	
	reports, err := models.GetAllOrderReportsShort()

	if err != nil {
		 context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch work order report."})
		 return
	}

	 context.JSON(http.StatusOK, reports)
	 
}

func getReportsByClient(c *gin.Context) {
	clientID, err := strconv.ParseInt(c.Param("clientId"), 10, 64)
	if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid client ID"})
			return
	}

	reports, err := models.GetOrderReportsShortForClient(clientID)
	if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch client reports"})
			return
	}

	c.JSON(http.StatusOK, reports)
}

func generateReportPDF(context *gin.Context) {
    workOrderId, err := strconv.ParseInt(context.Param("id"), 10, 64)
    fmt.Println("Received workOrderId:", workOrderId)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse work order id."})
        return
    }

    report, err := models.GetWorkOrderReport(workOrderId)
    if err != nil {
        context.JSON(http.StatusConflict, gin.H{"message": "Work order is not finished yet."})
        return
    }

    pdfBytes, err := models.GeneratePDF(report)
    if err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate PDF."})
        return
    }
    fmt.Println("PDF content length:", len(pdfBytes))

    
    context.Header("Content-Disposition", fmt.Sprintf("inline; filename=report_%d.pdf", workOrderId))
    context.Data(http.StatusOK, "application/pdf", pdfBytes)
}
