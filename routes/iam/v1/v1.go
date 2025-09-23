package iamv1

import (
	"github.com/Manuel-Leleuly/kanban-flow-go/controllers"
	"github.com/Manuel-Leleuly/kanban-flow-go/middlewares"
	"github.com/Manuel-Leleuly/kanban-flow-go/models"
	"github.com/gin-gonic/gin"
)

func IAMV1Routes(router *gin.RouterGroup, d *models.DBInstance) {
	v1 := router.Group("/v1")
	{
		v1.POST("/login", d.MakeHTTPHandleFunc(controllers.Login))
		v1.POST("/users", d.MakeHTTPHandleFunc(controllers.CreateUser))
	}

	withAccessToken := v1.Group("/", d.MakeHTTPHandleFunc(middlewares.CheckAccessToken))
	{
		withAccessToken.GET("/users/me", controllers.GetMe)
		withAccessToken.POST("/logout", controllers.Logout)
	}

	withRefreshToken := v1.Group("/", d.MakeHTTPHandleFunc(middlewares.CheckRefreshToken))
	{
		withRefreshToken.POST("/token/refresh", d.MakeHTTPHandleFunc(middlewares.CheckRefreshToken), controllers.RefreshToken)
	}
}
