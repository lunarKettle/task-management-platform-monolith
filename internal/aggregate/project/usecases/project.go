package usecases

import (
	"fmt"
	"time"

	"github.com/lunarKettle/task-management-platform-monolith/internal/aggregate/project/models"
)

type ProjectUseCases struct {
	repo ProjectRepository
}

func NewProjectUseCases(repo ProjectRepository) *ProjectUseCases {
	return &ProjectUseCases{
		repo: repo,
	}
}

// Команда для получения проекта по ID
type GetProjectByIDQuery struct {
	id uint32
}

func NewGetProjectByIDQuery(id uint32) *GetProjectByIDQuery {
	return &GetProjectByIDQuery{id: id}
}

func (p *ProjectUseCases) GetProjectByID(query *GetProjectByIDQuery) (*models.Project, error) {
	project, err := p.repo.GetById(query.id)
	if err != nil {
		return nil, fmt.Errorf("failed to get project by id: %w", err)
	}
	return project, nil
}

// Команда для создания проекта
type CreateProjectCommand struct {
	name           string
	description    string
	plannedEndDate time.Time
	status         string
	priority       uint32
	teamId         uint32
	budget         float64
}

func NewCreateProjectCommand(
	name string,
	description string,
	plannedEndDate time.Time,
	status string,
	priority uint32,
	teamId uint32,
	budget float64) *CreateProjectCommand {
	return &CreateProjectCommand{
		name:           name,
		description:    description,
		plannedEndDate: plannedEndDate,
		status:         status,
		priority:       priority,
		teamId:         teamId,
		budget:         budget,
	}
}

func (p *ProjectUseCases) CreateProject(cmd *CreateProjectCommand) (uint32, error) {
	project := &models.Project{
		Name:           cmd.name,
		Description:    cmd.description,
		StartDate:      time.Now(),
		PlannedEndDate: cmd.plannedEndDate,
		ActualEndDate:  time.Time{},
		Status:         cmd.status,
		Priority:       cmd.priority,
		TeamId:         cmd.teamId,
		Budget:         cmd.budget,
	}

	id, err := p.repo.Create(project)
	if err != nil {
		return 0, fmt.Errorf("failed to create project: %w", err)
	}
	return id, nil
}

// Команда для обновления проекта
type UpdateProjectCommand struct {
	id             uint32
	name           string
	description    string
	startDate      time.Time
	plannedEndDate time.Time
	actualEndDate  time.Time
	status         string
	priority       uint32
	teamId         uint32
	budget         float64
}

func NewUpdateProjectCommand(
	id uint32,
	name string,
	description string,
	startDate time.Time,
	plannedEndDate time.Time,
	actualEndDate time.Time,
	status string,
	priority uint32,
	teamId uint32,
	budget float64) *UpdateProjectCommand {
	return &UpdateProjectCommand{
		id:             id,
		name:           name,
		description:    description,
		startDate:      startDate,
		plannedEndDate: plannedEndDate,
		actualEndDate:  actualEndDate,
		status:         status,
		priority:       priority,
		teamId:         teamId,
		budget:         budget,
	}
}

func (p *ProjectUseCases) UpdateProject(cmd *UpdateProjectCommand) error {
	project := &models.Project{
		Id:             cmd.id,
		Name:           cmd.name,
		Description:    cmd.description,
		StartDate:      cmd.startDate,
		PlannedEndDate: cmd.plannedEndDate,
		ActualEndDate:  cmd.actualEndDate,
		Status:         cmd.status,
		Priority:       cmd.priority,
		TeamId:         cmd.teamId,
		Budget:         cmd.budget,
	}

	if err := p.repo.Update(project); err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}
	return nil
}

// Команда для удаления проекта
type DeleteProjectCommand struct {
	id uint32
}

func NewDeleteProjectCommand(id uint32) *DeleteProjectCommand {
	return &DeleteProjectCommand{id: id}
}

func (p *ProjectUseCases) DeleteProject(cmd *DeleteProjectCommand) error {
	if err := p.repo.Delete(cmd.id); err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}
	return nil
}
