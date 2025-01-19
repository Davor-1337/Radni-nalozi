package routes

import (
	"Diplomski/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {

	authenticatedAdmin := server.Group("/", middlewares.AdminMiddleware) 
	authenticatedTehnician := server.Group("/", middlewares.AdminOrAssignedTehnicianMiddleware)
	authenticatedRoles := server.Group("/", middlewares.RoleBasedAccess)
	{


		//klijenti
	authenticatedAdmin.POST("/klijenti", createClient)
	authenticatedAdmin.GET("/klijenti", getClients)
	authenticatedAdmin.GET("/klijenti/:id", getClient)
	authenticatedAdmin.PUT("/klijenti/:id", updateClient)
	authenticatedAdmin.DELETE("/klijenti/:id", deleteClient)
	//radni-nalozi
	authenticatedAdmin.GET("/radni-nalozi", getWorkOrders)
	authenticatedRoles.GET("/radni-nalozi/:id", getWorkOrder)
	authenticatedAdmin.POST("/radni-nalozi", createWorkOrder)
	authenticatedTehnician.PUT("/radni-nalozi/:id", updateWorkOrder)
	authenticatedAdmin.DELETE("/radni-nalozi/:id", deleteWorkOrder)
	authenticatedTehnician.PUT("/radni-nalozi/:id/zavrsi", finishWorkOrder)
	authenticatedTehnician.POST("/radni-nalozi/:id/materijal", inputMaterial)
	authenticatedTehnician.POST("/radni-nalozi/:id/sati", addWorkingHours)
	authenticatedTehnician.GET("/radni-nalozi/:id/sati", getWorkingHours)
	authenticatedAdmin.POST("/radni-nalozi/dodjela", assignWorkOrderToTehnician)
	authenticatedAdmin.PUT("/radni-nalozi/status/:id", updateWorkOrderStatus)
	authenticatedTehnician.PUT("/radni-nalozi/status", updateWorkOrderStatus)
	authenticatedAdmin.GET("/radni-nalozi/status", getWorkOrdersByStatus)
	authenticatedAdmin.GET("/radni-nalozi/filter", filterWorkOrders)
	authenticatedAdmin.GET("/radni-nalozi/arhiva", getArchivedWorkOrders)
	authenticatedAdmin.POST("/radni-nalozi/arhiva", archiveWorkOrder)
	authenticatedAdmin.GET("/radni-nalozi/arhiva/filter", filterArchivedWorkOrders)

	//serviseri
	authenticatedAdmin.GET("/serviseri", getTehnicians)
	authenticatedAdmin.POST("/serviseri", createTehnician)
	authenticatedAdmin.GET("/serviseri/:id", getTehnician)
	authenticatedAdmin.PUT("/serviseri/:id", updateTehnician)
	authenticatedAdmin.GET("/serviseri/:id/radni-nalozi", getWorkOrderByTehnician)
	authenticatedTehnician.GET("/serviseri/:id/sati", getTotalHours)
	//fakture
	authenticatedAdmin.GET("/fakture", getInvoices)
	authenticatedAdmin.POST("/fakture", createInvoice)
	authenticatedRoles.GET("/fakture/:id", getInvoice)
	authenticatedAdmin.PUT("/fakture/:id", updateInvoice)
	authenticatedAdmin.DELETE("/fakture/:id", deleteInvoice)
	authenticatedAdmin.GET("/fakture/generisi/:id", generateInvoice)
	authenticatedAdmin.POST("/fakture/generisi/:id", finalizeInvoice)

	//materijal
	authenticatedAdmin.POST("/materijali", addMaterial)
	authenticatedTehnician.GET("/materijali", getMaterials)
	authenticatedAdmin.GET("/materijali/:id", getMaterial)
	authenticatedTehnician.PUT("/materijali/:id", updateMaterial)
	authenticatedAdmin.DELETE("/materijali/:id", deleteMaterial)
	//users
	server.POST("/signup", signUp)
	server.POST("/login", login)
	}


}