package infrastructure

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lunarKettle/task-management-platform-monolith/internal/aggregate/project/models"
	"github.com/lunarKettle/task-management-platform-monolith/internal/aggregate/project/usecases"
	"github.com/stretchr/testify/assert"
)

func TestCreateProject(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewProjectRepository(db)

	project := &models.Project{
		Name:           "Project Test",
		Description:    "A test project",
		StartDate:      time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
		PlannedEndDate: time.Date(2024, time.December, 31, 0, 0, 0, 0, time.UTC),
		ActualEndDate:  time.Time{},
		Status:         "active",
		Priority:       1,
		Team: &models.Team{
			ID: 1,
		},
		Budget: 1000,
	}

	mock.ExpectQuery(`INSERT INTO projects`).
		WithArgs(
			project.Name,
			project.Description,
			project.StartDate,
			project.PlannedEndDate,
			project.ActualEndDate,
			project.Status,
			project.Priority,
			project.Team.ID,
			project.Budget,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	id, err := repo.CreateProject(project)
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), id)
}

func TestUpdateProject(t *testing.T) {
	// Подготовка
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewProjectRepository(db)

	project := &models.Project{
		Id:             1,
		Name:           "Updated Project",
		Description:    "Updated description",
		StartDate:      time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
		PlannedEndDate: time.Date(2024, time.December, 31, 0, 0, 0, 0, time.UTC),
		ActualEndDate:  time.Date(2024, time.February, 31, 0, 0, 0, 0, time.UTC),
		Status:         "active",
		Priority:       1,
		Team: &models.Team{
			ID: 1,
		},
		Budget: 1500,
	}

	// Mock запросы
	mock.ExpectExec(`UPDATE projects`).
		WithArgs(
			project.Name,
			project.Description,
			project.StartDate,
			project.PlannedEndDate,
			project.ActualEndDate,
			project.Status,
			project.Priority,
			project.Team.ID,
			project.Budget,
			project.Id,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Выполняем тест
	err = repo.UpdateProject(project)
	assert.NoError(t, err)
}

func TestDeleteProject(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewProjectRepository(db)

	projectID := uint32(1)

	mock.ExpectExec(`DELETE FROM projects WHERE id`).
		WithArgs(projectID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.DeleteProject(projectID)
	assert.NoError(t, err)
}

func TestGetAllProjects(t *testing.T) {
	// Подготовка
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &ProjectRepository{db: db}

	// Данные для проектов
	projectRows := sqlmock.NewRows([]string{
		"id", "name", "description", "start_date", "planned_end_date", "actual_end_date", "status", "priority", "team_id", "budget",
	}).
		AddRow(1, "Project 1", "Description 1", time.Now(), time.Now().AddDate(0, 1, 0), time.Time{}, "Active", 1, 1, 1000.0).
		AddRow(2, "Project 2", "Description 2", time.Now(), time.Now().AddDate(0, 2, 0), time.Now().AddDate(0, 1, 15), "Completed", 2, 2, 2000.0)

	// Данные для команды с id=1
	team1Rows := sqlmock.NewRows([]string{"id", "name", "manager_id"}).
		AddRow(1, "Team 1", 1)

	// Данные для команды с id=2
	team2Rows := sqlmock.NewRows([]string{"id", "name", "manager_id"}).
		AddRow(2, "Team 2", 2)

	// Mock запросы
	mock.ExpectQuery(`SELECT p.id, p.name, p.description, p.start_date, p.planned_end_date, p.actual_end_date, p.status, p.priority, p.team_id, p.budget FROM projects p;`).
		WillReturnRows(projectRows)

	mock.ExpectQuery(`SELECT id, name, manager_id FROM teams WHERE id=\$1`).
		WithArgs(1).
		WillReturnRows(team1Rows)

	mock.ExpectQuery(`SELECT id, name, manager_id FROM teams WHERE id=\$1`).
		WithArgs(2).
		WillReturnRows(team2Rows)

	// Выполняем тест
	projects, err := repo.GetAllProjects()
	assert.NoError(t, err)
	assert.NotNil(t, projects)
	assert.Equal(t, 2, len(projects))

	// Проверяем детали проектов
	project1 := projects[0]
	assert.Equal(t, uint32(1), project1.Id)
	assert.Equal(t, "Project 1", project1.Name)
	assert.Equal(t, "Description 1", project1.Description)
	assert.Equal(t, "Active", project1.Status)
	assert.Equal(t, uint32(1), project1.Priority)
	assert.Equal(t, float64(1000.0), project1.Budget)

	// Проверка команды проекта 1
	assert.NotNil(t, project1.Team)
	assert.Equal(t, uint32(1), project1.Team.ID)
	assert.Equal(t, "Team 1", project1.Team.Name)
	assert.Equal(t, uint32(1), project1.Team.ManagerID)

	project2 := projects[1]
	assert.Equal(t, uint32(2), project2.Id)
	assert.Equal(t, "Project 2", project2.Name)
	assert.Equal(t, "Description 2", project2.Description)
	assert.Equal(t, "Completed", project2.Status)
	assert.Equal(t, uint32(2), project2.Priority)
	assert.Equal(t, float64(2000.0), project2.Budget)

	// Проверка команды проекта 2
	assert.NotNil(t, project2.Team)
	assert.Equal(t, uint32(2), project2.Team.ID)
	assert.Equal(t, "Team 2", project2.Team.Name)
	assert.Equal(t, uint32(2), project2.Team.ManagerID)

	// Проверка выполнения всех ожидаемых запросов
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetProjectById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &ProjectRepository{db: db}

	projectID := uint32(1)
	projectName := "Project 1"
	description := "Description 1"
	startDate := time.Now()
	plannedEndDate := time.Now().AddDate(0, 1, 0)
	actualEndDate := time.Time{}
	status := "Active"
	priority := uint32(1)
	teamID := uint32(2)
	teamName := "Team 1"
	managerID := uint32(3)
	budget := float64(1000.0)

	projectRow := sqlmock.NewRows([]string{
		"id", "name", "description", "start_date", "planned_end_date", "actual_end_date", "status", "priority", "team_id", "team_name", "manager_id", "budget",
	}).AddRow(
		projectID, projectName, description, startDate, plannedEndDate, actualEndDate, status, priority, teamID, teamName, managerID, budget,
	)

	mock.ExpectQuery(`SELECT p.id, p.name, p.description, p.start_date, p.planned_end_date, p.actual_end_date, p.status, p.priority, p.team_id, t.name AS team_name, t.manager_id, p.budget FROM projects p LEFT JOIN teams t ON p.team_id = t.id WHERE p.id = \$1;`).
		WithArgs(projectID).
		WillReturnRows(projectRow)

	project, err := repo.GetProjectById(projectID)
	assert.NoError(t, err)
	assert.NotNil(t, project)

	assert.Equal(t, projectID, project.Id)
	assert.Equal(t, projectName, project.Name)
	assert.Equal(t, description, project.Description)
	assert.Equal(t, startDate, project.StartDate)
	assert.Equal(t, plannedEndDate, project.PlannedEndDate)
	assert.Equal(t, actualEndDate, project.ActualEndDate)
	assert.Equal(t, status, project.Status)
	assert.Equal(t, priority, project.Priority)
	assert.Equal(t, budget, project.Budget)

	assert.NotNil(t, project.Team)
	assert.Equal(t, teamID, project.Team.ID)
	assert.Equal(t, teamName, project.Team.Name)
	assert.Equal(t, managerID, project.Team.ManagerID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateTeam(t *testing.T) {
	// TODO
}
func TestUpdateTeam(t *testing.T) {
	// TODO
}
func TestDeleteTeam(t *testing.T) {
	// TODO
}

func TestGetAllTeams(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewProjectRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "manager_id", "member_id", "username", "role"}).
		AddRow(1, "Team 1", 1, 1, "User 1", "developer").
		AddRow(1, "Team 1", 1, 2, "User 2", "designer").
		AddRow(2, "Team 2", 2, 3, "User 3", "manager")

	mock.ExpectQuery(`SELECT .* FROM teams`).
		WillReturnRows(rows)

	teams, err := repo.GetAllTeams()
	assert.NoError(t, err)
	assert.Len(t, teams, 2)
	assert.Equal(t, uint32(1), teams[0].ID)
	assert.Equal(t, uint32(2), teams[1].ID)
}

func TestGetTeamById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewProjectRepository(db)

	teamID := uint32(1)

	rows := sqlmock.NewRows([]string{"id", "name", "manager_id"}).
		AddRow(1, "Team 1", 1)

	mock.ExpectQuery(`SELECT .* FROM teams`).
		WithArgs(teamID).
		WillReturnRows(rows)

	team, err := repo.GetTeamById(teamID)
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), team.ID)
}

func TestGetTeamIdByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewProjectRepository(db)

	userID := uint32(1)

	rows := sqlmock.NewRows([]string{"team_id"}).
		AddRow(1)

	mock.ExpectQuery(`SELECT team_id FROM users`).
		WithArgs(userID).
		WillReturnRows(rows)

	teamID, err := repo.GetTeamIdByUserID(userID)
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), teamID)
}

func TestGetMember(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewProjectRepository(db)

	userID := uint32(1)

	rows := sqlmock.NewRows([]string{"id", "username", "role", "team_id"}).
		AddRow(1, "User 1", "developer", 1)

	mock.ExpectQuery(`SELECT .* FROM users`).
		WithArgs(userID).
		WillReturnRows(rows)

	member, err := repo.GetMember(userID)
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), member.ID)
}

func TestGetMembers(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewProjectRepository(db)

	filter := usecases.MemberFilter{
		Role:   "developer",
		TeamID: 1,
	}

	rows := sqlmock.NewRows([]string{"id", "username", "role", "team_id"}).
		AddRow(1, "User 1", "developer", 1).
		AddRow(2, "User 2", "developer", 1)

	mock.ExpectQuery(`SELECT .* FROM users`).
		WithArgs(filter.Role, filter.TeamID).
		WillReturnRows(rows)

	members, err := repo.GetMembers(filter)
	assert.NoError(t, err)
	assert.Len(t, members, 2)
	assert.Equal(t, uint32(1), members[0].ID)
	assert.Equal(t, uint32(2), members[1].ID)
}

func TestGetTasksByEmployeeID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewProjectRepository(db)

	employeeID := uint32(1)

	rows := sqlmock.NewRows([]string{"id", "title", "description", "status", "assigned_to", "deadline"}).
		AddRow(1, "Task 1", "Description 1", "active", 1, "2024-01-01").
		AddRow(2, "Task 2", "Description 2", "completed", 1, "2024-01-02")

	mock.ExpectQuery(`SELECT .* FROM tasks`).
		WithArgs(employeeID).
		WillReturnRows(rows)

	tasks, err := repo.GetTasksByEmployeeID(employeeID)
	assert.NoError(t, err)
	assert.Len(t, tasks, 2)
	assert.Equal(t, uint32(1), tasks[0].ID)
	assert.Equal(t, uint32(1), tasks[1].ID)
}
