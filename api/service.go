package api

import (
	"errors"
	"studentadmin/utils"
)

type TeacherService interface {
	GetStudentID(studentEmail string) (int, error)
	GetTeacherID(teacherEmail string) (int, error)
	RegisterStudentsToTeacher(teacherEmail string, studentEmails []string) (error)
	GetCommonStudents(teacherEmails []string) ([]string, error)
	Suspend(studentEmail string) error
	RetrieveForNotifications(teacherEmail string, notification string) ([]string, error)
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
		switch err {
		case ErrTeacherNotFound:
			return ErrTeacherNotFound
		case ErrStudentNotFound:
			return ErrStudentNotFound
		default:
			return err
		}
	}

	return nil
}

func (s *teacherService) GetCommonStudents(teacherEmails []string) ([]string, error) {
	studentEmails, err := s.db.GetCommonStudents(teacherEmails)
	if err != nil {
		return nil, err
	}

	return studentEmails, nil
}

func (s *teacherService) Suspend(studentEmail string) error {
	if studentEmail == "" {
		return ErrInvalidInput
	}

	_, err := s.db.GetStudentID(studentEmail)
	if errors.Is(err, ErrStudentNotFound) {
		return ErrStudentNotFound
	} else if err != nil {
		return err
	}

	err = s.db.Suspend(studentEmail)
	if errors.Is(err, ErrStudentAlreadySuspended) {
		return ErrStudentAlreadySuspended
	} else if err != nil {
		return err
	}

	return nil
}

func (s *teacherService) RetrieveForNotifications(teacherEmail string, notification string) ([]string, error) {
	mentions, err := utils.ExtractEmailMentions(notification)
	if err != nil {
		return nil, err
	}

    // Check if the teacher exists
    _, err = s.db.GetTeacherID(teacherEmail)
    if err != nil {
        return nil, err
    }

	recipients, err := s.db.RetrieveForNotifications(teacherEmail, mentions)
	if err != nil {
		return nil, err
	}

	recipients = utils.RemoveDuplicates(recipients)

	return recipients, nil
}