package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	svc := NewTeacherService()

	t.Run("Register one student to a teacher", func(t *testing.T) {
		teacherEmail := ""
		studentEmails := []string {
			"",
		}
		err := svc.Register(teacherEmail, studentEmails)

		assert.NoError(t, err)
	})
	
}