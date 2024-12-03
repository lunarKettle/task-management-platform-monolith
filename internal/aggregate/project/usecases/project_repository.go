package usecases

import "github.com/lunarKettle/task-management-platform-monolith/internal/aggregate/project/models"

type ProjectRepository interface {
	CreateProject(project *models.Project) (uint32, error)
	UpdateProject(project *models.Project) error
	DeleteProject(projectID uint32) error
	GetProjectById(projectID uint32) (*models.Project, error)

	CreateTeam(team *models.Team) (uint32, error)
	UpdateTeam(team *models.Team) error
	DeleteTeam(teamID uint32) error
	GetTeamById(teamID uint32) (*models.Team, error)
	GetTeamIdByUserID(userID uint32) (uint32, error)
}
