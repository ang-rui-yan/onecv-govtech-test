package api

type TeacherService interface {

}

type teacherService struct {
}

func NewTeacherService () *teacherService {
	return &teacherService{}
}

func (s *teacherService) Register(teacherEmail string, studentEmails []string) (error) {
	return nil
}