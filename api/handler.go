package api

type teacherHandler struct {
    Service TeacherService
}

func NewTeacherHandler(svc TeacherService) teacherHandler{
    return teacherHandler{Service: svc}
}
