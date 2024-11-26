package usecases

import (
	"errors"
	"fmt"
	"time"

	"github.com/lunarKettle/task-management-platform-monolith/internal/aggregate/project/models"
	"github.com/lunarKettle/task-management-platform-monolith/pkg/common"
)

type ProjectUseCases struct {
	repo ProjectRepository
}

func NewProjectUseCases(repo ProjectRepository) *ProjectUseCases {
	return &ProjectUseCases{
		repo: repo,
	}
}

// Запрос для получения проекта по ID
type GetProjectByIDQuery struct {
	id uint32
}

func NewGetProjectByIDQuery(id uint32) *GetProjectByIDQuery {
	return &GetProjectByIDQuery{id: id}
}

func (p *ProjectUseCases) GetProjectByID(query *GetProjectByIDQuery) (*models.Project, error) {
	project, err := p.repo.GetProjectById(query.id)
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
		Team:           models.Team{ID: cmd.teamId},
		Budget:         cmd.budget,
	}

	id, err := p.repo.CreateProject(project)
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
	_, err := p.repo.GetProjectById(cmd.id)

	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return fmt.Errorf("project with id %d is not found: %w", cmd.id, err)
		}
		return err
	}

	project := &models.Project{
		Id:             cmd.id,
		Name:           cmd.name,
		Description:    cmd.description,
		StartDate:      cmd.startDate,
		PlannedEndDate: cmd.plannedEndDate,
		ActualEndDate:  cmd.actualEndDate,
		Status:         cmd.status,
		Priority:       cmd.priority,
		Team:           models.Team{ID: cmd.teamId},
		Budget:         cmd.budget,
	}

	if err := p.repo.UpdateProject(project); err != nil {
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
	if err := p.repo.DeleteProject(cmd.id); err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}
	return nil
}

// Запрос для получения команды по ID
type GetTeamByIDQuery struct {
	id uint32
}

func NewGetTeamByIDQuery(id uint32) *GetProjectByIDQuery {
	return &GetProjectByIDQuery{id: id}
}

func (p *ProjectUseCases) GetTeamByID(query *GetProjectByIDQuery) (*models.Team, error) {
	team, err := p.repo.GetTeamById(query.id)
	if err != nil {
		return nil, fmt.Errorf("failed to get team by id: %w", err)
	}
	return team, nil
}

// Команда для создания команды
type CreateTeamCommand struct {
	name      string
	members   []Member
	managerID uint32
}

func NewCreateTeamCommand(name string, members []Member, managerID uint32) *CreateTeamCommand {
	return &CreateTeamCommand{
		name:      name,
		members:   members,
		managerID: managerID,
	}
}

func (p *ProjectUseCases) CreateTeam(cmd *CreateTeamCommand) (uint32, error) {
	team := &models.Team{
		Name:      cmd.name,
		Members:   mapMembersToModels(cmd.members),
		ManagerID: cmd.managerID,
	}

	id, err := p.repo.CreateTeam(team)
	if err != nil {
		return 0, fmt.Errorf("failed to create team: %w", err)
	}
	return id, nil
}

// Команда для обновления команды
type UpdateTeamCommand struct {
	id        uint32
	name      string
	members   []Member
	managerID uint32
}

func NewUpdateTeamCommand(
	id uint32,
	name string,
	members []Member,
	managerID uint32) *UpdateTeamCommand {
	return &UpdateTeamCommand{
		id:        id,
		name:      name,
		members:   members,
		managerID: managerID,
	}
}

func (p *ProjectUseCases) UpdateTeam(cmd *UpdateTeamCommand) error {
	_, err := p.repo.GetTeamById(cmd.id)

	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return fmt.Errorf("team with id %d is not found: %w", cmd.id, err)
		}
		return fmt.Errorf("failed to get user with id %d: %w", cmd.id, err)
	}

	team := &models.Team{
		ID:        cmd.id,
		Name:      cmd.name,
		Members:   mapMembersToModels(cmd.members),
		ManagerID: cmd.managerID,
	}

	if err := p.repo.UpdateTeam(team); err != nil {
		return fmt.Errorf("failed to update team: %w", err)
	}
	return nil
}

// Команда для удаления команды
type DeleteTeamCommand struct {
	id uint32
}

func NewDeleteTeamCommand(id uint32) *DeleteTeamCommand {
	return &DeleteTeamCommand{id: id}
}

func (p *ProjectUseCases) DeleteTeam(cmd *DeleteTeamCommand) error {
	if err := p.repo.DeleteTeam(cmd.id); err != nil {
		return fmt.Errorf("failed to delete team: %w", err)
	}
	return nil
}
