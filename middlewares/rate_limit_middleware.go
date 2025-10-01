package middlewares

import (
	"net/http"

	"github.com/Manuel-Leleuly/kanban-flow-go/initializer"
	"github.com/Manuel-Leleuly/kanban-flow-go/models"
	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware(c *gin.Context) {
	if initializer.Limiter == nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorMessage{
			Message: "error when configuring limiter",
		})
	}

	if !initializer.Limiter.Allow() {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, models.ErrorMessage{
			Message: "rate limit exceeded",
		})
		return
	}

	c.Next()
}
