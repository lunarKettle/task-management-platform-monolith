package models

type Task struct {
	ID          uint32
	Description string
	EmployeeID  uint32
	ProjectID   uint32
	IsCompleted bool
}
