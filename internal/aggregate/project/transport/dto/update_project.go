package dto

import "time"

type UpdateProjectRequestDTO struct {
	ID             uint32    `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	StartDate      time.Time `json:"startDate"`
	PlannedEndDate time.Time `json:"plannedEndDate"`
	ActualEndDate  time.Time `json:"actualEndDate"`
	Status         string    `json:"status"`
	Priority       uint32    `json:"priority"`
	TeamId         uint32    `json:"teamId"`
	Budget         float64   `json:"budget"`
}
