package api

import (
	"studentadmin/api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(service TeacherService) *gin.Engine {
	router := gin.New()

	router.Use(middleware.LogMiddleware())
	router.Use(middleware.ContentTypeMiddleware())

	apiGroup := router.Group("/api")
	{
		handler := NewTeacherHandler(service)
		apiGroup.POST("/register", handler.RegisterHandler)
		apiGroup.GET("/commonstudents", handler.GetCommonStudentsHandler)
		apiGroup.POST("/suspend", handler.SuspendHandler)
		apiGroup.POST("/retrievefornotifications", handler.RetrieveForNotificationsHandler)
	}

	return router
}
