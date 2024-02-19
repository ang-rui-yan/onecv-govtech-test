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
	Teacher string `json:"teacher" binding:"required"`
	Students []string `json:"students" binding:"required"`
}

type SuspendRequestBody struct {
	StudentEmail string `json:"student" binding:"required"`
}