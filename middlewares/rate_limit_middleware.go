package middlewares

import (
	"net/http"

	"github.com/Manuel-Leleuly/kanban-flow-go/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(1, 5)

func RateLimitMiddleware(c *gin.Context) {
	if !limiter.Allow() {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, models.ErrorMessage{
			Message: "rate limit exceeded",
		})
		return
	}

	c.Next()
}
