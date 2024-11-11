package usecases

import "github.com/lunarKettle/task-management-platform-monolith/internal/project/models"

type ProjectRepository interface {
	Create(project *models.Project) (uint32, error)
	Update(project *models.Project) error
	Delete(projectId uint32) error
	GetById(projectIdS uint32) (*models.Project, error)
}
