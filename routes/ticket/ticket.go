package ticket

import (
	"github.com/Manuel-Leleuly/kanban-flow-go/middlewares"
	"github.com/Manuel-Leleuly/kanban-flow-go/models"
	ticketv1 "github.com/Manuel-Leleuly/kanban-flow-go/routes/ticket/v1"
	"github.com/gin-gonic/gin"
)

func TicketRoutes(router *gin.Engine, d *models.DBInstance) {
	tickets := router.Group("/tickets", d.MakeHTTPHandleFunc(middlewares.CheckAccessToken))

	ticketv1.TicketV1Routes(tickets, d)
}
