package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS: Allow all cross-origin requests and permit all OPTIONS methods
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		// origin := c.Request.Header.Get("Origin")
		c.Header("Access-Control-Allow-Origin", "*")
		requestHeaders := c.Request.Header.Get("Access-Control-Request-Headers")
		c.Header("Access-Control-Allow-Headers", requestHeaders) //"Content-Type,Content-Length,AccessToken,X-CSRF-Token, Authorization, Token,X-Token,X-User-Id,apifoxtoken"
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS,DELETE,PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, New-Token, New-Expires-At")
		c.Header("Access-Control-Allow-Credentials", "true")

		// Allow all OPTIONS methods
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// Handle request
		c.Next()
	}
}
