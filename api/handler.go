package api

import (
	"github.com/gin-gonic/gin"
)

type teacherHandler struct {
    Service TeacherService
}

func NewTeacherHandler(svc TeacherService) teacherHandler{
    return teacherHandler{Service: svc}
}

func (h *teacherHandler) RegisterHandler(c *gin.Context) {
	panic("not implemented")
}


func (h *teacherHandler) GetCommonStudentsHandler(c *gin.Context) {
	panic("not implemented")
}


func (h *teacherHandler) SuspendHandler(c *gin.Context) {
	panic("not implemented")
}


func (h *teacherHandler) RetrieveForNotificationsHandler(c *gin.Context) {
	panic("not implemented")
}