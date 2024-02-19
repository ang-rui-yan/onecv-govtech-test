package middleware

import (
	"github.com/gin-gonic/gin"
)

func ContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// if we are posting or anything similar with a request body, we will check for content type
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			if c.GetHeader("Content-Type") != "application/json" {
				c.AbortWithStatusJSON(415, gin.H{
					"error": "Content-Type must be application/json",
				})
				return
			}
		}

		c.Next()
	}
}