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
		return 0, fmt.Errorf("could not find student with email %s: %v", studentEmail, err)
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
		_, err = tx.Exec(ctx, "INSERT INTO teacher_students (teacher_id, student_id) VALUES ($1, $2)", teacherID, studentID)
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