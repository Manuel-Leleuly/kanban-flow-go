package context

import (
	"errors"

	"github.com/Manuel-Leleuly/kanban-flow-go/models"
	"github.com/gin-gonic/gin"
)

func GetUserFromContext(c *gin.Context) (*models.User, error) {
	contextUser, exist := c.Get("me")
	if !exist {
		return nil, errors.New("user doesn't exist in context")
	}

	user, ok := contextUser.(*models.User)
	if !ok {
		return nil, errors.New("user doesn't exist in context")
	}

	return user, nil
}

func RemoveUserFromContext(c *gin.Context) error {
	_, err := GetUserFromContext(c)
	if err != nil {
		return errors.New("user already removed")
	}

	delete(c.Keys, "me")
	return nil
}
