package api

import (
	"errors"
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

	err := h.Service.RegisterStudentsToTeacher(requestBody.TeacherEmail, requestBody.StudentEmails)
	if err != nil {
		var statusCode int
		var errorMessage string

		switch err {
		case ErrTeacherNotFound:
			statusCode = http.StatusBadRequest
			errorMessage = ErrTeacherNotFound.Error()
		case ErrStudentNotFound:
			statusCode = http.StatusNotFound
			errorMessage = ErrStudentNotFound.Error()
		default:
			statusCode = http.StatusInternalServerError
			errorMessage = "Internal server error"
		}
		
		c.JSON(statusCode, gin.H{"error": errorMessage})
		return
	}

	c.Status(http.StatusNoContent)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrBadRequest.Error()})
		return
	}

	err := h.Service.Suspend(requestBody.StudentEmail)
	if err != nil {
		if errors.Is(err, ErrStudentAlreadySuspended) {
			c.Status(http.StatusNoContent)
		}

		var statusCode int
		var errorMessage string

		switch err {
		case ErrInvalidInput:
			statusCode = http.StatusBadRequest
			errorMessage = "Invalid input provided"
		case ErrStudentNotFound:
			statusCode = http.StatusNotFound
			errorMessage = "Student is not found"
		default:
			statusCode = http.StatusInternalServerError
			errorMessage = "Internal server error"
		}

		c.JSON(statusCode, gin.H{"error": errorMessage})
		return
	}

	c.Status(http.StatusNoContent)
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