package kanban

import (
	"github.com/Manuel-Leleuly/kanban-flow-go/middlewares"
	"github.com/Manuel-Leleuly/kanban-flow-go/models"
	kanbanv1 "github.com/Manuel-Leleuly/kanban-flow-go/routes/kanban/v1"
	"github.com/gin-gonic/gin"
)

func KanbanRoutes(router *gin.Engine, d *models.DBInstance) {
	kanban := router.Group("/kanban", d.MakeHTTPHandleFunc(middlewares.CheckAccessToken))

	kanbanv1.KanbanV1Routes(kanban, d)
}
