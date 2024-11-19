package usecases

import "github.com/lunarKettle/task-management-platform-monolith/internal/aggregate/project/models"

type ProjectRepository interface {
	CreateProject(project *models.Project) (uint32, error)
	UpdateProject(project *models.Project) error
	DeleteProject(projectId uint32) error
	GetProjectById(projectIdS uint32) (*models.Project, error)

	CreateTeam(team *models.Team) (uint32, error)
	UpdateTeam(team *models.Team) error
	DeleteTeam(teamId uint32) error
	GetTeamById(teamId uint32) (*models.Team, error)
}
