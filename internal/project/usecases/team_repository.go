package usecases

import "github.com/lunarKettle/task-management-platform-monolith/internal/project/models"

type TeamRepository interface {
	Create(team models.Team) (uint32, error)
	Update(team models.Team) error
	Delete(teamId uint32) error
	GetById(teamId uint32) (models.Team, error)
}
