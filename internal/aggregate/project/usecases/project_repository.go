package usecases

import "github.com/lunarKettle/task-management-platform-monolith/internal/aggregate/project/models"

type ProjectRepository interface {
	CreateProject(project *models.Project) (uint32, error)
	UpdateProject(project *models.Project) error
	DeleteProject(projectID uint32) error
	GetAllProjects() ([]*models.Project, error)
	GetProjectById(projectID uint32) (*models.Project, error)

	CreateTeam(team *models.Team) (uint32, error)
	UpdateTeam(team *models.Team) error
	DeleteTeam(teamID uint32) error
	GetTeamById(teamID uint32) (*models.Team, error)
	GetTeamIdByUserID(userID uint32) (uint32, error)

	CreateTask(task *models.Task) (uint32, error)
	UpdateTask(task *models.Task) error
	DeleteTask(taskID uint32) error
	GetTaskById(taskID uint32) (*models.Task, error)
	GetTasksByEmployeeID(employeeID uint32) ([]*models.Task, error)
}
