package api

import "errors"

var (
	ErrStudentAlreadySuspended = errors.New("student is already suspended")
	ErrStudentNotFound = errors.New("student not found")
	ErrTeacherNotFound = errors.New("teacher not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrFailedToUpdateSuspension = errors.New("failed to update student suspension status")
	ErrBadRequest = errors.New("bad request")
)