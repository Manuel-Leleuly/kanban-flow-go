package kanbanv1

import (
	"github.com/Manuel-Leleuly/kanban-flow-go/controllers"
	"github.com/Manuel-Leleuly/kanban-flow-go/models"
	"github.com/gin-gonic/gin"
)

func KanbanV1Routes(router *gin.RouterGroup, d *models.DBInstance) {
	v1 := router.Group("/v1")
	{
		v1.POST("/tickets", d.MakeHTTPHandleFunc(controllers.CreateTicket))
		v1.GET("/tickets", d.MakeHTTPHandleFunc(controllers.GetTicketList))
		v1.GET("/tickets/:ticketId", d.MakeHTTPHandleFunc(controllers.GetTicketById))
		v1.PUT("/tickets/:ticketId", d.MakeHTTPHandleFunc(controllers.UpdateTicket))
		v1.DELETE("/tickets/:ticketId", d.MakeHTTPHandleFunc(controllers.DeleteTicket))
	}
}
