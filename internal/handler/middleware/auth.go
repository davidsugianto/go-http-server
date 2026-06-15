package middleware

import (
	"fmt"

	"github.com/davidsugianto/go-pkgs/response"
	"github.com/gin-gonic/gin"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "Bearer secret-token" {
			response.GinUnauthorized(c, fmt.Errorf("authorization header required"))
			c.Abort()
			return
		}
		c.Next()
	}
}
