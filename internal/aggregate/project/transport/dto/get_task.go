package dto

type GetTaskResponseDTO struct {
	ID          uint32 `json:"id"`
	Description string `json:"description"`
	EmployeeID  uint32 `json:"employee_id"`
	ProjectID   uint32 `json:"project_id"`
	IsCompleted bool   `json:"is_completed"`
}
