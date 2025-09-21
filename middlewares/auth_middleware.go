package middlewares

import (
	"net/http"

	"github.com/Manuel-Leleuly/kanban-flow-go/context"
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

	user, err := jwthelper.ValidateToken(d, accessToken, false)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorMessage{
			Message: "unauthorized access",
		})
		return
	}

	// add user to context
	contextUser, err := context.GetUserFromContext(c)
	if err != nil || contextUser.ID != user.ID {
		c.Set("me", user)
	}

	c.Next()
}

func CheckRefreshToken(d *models.DBInstance, c *gin.Context) {
	bearerToken := c.GetHeader("Authorization")

	refreshToken, err := jwthelper.GetTokenStringFromHeader(bearerToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorMessage{
			Message: "unauthorized access",
		})
		return
	}

	user, err := jwthelper.ValidateToken(d, refreshToken, true)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorMessage{
			Message: "unauthorized access",
		})
		return
	}

	// add user to context
	contextUser, err := context.GetUserFromContext(c)
	if err != nil || contextUser.ID != user.ID {
		c.Set("me", user)
	}

	c.Next()
}
