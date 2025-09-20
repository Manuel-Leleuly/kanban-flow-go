package middlewares

import "github.com/gin-gonic/gin"

func SecurityHeadersMiddleware(c *gin.Context) {
	c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
	c.Writer.Header().Set("X-Frame-Options", "DENY")
	c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
	c.Writer.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'")
	c.Next()
}
