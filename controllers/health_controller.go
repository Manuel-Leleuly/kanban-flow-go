package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckServerHealth(c *gin.Context) {
	req := c.Request
	log.Printf("[HEALTHZ] method=%s proto=%s remote=%s host=%s path=%s\n", req.Method, req.Proto, req.RemoteAddr, req.Host, req.URL.Path)
	for k, v := range req.Header {
		log.Printf("[HEALTHZ-HEADER] %s: %v\n", k, v)
	}

	c.Status(http.StatusOK)
	if req.Method == http.MethodGet {
		c.Writer.Write([]byte("ok"))
	}
}
