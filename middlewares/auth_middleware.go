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

	// get token from cookie if it's not in header
	if bearerToken == "" {
		cookies := c.Request.Cookies()
		cookieMap := make(map[string]string)
		for _, cookie := range cookies {
			cookieMap[cookie.Name] = cookie.Value
		}

		accessToken, ok := cookieMap["access_token"]
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorMessage{
				Message: "token is expired",
			})
			return
		}
		bearerToken = "Bearer " + accessToken
	}

	accessToken, err := jwthelper.GetTokenStringFromHeader(bearerToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorMessage{
			Message: "token is expired",
		})
		return
	}

	user, err := jwthelper.ValidateToken(d, accessToken, false)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorMessage{
			Message: "token is expired",
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

	// get token from cookie if it's not in header
	if bearerToken == "" {
		cookies := c.Request.Cookies()
		cookieMap := make(map[string]string)
		for _, cookie := range cookies {
			cookieMap[cookie.Name] = cookie.Value
		}

		token, ok := cookieMap["refresh_token"]
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorMessage{
				Message: "token is expired",
			})
			return
		}
		bearerToken = "Bearer " + token
	}

	refreshToken, err := jwthelper.GetTokenStringFromHeader(bearerToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorMessage{
			Message: "token is expired",
		})
		return
	}

	user, err := jwthelper.ValidateToken(d, refreshToken, true)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorMessage{
			Message: "token is expired",
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
