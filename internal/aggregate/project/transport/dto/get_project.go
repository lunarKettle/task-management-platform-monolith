package dto

import "github.com/lunarKettle/task-management-platform-monolith/internal/aggregate/project/models"

type GetProjectResponseDTO struct {
	ID      uint32          `json:"id"`
	Name    string          `json:"name"`
	Members []models.Member `json:"members"`
}
