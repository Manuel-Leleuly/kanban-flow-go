package iam

import (
	"github.com/Manuel-Leleuly/kanban-flow-go/models"
	iamv1 "github.com/Manuel-Leleuly/kanban-flow-go/routes/iam/v1"
	"github.com/gin-gonic/gin"
)

func IAMRoutes(router *gin.Engine, d *models.DBInstance) {
	iam := router.Group("/iam")

	iamv1.IAMV1Routes(iam, d)
}
