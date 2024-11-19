package dto

import "time"

type CreateProjectRequestDTO struct {
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	PlannedEndDate time.Time `json:"plannedEndDate"`
	Status         string    `json:"status"`
	Priority       uint32    `json:"priority"`
	TeamId         uint32    `json:"teamId"`
	Budget         float64   `json:"budget"`
}
