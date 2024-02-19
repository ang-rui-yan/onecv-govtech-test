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

	if err := h.Service.RegisterStudentsToTeacher(requestBody.Teacher, requestBody.Students); err != nil {
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
	panic("not implemented")
}


func (h *teacherHandler) RetrieveForNotificationsHandler(c *gin.Context) {
	panic("not implemented")
}