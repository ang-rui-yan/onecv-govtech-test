package api

import (
	"regexp"
	"testing"

	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
)

func TestGetStudentID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	database := NewDatabase(mock)
	teacherService := NewTeacherService(database)
	studentEmail := "studentjon@gmail.com"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM students WHERE email = $1")).
		WithArgs(studentEmail).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))

	id, err := teacherService.GetStudentID(studentEmail)
	assert.NoError(t, err, "Error was not expected while getting student ID")
	assert.Equal(t, 1, id, "Expected student ID does not match the actual ID")

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetTeacherID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	database := NewDatabase(mock)
	teacherService := NewTeacherService(database)
	teacherEmail := "teacherken@gmail.com"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM teachers WHERE email = $1")).
		WithArgs(teacherEmail).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))

	id, err := teacherService.GetTeacherID(teacherEmail)
	assert.NoError(t, err, "Error was not expected while getting teacher ID")
	assert.Equal(t, 1, id, "Expected teacher ID does not match the actual ID")

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRegisterStudentsToTeacher(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	database := NewDatabase(mock)
	teacherService := NewTeacherService(database)

	teacherEmail := "teacherken@gmail.com"
	studentEmails := []string{"studentjon@gmail.com", "studenthon@gmail.com"}

	// Begin transaction mock
	mock.ExpectBegin()

	// Mocking the teacher ID query
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM teachers WHERE email = $1")).
		WithArgs(teacherEmail).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM students WHERE email = $1")).
		WithArgs("studentjon@gmail.com").
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectExec("INSERT INTO teacher_students").WithArgs(1, 1).WillReturnResult(pgxmock.NewResult("INSERT", 1))

	// Mock the second student ID fetch.
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM students WHERE email = $1")).
		WithArgs("studenthon@gmail.com").
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(2))

	// Mocking the INSERT operation
	mock.ExpectExec("INSERT INTO teacher_students").WithArgs(1, 2).WillReturnResult(pgxmock.NewResult("INSERT", 1))

	// Commit transaction mock
	mock.ExpectCommit()

	// Call the service method
	err = teacherService.RegisterStudentsToTeacher(teacherEmail, studentEmails)

	// Assert there was no error
	assert.NoError(t, err)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}