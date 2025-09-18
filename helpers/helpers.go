package helpers

import "github.com/gin-gonic/gin"

func GetBaseUrl(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	return scheme + "://" + c.Request.Host
}
