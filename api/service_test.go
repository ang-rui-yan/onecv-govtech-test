package api

import (
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

	mock.ExpectQuery("SELECT id FROM students WHERE email = \\$1").
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
// func TestGetTeacherID(t *testing.T) {
// 	mock, err := pgxmock.NewPool()

//     if err != nil {
//         t.Errorf("error creating stub connection: %v\n", err)
//     }
//     defer mock.Close()

//     db := NewDatabase(mock)
// 	svc := NewTeacherService(db)

// 	teacherEmail := "teacherken@gmail.com"

// 	// Begin transaction mock
// 	mock.ExpectBegin()

// 	// Mocking the teacher ID query
// 	mock.ExpectQuery("SELECT id FROM teachers WHERE email = $1").
// 		WithArgs(teacherEmail).
// 		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))

// 	// Commit transaction mock
// 	mock.ExpectCommit()

// 	// Call the service method
// 	id, err := svc.GetTeacherID(teacherEmail)

// 	// Assert there was no error
// 	assert.NoError(t, err)
// 	assert.Equal(t, 1, id)

// 	// Ensure all expectations were met
// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("there were unfulfilled expectations: %s", err)
// 	}

// }

// func TestRegister(t *testing.T) {
// 	mock, err := pgxmock.NewPool()

//     if err != nil {
//         t.Errorf("error creating stub connection: %v\n", err)
//     }
//     defer mock.Close()

//     db := NewDatabase(mock)
// 	svc := NewTeacherService(db)

// 	teacherEmail := "teacherken@gmail.com"
// 	studentEmails := []string{"studentjon@gmail.com", "studenthon@gmail.com"}

// 	// Begin transaction mock
// 	mock.ExpectBegin()

// 	// Mocking the teacher ID query
// 	mock.ExpectQuery("SELECT id FROM teachers WHERE email = $1").
// 		WithArgs(teacherEmail).
// 		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))

// 	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM students WHERE email = $1")).
// 		WithArgs("studentjon@gmail.com").
// 		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))

// 	// Mock the second student ID fetch.
// 	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM students WHERE email = $1")).
// 		WithArgs("studenthon@gmail.com").
// 		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(2))

// 	// Mocking the INSERT operation
// 	mock.ExpectExec("INSERT INTO teacher_students").WithArgs(1, 1).WillReturnResult(pgxmock.NewResult("INSERT", 1))
// 	mock.ExpectExec("INSERT INTO teacher_students").WithArgs(1, 2).WillReturnResult(pgxmock.NewResult("INSERT", 1))

// 	// Commit transaction mock
// 	mock.ExpectCommit()

// 	// Call the service method
// 	err = svc.Register(teacherEmail, studentEmails)

// 	// Assert there was no error
// 	assert.NoError(t, err)

// 	// Ensure all expectations were met
// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("there were unfulfilled expectations: %s", err)
// 	}

// }

// func TestRegister(t *testing.T) {
// 	// Initialize the mock for the PosgresInterface
// 	dbMock := new(mocks.PosgresInterface)

// 	// Set up your service, injecting the mock
// 	service := NewTeacherService(dbMock)

// 	// Define your test case
// 	testCases := []struct {
// 		name           string
// 		teacherEmail   string
// 		studentEmails  []string
// 		mockSetup      func()
// 		expectedError  error
// 	}{
// 		{
// 			name:         "successful registration",
// 			teacherEmail: "teacher@example.com",
// 			studentEmails: []string{
// 				"student1@example.com",
// 				"student2@example.com",
// 			},
// 			mockSetup: func() {
// 				dbMock.On("RegisterStudentsToTeacher", mock.Anything, "teacher@example.com", []string{"student1@example.com", "student2@example.com"}).Return(nil)
// 			},
// 			expectedError: nil,
// 		},
// 		{
// 			name:         "failed registration due to DB error",
// 			teacherEmail: "teacher@example.com",
// 			studentEmails: []string{
// 				"student1@example.com",
// 				"student2@example.com",
// 			},
// 			mockSetup: func() {
// 				dbMock.On("RegisterStudentsToTeacher", mock.Anything, "teacher@example.com", []string{"student1@example.com", "student2@example.com"}).Return(errors.New("db error"))
// 			},
// 			expectedError: errors.New("db error"),
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			// Set up the mock behavior for this test case
// 			tc.mockSetup()

// 			// Call the method under test
// 			err := service.Register(tc.teacherEmail, tc.studentEmails)

// 			// Assert the expectations
// 			if tc.expectedError != nil {
// 				assert.EqualError(t, err, tc.expectedError.Error())
// 			} else {
// 				assert.NoError(t, err)
// 			}

// 			// Verify that the expectations were met
// 			dbMock.AssertExpectations(t)
// 		})
// 	}
// }
