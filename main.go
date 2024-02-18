package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.New()

	api := router.Group("/api")
	
	api.POST("/register")
	api.GET("/commonstudents")
	api.POST("/suspend")
	api.POST("/retrievefornotifications")

	router.Run(":8080")
}