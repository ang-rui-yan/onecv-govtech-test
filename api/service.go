package api

type TeacherService interface {
	// Register(teacherEmail string, studentEmails []string) (error)
	// GetTeacherID(teacherEmail string) (int, error)
	GetStudentID(studentEmail string) (int, error)
}

type teacherService struct {
	db Database
}

func NewTeacherService (db Database) *teacherService {
	return &teacherService{db: db}
}

func (s *teacherService) GetStudentID(studentEmail string) (int, error) {
	id, err := s.db.GetStudentID(studentEmail)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *teacherService) GetTeacherID(teacherEmail string) (int, error) {
	id, err := s.db.GetTeacherID(teacherEmail)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *teacherService) RegisterStudentsToTeacher(teacherEmail string, studentEmails []string) (error) {
	err := s.db.RegisterStudentsToTeacher(teacherEmail, studentEmails)
	if err != nil {
		return err
	}

	return nil
}