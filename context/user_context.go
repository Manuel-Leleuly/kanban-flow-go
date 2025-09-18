package context

import (
	"errors"

	"github.com/Manuel-Leleuly/kanban-flow-go/models"
	"github.com/gin-gonic/gin"
)

func GetUserFromContext(c *gin.Context) (*models.User, error) {
	contextUser, exist := c.Get("user")
	if !exist {
		return &models.User{}, errors.New("user doesn't exist in context")
	}

	user, ok := contextUser.(*models.User)
	if !ok {
		return &models.User{}, errors.New("user doesn't exist in context")
	}

	return user, nil
}
