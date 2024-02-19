package api

type Student struct {
	ID int
	Email string
	Suspended bool
}

type Teacher struct {
	ID int
	Email string
}

type RegistrationRequestBody struct {
	TeacherEmail string `json:"teacher" binding:"required"`
	StudentEmails []string `json:"students" binding:"required"`
}

type SuspendRequestBody struct {
	StudentEmail string `json:"student" binding:"required"`
}

type NotificationRequestBody struct {
	TeacherEmail string `json:"teacher" binding:"required"`
	Notification string `json:"notification" binding:"required"`
}