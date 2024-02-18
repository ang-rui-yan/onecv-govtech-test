package api

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

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
	query := `
			SELECT s.email 
			FROM student s
			INNER JOIN teacher_students ts ON s.id = ts.student_id
			INNER JOIN teacher t ON ts.teacher_id = t.id 
			WHERE t.email = ANY($1)
			GROUP BY s.email
			HAVING COUNT(DISTINCT t.id) = (SELECT COUNT(DISTINCT id) FROM teacher WHERE email = ANY($1))`

	rows, err := pool.DB.Query(context.Background(), query, teacherEmails)
	if err != nil {
		return nil, fmt.Errorf("could not find students for teachers %v: %v", teacherEmails, err)
	}
	defer rows.Close()

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
	ctx := context.Background()

	var suspended bool
	err := pool.DB.QueryRow(ctx, "SELECT suspended FROM students WHERE email = $1", studentEmail).
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