package ticketv1

import (
	"github.com/Manuel-Leleuly/kanban-flow-go/controllers"
	"github.com/Manuel-Leleuly/kanban-flow-go/models"
	"github.com/gin-gonic/gin"
)

func TicketV1Routes(router *gin.RouterGroup, d *models.DBInstance) {
	v1 := router.Group("/v1")
	{
		v1.POST("/", d.MakeHTTPHandleFunc(controllers.CreateTicket))
		v1.GET("/", d.MakeHTTPHandleFunc(controllers.GetTicketList))
		v1.GET("/:ticketId", d.MakeHTTPHandleFunc(controllers.GetTicketById))
		v1.PUT("/:ticketId", d.MakeHTTPHandleFunc(controllers.UpdateTicket))
		v1.DELETE("/:ticketId", d.MakeHTTPHandleFunc(controllers.DeleteTicket))
	}
}
