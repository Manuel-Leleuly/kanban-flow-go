package helpers

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetBaseUrl(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	return scheme + "://" + c.Request.Host
}

func GenerateUUIDWithoutHyphen() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}
