package routes

import (
	"net/http"

	"github.com/Manuel-Leleuly/kanban-flow-go/controllers"
	"github.com/Manuel-Leleuly/kanban-flow-go/helpers"
	"github.com/Manuel-Leleuly/kanban-flow-go/middlewares"
	"github.com/Manuel-Leleuly/kanban-flow-go/models"
	"github.com/Manuel-Leleuly/kanban-flow-go/routes/iam"
	"github.com/Manuel-Leleuly/kanban-flow-go/routes/kanban"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func GetRoutes(d *models.DBInstance) *gin.Engine {
	router := gin.New()

	// use custom logger but keep the default recovery
	router.Use(middlewares.LoggerMiddleware, gin.Recovery())

	// implement swagger
	router.GET("/apidocs/*any", func(c *gin.Context) {
		if c.Request.RequestURI == "/apidocs/" {
			c.Redirect(http.StatusFound, "/apidocs/index.html")
		}
		ginSwagger.WrapHandler(swaggerfiles.Handler, ginSwagger.URL(helpers.GetBaseUrl(c)+"/apidocs/doc.json"))(c)
	})

	// implement rate limit and security
	router.Use(middlewares.RateLimitMiddleware, middlewares.SecurityHeadersMiddleware)

	// health check
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// implement websocket
	router.GET("/ws", controllers.WebSocketHandler)

	iam.IAMRoutes(router, d)
	kanban.KanbanRoutes(router, d)

	return router
}
