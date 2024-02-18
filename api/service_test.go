package api

import (
	"errors"
	"regexp"
	"testing"

	"github.com/jackc/pgx/v5"
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

	t.Run("successful registration", func(t *testing.T) {
		mock.ExpectBegin()

		// mock teacher
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM teachers WHERE email = $1")).
			WithArgs(teacherEmail).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))

		// mock first student
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM students WHERE email = $1")).
			WithArgs("studentjon@gmail.com").
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectExec("INSERT INTO teacher_students").WithArgs(1, 1).WillReturnResult(pgxmock.NewResult("INSERT", 1))

		// mock second student
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM students WHERE email = $1")).
			WithArgs("studenthon@gmail.com").
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(2))
		mock.ExpectExec("INSERT INTO teacher_students").WithArgs(1, 2).WillReturnResult(pgxmock.NewResult("INSERT", 1))

		mock.ExpectCommit()

		err := teacherService.RegisterStudentsToTeacher(teacherEmail, studentEmails)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("registration fails due to missing teacher", func(t *testing.T) {
		mock.ExpectBegin()

		// Simulate teacher not found
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM teachers WHERE email = $1")).
			WithArgs(teacherEmail).
			WillReturnError(pgx.ErrNoRows) 

		err := teacherService.RegisterStudentsToTeacher(teacherEmail, studentEmails)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	
	t.Run("registration fails on student insertion", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM teachers WHERE email = $1")).
			WithArgs(teacherEmail).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM students WHERE email = $1")).
			WithArgs(studentEmails[0]).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
		// simulate an error on INSERT
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO teacher_students")).
			WithArgs(1, 1).
			WillReturnError(errors.New("insert error"))
		// transaction should be rolled back, not committed, due to the error
		mock.ExpectRollback()

		err := teacherService.RegisterStudentsToTeacher(teacherEmail, studentEmails)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("duplicate registration", func(t *testing.T) {

		studentEmails := []string{"studentjon@gmail.com", "studentjon@gmail.com"}
		mock.ExpectBegin()

		// mock teacher
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM teachers WHERE email = $1")).
			WithArgs(teacherEmail).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))

		// mock first student
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM students WHERE email = $1")).
			WithArgs("studentjon@gmail.com").
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectExec("INSERT INTO teacher_students").WithArgs(1, 1).WillReturnResult(pgxmock.NewResult("INSERT", 1))

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM students WHERE email = $1")).
			WithArgs("studentjon@gmail.com").
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectExec("INSERT INTO teacher_students").WithArgs(1, 1).WillReturnResult(pgxmock.NewResult("INSERT", 0))

		mock.ExpectCommit()

		err := teacherService.RegisterStudentsToTeacher(teacherEmail, studentEmails)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())

	})

}