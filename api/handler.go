package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type teacherHandler struct {
    Service TeacherService
}

func NewTeacherHandler(svc TeacherService) teacherHandler{
    return teacherHandler{Service: svc}
}

func (h *teacherHandler) RegisterHandler(c *gin.Context) {
	// request body
	var requestBody RegistrationRequestBody

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.RegisterStudentsToTeacher(requestBody.TeacherEmail, requestBody.StudentEmails); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Successfully registered!"})
}


func (h *teacherHandler) GetCommonStudentsHandler(c *gin.Context) {
	// querystring
	queryString := c.Request.URL.Query()

	// read the querystring for teachers
	teacherEmails := queryString["teacher"]

	studentEmails, err := h.Service.GetCommonStudents(teacherEmails)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"students": studentEmails})
}


func (h *teacherHandler) SuspendHandler(c *gin.Context) {
	// read request body
	var requestBody SuspendRequestBody

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.Suspend(requestBody.StudentEmail); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Successfully suspended!"})
}


func (h *teacherHandler) RetrieveForNotificationsHandler(c *gin.Context) {
	// read request body
	var requestBody NotificationRequestBody

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	studentEmails, err := h.Service.RetrieveForNotifications(requestBody.TeacherEmail, requestBody.Notification)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"recipients": studentEmails})
}