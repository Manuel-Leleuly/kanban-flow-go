package middlewares

import (
	"net/http"

	jwthelper "github.com/Manuel-Leleuly/kanban-flow-go/helpers/jwt"
	"github.com/Manuel-Leleuly/kanban-flow-go/models"
	"github.com/gin-gonic/gin"
)

func CheckAccessToken(d *models.DBInstance, c *gin.Context) {
	bearerToken := c.GetHeader("Authorization")

	accessToken, err := jwthelper.GetTokenStringFromHeader(bearerToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorMessage{
			Message: "unauthorized access",
		})
		return
	}

	if err := jwthelper.ValidateAccessToken(d, accessToken); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorMessage{
			Message: "unauthorized access",
		})
		return
	}

	c.Next()
}
