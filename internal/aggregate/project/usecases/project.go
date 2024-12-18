package usecases

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/lunarKettle/task-management-platform-monolith/internal/aggregate/project/models"
	"github.com/lunarKettle/task-management-platform-monolith/pkg/common"
)

const adminRole string = "admin"

type ProjectUseCases struct {
	repo   ProjectRepository
	logger *log.Logger
	mu     sync.Mutex
}

func NewProjectUseCases(repo ProjectRepository) *ProjectUseCases {
	logFile, err := os.OpenFile("business_operations.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("failed to create log file: %v", err))
	}

	logger := log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	return &ProjectUseCases{
		repo:   repo,
		logger: logger,
	}
}

func (uc *ProjectUseCases) logMessage(ctx context.Context, message string) {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	claims := ctx.Value(common.ContextKeyClaims).(*common.Claims)
	uc.logger.Printf("UserID: %d, Role: %s - %s", claims.UserID, claims.Role, message)
}

// Запрос для получения всех проектов
func (uc *ProjectUseCases) GetAllProjects(ctx context.Context) ([]*models.Project, error) {
	claims := ctx.Value(common.ContextKeyClaims).(*common.Claims)
	projects, err := uc.repo.GetAllProjects()
	var result []*models.Project
	if err != nil {
		return nil, fmt.Errorf("failed to get all projects: %w", err)
	}

	if claims.Role != adminRole {
		teamID, err := uc.repo.GetTeamIdByUserID(claims.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to get team by userID %d: %w", claims.UserID, err)
		}

		for _, project := range projects {
			if project.Team.ID == teamID {
				result = append(result, project)
			}
		}

		return result, nil
	}

	uc.logMessage(ctx, "Fetching all projects")
	return projects, nil
}

// Запрос для получения проекта по ID
type GetProjectByIDQuery struct {
	id uint32
}

func NewGetProjectByIDQuery(id uint32) *GetProjectByIDQuery {
	return &GetProjectByIDQuery{id: id}
}

func (uc *ProjectUseCases) GetProjectByID(ctx context.Context, query *GetProjectByIDQuery) (*models.Project, error) {
	claims := ctx.Value(common.ContextKeyClaims).(*common.Claims)

	project, err := uc.repo.GetProjectById(query.id)
	if err != nil {
		return nil, fmt.Errorf("failed to get project by id: %w", err)
	}

	if claims.Role != adminRole {
		teamID, err := uc.repo.GetTeamIdByUserID(claims.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to get team by userID %d: %w", claims.UserID, err)
		}

		if project.Team.ID != teamID {
			return nil, common.ErrForbidden
		}
	}

	uc.logMessage(ctx, fmt.Sprintf("Fetching project by id: %d", query.id))
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

func (uc *ProjectUseCases) CreateProject(ctx context.Context, cmd *CreateProjectCommand) (uint32, error) {
	claims := ctx.Value(common.ContextKeyClaims).(*common.Claims)

	if claims.Role != adminRole {
		return 0, common.ErrForbidden
	}

	project := &models.Project{
		Name:           cmd.name,
		Description:    cmd.description,
		StartDate:      time.Now(),
		PlannedEndDate: cmd.plannedEndDate,
		ActualEndDate:  time.Time{},
		Status:         cmd.status,
		Priority:       cmd.priority,
		Team:           &models.Team{ID: cmd.teamId},
		Budget:         cmd.budget,
	}

	id, err := uc.repo.CreateProject(project)
	if err != nil {
		return 0, fmt.Errorf("failed to create project: %w", err)
	}

	uc.logMessage(ctx, fmt.Sprintf("Creating new project (id: %d)", id))
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

func (uc *ProjectUseCases) UpdateProject(ctx context.Context, cmd *UpdateProjectCommand) error {
	claims := ctx.Value(common.ContextKeyClaims).(*common.Claims)

	if claims.Role != adminRole {
		return common.ErrForbidden
	}

	_, err := uc.repo.GetProjectById(cmd.id)

	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return fmt.Errorf("project with id %d is not found: %w", cmd.id, err)
		}
		return fmt.Errorf("failed to get project with id %d: %w", cmd.id, err)
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
		Team:           &models.Team{ID: cmd.teamId},
		Budget:         cmd.budget,
	}

	if err := uc.repo.UpdateProject(project); err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}

	uc.logMessage(ctx, fmt.Sprintf("Updating project (id: %d)", cmd.id))
	return nil
}

// Команда для удаления проекта
type DeleteProjectCommand struct {
	id uint32
}

func NewDeleteProjectCommand(id uint32) *DeleteProjectCommand {
	return &DeleteProjectCommand{id: id}
}

func (uc *ProjectUseCases) DeleteProject(ctx context.Context, cmd *DeleteProjectCommand) error {
	claims := ctx.Value(common.ContextKeyClaims).(*common.Claims)

	if claims.Role != adminRole {
		return common.ErrForbidden
	}

	if err := uc.repo.DeleteProject(cmd.id); err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	uc.logMessage(ctx, fmt.Sprintf("Deleting project (id: %d)", cmd.id))
	return nil
}

func (uc *ProjectUseCases) GetAllTeams(ctx context.Context) ([]*models.Team, error) {
	claims := ctx.Value(common.ContextKeyClaims).(*common.Claims)

	teams, err := uc.repo.GetAllTeams()
	var result []*models.Team
	if err != nil {
		return nil, fmt.Errorf("failed to get all teams: %w", err)
	}

	if claims.Role != adminRole {
		for _, team := range teams {
			for _, member := range team.Members {
				if member.ID == claims.UserID {
					result = append(result, team)
				}
			}
		}
		return result, nil
	}

	uc.logMessage(ctx, "Fetching all teams")
	return teams, nil
}

// Запрос для получения команды по ID
type GetTeamByIDQuery struct {
	id uint32
}

func NewGetTeamByIDQuery(id uint32) *GetProjectByIDQuery {
	return &GetProjectByIDQuery{id: id}
}

func (uc *ProjectUseCases) GetTeamByID(ctx context.Context, query *GetProjectByIDQuery) (*models.Team, error) {
	team, err := uc.repo.GetTeamById(query.id)
	if err != nil {
		return nil, fmt.Errorf("failed to get team by id: %w", err)
	}

	uc.logMessage(ctx, fmt.Sprintf("Fetching team (id: %d)", query.id))
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

func (uc *ProjectUseCases) CreateTeam(ctx context.Context, cmd *CreateTeamCommand) (uint32, error) {
	claims := ctx.Value(common.ContextKeyClaims).(*common.Claims)

	if claims.Role != adminRole {
		return 0, common.ErrForbidden
	}

	team := &models.Team{
		Name:      cmd.name,
		Members:   mapMembersToModels(cmd.members),
		ManagerID: cmd.managerID,
	}

	id, err := uc.repo.CreateTeam(team)
	if err != nil {
		return 0, fmt.Errorf("failed to create team: %w", err)
	}

	uc.logMessage(ctx, fmt.Sprintf("Creating new team (id: %d)", id))
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

func (uc *ProjectUseCases) UpdateTeam(ctx context.Context, cmd *UpdateTeamCommand) error {
	claims := ctx.Value(common.ContextKeyClaims).(*common.Claims)

	if claims.Role != adminRole {
		return common.ErrForbidden
	}

	_, err := uc.repo.GetTeamById(cmd.id)

	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return fmt.Errorf("team with id %d is not found: %w", cmd.id, err)
		}
		return fmt.Errorf("failed to get user with id %d: %w", cmd.id, err)
	}

	for _, value := range cmd.members {
		member, err := uc.repo.GetMember(value.id)

		if err != nil {
			return fmt.Errorf("failed to get member by id %d: %w", value.id, err)
		}

		if member.Name != value.name {
			return common.ErrForbidden
		}
	}

	team := &models.Team{
		ID:        cmd.id,
		Name:      cmd.name,
		Members:   mapMembersToModels(cmd.members),
		ManagerID: cmd.managerID,
	}

	if err := uc.repo.UpdateTeam(team); err != nil {
		return fmt.Errorf("failed to update team: %w", err)
	}

	uc.logMessage(ctx, fmt.Sprintf("Updating team (id: %d)", cmd.id))
	return nil
}

// Команда для удаления команды
type DeleteTeamCommand struct {
	id uint32
}

func NewDeleteTeamCommand(id uint32) *DeleteTeamCommand {
	return &DeleteTeamCommand{id: id}
}

func (uc *ProjectUseCases) DeleteTeam(ctx context.Context, cmd *DeleteTeamCommand) error {
	claims := ctx.Value(common.ContextKeyClaims).(*common.Claims)

	if claims.Role != adminRole {
		return common.ErrForbidden
	}

	if err := uc.repo.DeleteTeam(cmd.id); err != nil {
		return fmt.Errorf("failed to delete team: %w", err)
	}

	uc.logMessage(ctx, fmt.Sprintf("Deleting team (id: %d)", cmd.id))
	return nil
}

type MemberFilter struct {
	Role   string
	TeamID uint32
}

func (uc *ProjectUseCases) GetMembers(ctx context.Context, filter MemberFilter) ([]*models.Member, error) {
	members, err := uc.repo.GetMembers(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get members: %w", err)
	}

	uc.logMessage(ctx, "Fetching teams members")
	return members, nil
}

// Запрос для получения задачи по ID
type GetTaskByIDQuery struct {
	id uint32
}

func NewGetTaskByIDQuery(id uint32) *GetTaskByIDQuery {
	return &GetTaskByIDQuery{id: id}
}

func (uc *ProjectUseCases) GetTaskByID(ctx context.Context, query *GetTaskByIDQuery) (*models.Task, error) {
	task, err := uc.repo.GetTaskById(query.id)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return nil, fmt.Errorf("task with id %d is not found: %w", query.id, err)
		}
		return nil, fmt.Errorf("failed to get task by id: %w", err)
	}

	uc.logMessage(ctx, fmt.Sprintf("Fetching task (id: %d)", query.id))
	return task, nil
}

// Команда для создания задачи
type CreateTaskCommand struct {
	description string
	employeeID  uint32
	projectID   uint32
	isCompleted bool
}

func NewCreateTaskCommand(description string, employeeID, projectID uint32, isCompleted bool) *CreateTaskCommand {
	return &CreateTaskCommand{
		description: description,
		employeeID:  employeeID,
		projectID:   projectID,
		isCompleted: isCompleted,
	}
}

func (uc *ProjectUseCases) CreateTask(ctx context.Context, cmd *CreateTaskCommand) (uint32, error) {
	claims := ctx.Value(common.ContextKeyClaims).(*common.Claims)

	if claims.Role != adminRole {
		return 0, common.ErrForbidden
	}

	task := &models.Task{
		Description: cmd.description,
		EmployeeID:  cmd.employeeID,
		ProjectID:   cmd.projectID,
		IsCompleted: cmd.isCompleted,
	}

	id, err := uc.repo.CreateTask(task)
	if err != nil {
		return 0, fmt.Errorf("failed to create task: %w", err)
	}

	uc.logMessage(ctx, fmt.Sprintf("Creating new task (id: %d)", id))
	return id, nil
}

// Команда для обновления задачи
type UpdateTaskCommand struct {
	id          uint32
	description string
	employeeID  uint32
	projectID   uint32
	isCompleted bool
}

func NewUpdateTaskCommand(id uint32, description string, employeeID, projectID uint32, isCompleted bool) *UpdateTaskCommand {
	return &UpdateTaskCommand{
		id:          id,
		description: description,
		employeeID:  employeeID,
		projectID:   projectID,
		isCompleted: isCompleted,
	}
}

func (uc *ProjectUseCases) UpdateTask(ctx context.Context, cmd *UpdateTaskCommand) error {
	claims := ctx.Value(common.ContextKeyClaims).(*common.Claims)

	if claims.Role != adminRole {
		return common.ErrForbidden
	}

	_, err := uc.repo.GetTaskById(cmd.id)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return fmt.Errorf("task with id %d is not found: %w", cmd.id, err)
		}
		return fmt.Errorf("failed to get task with id %d: %w", cmd.id, err)
	}

	task := &models.Task{
		ID:          cmd.id,
		Description: cmd.description,
		EmployeeID:  cmd.employeeID,
		ProjectID:   cmd.projectID,
		IsCompleted: cmd.isCompleted,
	}

	if err := uc.repo.UpdateTask(task); err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	uc.logMessage(ctx, fmt.Sprintf("Updating task (id: %d)", cmd.id))
	return nil
}

// Команда для удаления задачи
type DeleteTaskCommand struct {
	id uint32
}

func NewDeleteTaskCommand(id uint32) *DeleteTaskCommand {
	return &DeleteTaskCommand{id: id}
}

func (uc *ProjectUseCases) DeleteTask(ctx context.Context, cmd *DeleteTaskCommand) error {
	claims := ctx.Value(common.ContextKeyClaims).(*common.Claims)

	if claims.Role != adminRole {
		return common.ErrForbidden
	}

	if err := uc.repo.DeleteTask(cmd.id); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	uc.logMessage(ctx, fmt.Sprintf("Deleting task (id: %d)", cmd.id))
	return nil
}

// Запрос для получения задач сотрудника
type GetTasksByEmployeeIDQuery struct {
	employeeID uint32
}

func NewGetTasksByEmployeeIDQuery(employeeID uint32) *GetTasksByEmployeeIDQuery {
	return &GetTasksByEmployeeIDQuery{employeeID: employeeID}
}

func (uc *ProjectUseCases) GetTasksByEmployeeID(ctx context.Context, query *GetTasksByEmployeeIDQuery) ([]*models.Task, error) {
	tasks, err := uc.repo.GetTasksByEmployeeID(query.employeeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks for employee id %d: %w", query.employeeID, err)
	}

	uc.logMessage(ctx, fmt.Sprintf("Fetching tasks by employee id (id: %d)", query.employeeID))
	return tasks, nil
}

type TaskFilter struct {
	EmployeeID  uint32
	ProjectID   uint32
	IsCompleted *bool
}

func (uc *ProjectUseCases) GetTasks(ctx context.Context, filter TaskFilter) ([]*models.Task, error) {
	uc.logMessage(ctx, "Fetching tasks")
	return uc.repo.GetTasks(filter)
}
