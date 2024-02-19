package api

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// interface for mocks and actual
type PgxPoolIface interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	Ping(context.Context) error
	Close()
}

type Database struct {
    DB PgxPoolIface
}

func NewDatabase(ds PgxPoolIface) Database {
    return Database{DB: ds}
}

// get student id from email
func (pool Database) GetStudentID (studentEmail string) (int, error) {
	query := "SELECT id FROM students WHERE email = $1"

	row := pool.DB.QueryRow(context.Background(), query, studentEmail)

	var studentID int
	err := row.Scan(&studentID)
	if err != nil {
		return 0, fmt.Errorf("could not find student with email %s", studentEmail)
	}

	return studentID, nil
}

// Get teacher id from email
func (pool Database) GetTeacherID(teacherEmail string) (int, error) {
	query := "SELECT id FROM teachers WHERE email = $1"

	row := pool.DB.QueryRow(context.Background(), query, teacherEmail)

	var teacherID int
	err := row.Scan(&teacherID)
	if err != nil {
		return 0, fmt.Errorf("could not find teacher with email %s: %v", teacherEmail, err)
	}

	return teacherID, nil
}

func (pool Database) RegisterStudentsToTeacher(teacherEmail string, studentEmails []string) error {
	// begin transaction
	ctx := context.Background()
	tx, err := pool.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Get the teacher ID
	var teacherID int
	err = tx.QueryRow(ctx, "SELECT id FROM teachers WHERE email = $1", teacherEmail).Scan(&teacherID)
	if err != nil {
		return fmt.Errorf("could not find teacher with email %s: %v", teacherEmail, err)
	}

	for _, studentEmail := range studentEmails {
		// Get the student ID
		var studentID int
		err = tx.QueryRow(ctx, "SELECT id FROM students WHERE email = $1", studentEmail).Scan(&studentID)
		if err != nil {
			return fmt.Errorf("could not find student with email %s: %v", studentEmail, err)
		}

		// Prepare the insert statement
		query := "INSERT INTO teacher_students (teacher_id, student_id) VALUES ($1, $2) ON CONFLICT (teacher_id, student_id) DO NOTHING"
		_, err = tx.Exec(ctx, query, teacherID, studentID)
		if err != nil {
			return fmt.Errorf("could not insert teacher-student relationship for teacher %s and student %s: %v", teacherEmail, studentEmail, err)
		}
	}

	// Commit the transaction
	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (pool Database) GetCommonStudents(teacherEmails []string) ([]string, error) {
	// join all the tables and see where the teacher email matches
	// get the distinct ones
	query := `
			SELECT s.email 
			FROM students s
			INNER JOIN teacher_students ts ON s.id = ts.student_id
			INNER JOIN teachers t ON ts.teacher_id = t.id 
			WHERE t.email = ANY($1)
			GROUP BY s.email
			HAVING COUNT(DISTINCT t.id) = (SELECT COUNT(DISTINCT id) FROM teachers WHERE email = ANY($1))`

	rows, err := pool.DB.Query(context.Background(), query, teacherEmails)
	if err != nil {
		return nil, fmt.Errorf("could not find students for teachers %v: %v", teacherEmails, err)
	}
	defer rows.Close()

	// loop through the students retrieved
	var students []string
	for rows.Next() {
		var studentEmail string
		if err := rows.Scan(&studentEmail); err != nil {
			return nil, err
		}
		students = append(students, studentEmail)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}


func (pool Database) Suspend(studentEmail string) error {
	var suspended bool

	// check if the student is already suspended
	err := pool.DB.QueryRow(context.Background(), "SELECT suspended FROM students WHERE email = $1", studentEmail).
			Scan(&suspended)
	if err != nil {
		return fmt.Errorf("could not query student suspension status: %v", err)
	}
	if suspended {
		return fmt.Errorf("student with email %s is already suspended", studentEmail)
	}
	
	query := `
		UPDATE students 
		SET suspended = true
		WHERE email = $1 AND NOT suspended`

	cmdTag, err := pool.DB.Exec(context.Background(), query, studentEmail)
	if err != nil {
		return fmt.Errorf("failed to suspend student: %v", err)
	}

    if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no action taken for student %s, they might have been suspended", studentEmail)
    }
	
	return nil
}

func (pool Database) RetrieveForNotifications(teacherEmail string, mentions []string) ([]string, error) {
	query := `
		SELECT DISTINCT email 
		FROM students 
		WHERE NOT suspended AND (
			email = ANY($1) OR
			id IN (
				SELECT student_id 
				FROM teacher_students 
				WHERE teacher_id = (SELECT id FROM teachers WHERE email = $2)
			)
		)`


	rows, err := pool.DB.Query(context.Background(), query, mentions, teacherEmail)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve students for notification: %v", err)
	}
	defer rows.Close()

	var recipients []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, fmt.Errorf("failed to scan student email: %v", err)
		}
		recipients = append(recipients, email)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed during rows iteration: %v", err)
	}

	return recipients, nil
}