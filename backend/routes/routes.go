package routes

import (
	"Diplomski/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {

	authenticatedAdmin := server.Group("/", middlewares.AdminMiddleware) 
	authenticatedTehnician := server.Group("/", middlewares.AdminOrAssignedTehnicianMiddleware)
	authenticatedRoles := server.Group("/", middlewares.RoleBasedAccess)
	authenticatedClient := server.Group("/", middlewares.AuthenticateAdminOrClient)
	{


		//klijenti
	authenticatedAdmin.POST("/api/klijenti", createClient)
	authenticatedAdmin.POST("/api/klijenti/filter", filterClients)
	authenticatedAdmin.GET("/api/klijenti", getClients)
	authenticatedAdmin.GET("/api/klijenti/:id", getClient)
	authenticatedAdmin.PUT("/api/klijenti/:id", updateClient)
	authenticatedAdmin.DELETE("/api/klijenti/:id", deleteClient)
	authenticatedAdmin.GET("/api/klijenti/:id/nalozi", getOrdersForClient)
	authenticatedClient.GET("/api/klijenti/radni-nalozi", getOrdersForClientFromToken)
	//radni-nalozi

	authenticatedTehnician.GET("/api/radni-nalozi", getWorkOrders)
	authenticatedTehnician.GET("/api/radni-nalozi/ukupno", getTotalWorkOrders)
	authenticatedTehnician.GET("/api/radni-nalozi/stats", getWorkOrderStats)
	authenticatedTehnician.GET("/api/radni-nalozi/status-count", getWorkOrderStatusHandler)
	authenticatedTehnician.GET("/api/radni-nalozi/4", get4WorkOrders)
	authenticatedAdmin.GET("/api/radni-nalozi/na-cekanju", getPendingWorkOrders)
	authenticatedTehnician.GET("/api/radni-nalozi/aktivni", getActiveWorkOrders)
	authenticatedRoles.GET("/api/radni-nalozi/:id", getWorkOrder)
	authenticatedClient.POST("/api/radni-nalozi", createWorkOrder)
	authenticatedAdmin.POST("/api/radni-nalozi/odaberi/:id", handleWorkOrderStatus)
	authenticatedAdmin.PUT("/api/radni-nalozi/:id", updateWorkOrder)
	authenticatedAdmin.DELETE("/api/radni-nalozi/:id", deleteWorkOrder)
	authenticatedTehnician.PUT("/api/radni-nalozi/:id/zavrsi", finishWorkOrder)
	authenticatedTehnician.POST("/api/radni-nalozi/:id/materijal", inputMaterial)
	authenticatedTehnician.POST("/api/radni-nalozi/:id/sati", addWorkingHours)
	authenticatedTehnician.GET("/api/radni-nalozi/:id/sati", getWorkingHours)
	authenticatedAdmin.POST("/api/radni-nalozi/dodjela", assignWorkOrderToTechnician)
	authenticatedAdmin.PUT("/api/radni-nalozi/status/:id", updateWorkOrderStatus)
	authenticatedTehnician.PUT("/api/radni-nalozi/status", updateWorkOrderStatus)
	authenticatedAdmin.GET("/api/radni-nalozi/status", getWorkOrdersByStatus)
	authenticatedAdmin.POST("/api/radni-nalozi/filter", filterWorkOrders)
	authenticatedAdmin.GET("/api/radni-nalozi/arhiva", getArchivedWorkOrders)
	authenticatedAdmin.POST("/api/radni-nalozi/arhiva", archiveWorkOrder)
	authenticatedAdmin.GET("/api/radni-nalozi/arhiva/filter", filterArchivedWorkOrders)

	//serviseri
	authenticatedAdmin.GET("/api/serviseri", getTehnicians)
	authenticatedAdmin.POST("/api/serviseri", createTehnician)
	authenticatedAdmin.GET("/api/serviseri/:id", getTehnician)
	authenticatedAdmin.PUT("/api/serviseri/:id", updateTehnician)
	authenticatedTehnician.GET("/api/serviseri/:id/radni-nalozi", getWorkOrderByTehnician)
	authenticatedTehnician.GET("/api/serviseri/radni-nalozi", getWorkOrderForTehnician)
	authenticatedAdmin.GET("/api/serviseri/:id/radni-nalozi/details", getDetailsForTehnician)
	authenticatedTehnician.GET("/api/serviseri/:id/sati", getTotalHours)
	//fakture
	authenticatedAdmin.GET("/api/fakture", getInvoices)
	authenticatedAdmin.POST("/api/fakture", createInvoice)
	authenticatedRoles.GET("/api/fakture/:id", getInvoice)
	authenticatedAdmin.PUT("/api/fakture/:id", updateInvoice)
	authenticatedAdmin.DELETE("/api/fakture/:id", deleteInvoice)
	authenticatedClient.GET("/api/fakture/klijent/:clientId", getInvoicesByClient)
	authenticatedAdmin.GET("/api/fakture/generisi/:id", generateInvoice)
	authenticatedAdmin.POST("/api/fakture/generisi/:id", finalizeInvoice)

	//materijal
	authenticatedAdmin.POST("/api/materijali", addMaterial)
	authenticatedTehnician.POST("/api/materijali/filter", filterMaterials)
	authenticatedTehnician.GET("/api/materijali", getMaterials)
	authenticatedAdmin.GET("/api/materijali/:id", getMaterial)
	authenticatedTehnician.GET("/api/materijali/kategorija", getMaterialCategoryStats)
	authenticatedTehnician.PUT("/api/materijali/:id", updateMaterial)
	authenticatedAdmin.DELETE("/api/materijali/:id", deleteMaterial)
	//izvjestaj
	authenticatedRoles.GET("/api/izvjestaji/radni-nalog/:id", getReport)
	authenticatedAdmin.GET("/api/izvjestaji/radni-nalog/skraceni", getReports)
	authenticatedRoles.GET("/api/izvjestaji/radni-nalog/pdf/:id", generateReportPDF)
	authenticatedClient.GET("/api/izvjestaji/klijent/:clientId", getReportsByClient)
	//obavjestenja 
	authenticatedRoles.GET("/api/obavjestenja", getNotifications)
	authenticatedRoles.DELETE("/api/obavjestenja/:id", DeleteNotification)
	//users
	server.POST("/api/zahtjevi", createZahtjev)
	server.POST("/api/zahtjevi/azuriraj", handleRequest)
	server.GET("/api/zahtjevi", GetAllZahtjevi)
	server.POST("/api/signup", signUp)
	server.POST("/api/login", login)
	server.PUT("/api/updatePassword", updatePassword)
	}


}