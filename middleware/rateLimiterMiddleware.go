package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"note_app_server/utils"
)

// RateLimiterMiddleware 限流
func RateLimiterMiddleware(rate, capacity int) gin.HandlerFunc {
	return func(c *gin.Context) {
		pass, err := utils.RequestLimiter(c, rate, capacity)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  "请稍后尝试访问",
			})
		}
		if !pass {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"code": http.StatusServiceUnavailable,
				"msg":  "请勿频繁请求数据",
			})
		}
		c.Next()
	}
}
