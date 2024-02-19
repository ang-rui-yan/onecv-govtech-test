package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// take note of time now
		start := time.Now()

		// process request
		c.Next()

		// log details
		duration := time.Since(start)
		statusCode := c.Writer.Status()
		logMessage := fmt.Sprintf("[%d] %s %s %v", statusCode, c.Request.Method, c.Request.URL.Path, duration)
		fmt.Println(logMessage)
	}
}

