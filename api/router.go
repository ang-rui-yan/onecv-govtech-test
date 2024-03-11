package api

import (
	"net/http"
	"studentadmin/api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(service TeacherService) *gin.Engine {
	router := gin.New()

	router.Use(middleware.LogMiddleware())
	router.Use(middleware.ContentTypeMiddleware())

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "service is running")
	})

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
