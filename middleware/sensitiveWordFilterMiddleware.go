package middleware

import "github.com/gin-gonic/gin"

func SensitiveWordFilterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
