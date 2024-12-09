package infrastructure

import (
	"database/sql"
	"fmt"

	"github.com/lunarKettle/task-management-platform-monolith/internal/aggregate/project/models"
	"github.com/lunarKettle/task-management-platform-monolith/pkg/common"
)

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(database *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: database}
}

func (r *ProjectRepository) CreateProject(project *models.Project) (uint32, error) {
	query := `INSERT INTO projects (name, description, start_date, planned_end_date, actual_end_date, status, priority, team_id, budget)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	var id uint32
	err := r.db.QueryRow(query,
		project.Name,
		project.Description,
		project.StartDate,
		project.PlannedEndDate,
		project.ActualEndDate,
		project.Status,
		project.Priority,
		project.Team.ID,
		project.Budget).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error inserting project: %v", err)
	}
	return id, nil
}

func (r *ProjectRepository) UpdateProject(project *models.Project) error {
	query := `UPDATE projects
		SET name = $1, description = $2, start_date = $3, planned_end_date = $4, actual_end_date = $5,
		    status = $6, priority = $7, team_id = $8, budget = $9
		WHERE id = $10`

	_, err := r.db.Exec(query,
		project.Name,
		project.Description,
		project.StartDate,
		project.PlannedEndDate,
		project.ActualEndDate,
		project.Status,
		project.Priority,
		project.Team.ID,
		project.Budget,
		project.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("project with id %d not found: %w", project.Id, common.ErrNotFound)
		}
		return fmt.Errorf("error updating project: %v", err)
	}
	return nil
}

func (r *ProjectRepository) DeleteProject(projectId uint32) error {
	query := `DELETE FROM projects WHERE id = $1`

	_, err := r.db.Exec(query, projectId)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("project with id %d not found: %w", projectId, common.ErrNotFound)
		}
		return fmt.Errorf("error deleting project: %v", err)
	}
	return nil
}

func (r *ProjectRepository) GetAllProjects() ([]*models.Project, error) {
	query := `
	SELECT 
    	p.id, 
		p.name, 
		p.description, 
		p.start_date, 
		p.planned_end_date, 
		p.actual_end_date, 
		p.status, 
		p.priority,  
    	p.team_id, 
		p.budget
	FROM 
    	projects p;`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query projects: %w", err)
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		project := &models.Project{}
		var teamID sql.NullInt64

		err := rows.Scan(
			&project.Id,
			&project.Name,
			&project.Description,
			&project.StartDate,
			&project.PlannedEndDate,
			&project.ActualEndDate,
			&project.Status,
			&project.Priority,
			&teamID,
			&project.Budget,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project row: %w", err)
		}

		if teamID.Valid {
			team, err := r.GetTeamById(uint32(teamID.Int64))
			if err != nil {
				return nil, fmt.Errorf("failed to get team for project %d: %w", project.Id, err)
			}
			project.Team = team
		} else {
			project.Team = nil
		}

		projects = append(projects, project)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during row iteration: %w", err)
	}

	return projects, nil
}

func (r *ProjectRepository) GetProjectById(projectId uint32) (*models.Project, error) {
	query := `
    SELECT 
        p.id, 
        p.name, 
        p.description, 
        p.start_date, 
        p.planned_end_date, 
        p.actual_end_date, 
        p.status, 
        p.priority,  
        p.team_id, 
        t.name AS team_name,
        t.manager_id,
        p.budget
    FROM 
        projects p
    LEFT JOIN
        teams t
    ON
        p.team_id = t.id
    WHERE 
        p.id = $1;`

	project := &models.Project{}
	var teamID sql.NullInt64
	var teamName sql.NullString
	var managerID sql.NullInt64

	row := r.db.QueryRow(query, projectId)

	err := row.Scan(
		&project.Id,
		&project.Name,
		&project.Description,
		&project.StartDate,
		&project.PlannedEndDate,
		&project.ActualEndDate,
		&project.Status,
		&project.Priority,
		&teamID,
		&teamName,
		&managerID,
		&project.Budget)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("project with id %d not found: %w", projectId, common.ErrNotFound)
		}
		return nil, err
	}

	if teamID.Valid {
		project.Team = &models.Team{
			ID:        uint32(teamID.Int64),
			Name:      teamName.String,
			ManagerID: uint32(managerID.Int64),
		}
	} else {
		project.Team = nil
	}

	return project, nil
}

func (r *ProjectRepository) CreateTeam(team *models.Team) (uint32, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	checkQuery := "SELECT team_id FROM users WHERE id = $1 FOR UPDATE"
	var existingValue sql.NullString
	err = tx.QueryRow(checkQuery, team.ManagerID).Scan(&existingValue)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("record with id %d not found: %w", team.ManagerID, common.ErrNotFound)
		}
		return 0, fmt.Errorf("failed to check column: %w", err)
	}

	if existingValue.Valid {
		return 0, fmt.Errorf("team_id is not NULL for id %d", team.ManagerID)
	}

	insertTeamQuery := `INSERT INTO teams (name, manager_id) VALUES ($1, $2) RETURNING id`
	var teamID uint32
	err = tx.QueryRow(insertTeamQuery, team.Name, team.ManagerID).Scan(&teamID)
	if err != nil {
		return 0, fmt.Errorf("error inserting team: %w", err)
	}

	updateQuery := "UPDATE users SET team_id = $1 WHERE id = $2"
	_, err = tx.Exec(updateQuery, teamID, team.ManagerID)
	if err != nil {
		return 0, fmt.Errorf("failed to set team_id: %w", err)
	}

	return teamID, nil
}

func (r *ProjectRepository) UpdateTeam(team *models.Team) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	checkQuery := "SELECT team_id FROM users WHERE id = $1 FOR UPDATE"
	var existingTeamID sql.NullInt32
	err = tx.QueryRow(checkQuery, team.ManagerID).Scan(&existingTeamID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("record with id %d not found: %w", team.ManagerID, common.ErrNotFound)
		}
		return fmt.Errorf("failed to check team_id: %w", err)
	}

	if !existingTeamID.Valid || uint32(existingTeamID.Int32) != team.ID {
		return fmt.Errorf("manager with id %d is not associated with team %d", team.ManagerID, team.ID)
	}

	updateTeamQuery := `
		UPDATE teams
		SET name = $1, manager_id = $2
		WHERE id = $3
	`
	_, err = tx.Exec(updateTeamQuery, team.Name, team.ManagerID, team.ID)
	if err != nil {
		return fmt.Errorf("failed to update team: %w", err)
	}

	updateUserQuery := `
		UPDATE users
		SET team_id = $1
		WHERE id = $2
	`
	_, err = tx.Exec(updateUserQuery, team.ID, team.ManagerID)
	if err != nil {
		return fmt.Errorf("failed to update user team_id: %w", err)
	}

	return nil
}

func (r *ProjectRepository) DeleteTeam(teamId uint32) error {
	query := `DELETE FROM teams WHERE id = $1`

	_, err := r.db.Exec(query, teamId)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("team with id %d not found: %w", teamId, common.ErrNotFound)
		}
		return fmt.Errorf("error deleting team: %v", err)
	}
	return nil
}

func (r *ProjectRepository) GetAllTeams() ([]*models.Team, error) {
	query := `
	SELECT 
		t.id, 
		t.name, 
		t.manager_id, 
		u.id AS member_id,
		u.username, 
		u.role 
	FROM 
		teams t
	LEFT JOIN 
		users u 
	ON 
		u.team_id = t.id
	ORDER BY 
		t.id, u.id`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query teams: %w", err)
	}
	defer rows.Close()

	var teams []*models.Team
	var currentTeam *models.Team
	var lastTeamID uint32

	for rows.Next() {
		var (
			teamID     uint32
			teamName   string
			managerID  sql.NullInt64
			memberID   sql.NullInt64
			memberName string
			role       sql.NullString
		)

		err = rows.Scan(&teamID, &teamName, &managerID, &memberID, &memberName, &role)
		if err != nil {
			return nil, fmt.Errorf("failed to scan team row: %w", err)
		}

		// Если мы перешли к новой команде
		if currentTeam == nil || teamID != lastTeamID {
			if currentTeam != nil {
				teams = append(teams, currentTeam)
			}
			currentTeam = &models.Team{
				ID:      teamID,
				Name:    teamName,
				Members: []models.Member{},
			}
			if managerID.Valid {
				currentTeam.ManagerID = uint32(managerID.Int64)
			}
			lastTeamID = teamID
		}

		// Добавляем участника, если он есть
		if memberID.Valid {
			currentTeam.Members = append(currentTeam.Members, models.Member{
				ID:   uint32(memberID.Int64),
				Name: memberName,
				Role: role.String,
			})
		}
	}

	// Добавляем последнюю команду, если есть
	if currentTeam != nil {
		teams = append(teams, currentTeam)
	}

	// Проверяем ошибки при итерации
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over team rows: %w", err)
	}

	return teams, nil
}

func (r *ProjectRepository) GetTeamById(teamId uint32) (*models.Team, error) {
	query := `
	SELECT 
		id, name, manager_id
	FROM 
		teams
	WHERE id=$1`

	team := &models.Team{}

	var managerID sql.NullInt64

	err := r.db.QueryRow(query, teamId).Scan(&team.ID, &team.Name, &managerID)

	if managerID.Valid {
		team.ManagerID = uint32(managerID.Int64)
	} else {
		team.ManagerID = 0
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("team with id %d not found: %w", teamId, common.ErrNotFound)
		}
		return nil, err
	}
	return team, nil
}

func (r *ProjectRepository) GetTeamIdByUserID(userID uint32) (uint32, error) {
	var teamID sql.NullInt32
	query := "SELECT team_id FROM users WHERE id = $1"

	err := r.db.QueryRow(query, userID).Scan(&teamID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no user found with id %d: %w", userID, common.ErrNotFound)
		}
		return 0, err
	}

	if !teamID.Valid {
		return 0, fmt.Errorf("team_id is NULL for user id %d: %w", userID, common.ErrNotFound)
	}

	return uint32(teamID.Int32), nil
}

func (r *ProjectRepository) GetAllMembers() ([]*models.Member, error) {
	query := `
	SELECT 
		id, 
		username, 
		role, 
		team_id 
	FROM 
		users`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query members: %w", err)
	}
	defer rows.Close()

	var members []*models.Member

	for rows.Next() {
		var (
			memberID uint32
			name     string
			role     sql.NullString
			teamID   sql.NullInt64
		)

		err = rows.Scan(&memberID, &name, &role, &teamID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan member row: %w", err)
		}

		member := &models.Member{
			ID:   memberID,
			Name: name,
			Role: role.String, // Если `role` NULL, будет пустая строка
		}

		if teamID.Valid {
			member.TeamID = uint32(teamID.Int64)
		} else {
			member.TeamID = 0 // Если `team_id` NULL, используем 0 как значение по умолчанию
		}

		members = append(members, member)
	}

	// Проверяем ошибки, возникшие при итерации
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over member rows: %w", err)
	}

	return members, nil
}

func (r *ProjectRepository) CreateTask(task *models.Task) (uint32, error) {
	query := `INSERT INTO tasks (description, employee_id, project_id, is_completed)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	var id uint32
	err := r.db.QueryRow(query,
		task.Description,
		task.EmployeeID,
		task.ProjectID,
		task.IsCompleted).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error inserting task: %v", err)
	}
	return id, nil
}

func (r *ProjectRepository) UpdateTask(task *models.Task) error {
	query := `UPDATE tasks
		SET description = $1, employee_id = $2, project_id = $3, is_completed = $4
		WHERE id = $5`

	_, err := r.db.Exec(query,
		task.Description,
		task.EmployeeID,
		task.ProjectID,
		task.IsCompleted,
		task.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("task with id %d not found: %w", task.ID, common.ErrNotFound)
		}
		return fmt.Errorf("error updating task: %v", err)
	}
	return nil
}

func (r *ProjectRepository) DeleteTask(taskID uint32) error {
	query := `DELETE FROM tasks WHERE id = $1`

	_, err := r.db.Exec(query, taskID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("task with id %d not found: %w", taskID, common.ErrNotFound)
		}
		return fmt.Errorf("error deleting task: %v", err)
	}
	return nil
}

func (r *ProjectRepository) GetTaskById(taskID uint32) (*models.Task, error) {
	query := `
	SELECT 
    	t.id, 
		t.description, 
		t.employee_id, 
		t.project_id,
		t.is_completed
	FROM 
    	tasks t
	WHERE 
    	t.id = $1;`

	task := &models.Task{}

	row := r.db.QueryRow(query, taskID)

	err := row.Scan(
		&task.ID,
		&task.Description,
		&task.EmployeeID,
		&task.ProjectID,
		&task.IsCompleted,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task with id %d not found: %w", taskID, common.ErrNotFound)
		}
		return nil, err
	}

	return task, nil
}

func (r *ProjectRepository) GetTasksByEmployeeID(employeeID uint32) ([]*models.Task, error) {
	query := `
	SELECT 
    	t.id, 
		t.description, 
		t.employee_id, 
		t.project_id,
		t.is_completed
	FROM 
    	tasks t
	WHERE 
    	t.employee_id = $1;`

	rows, err := r.db.Query(query, employeeID)
	if err != nil {
		return nil, fmt.Errorf("error querying tasks for employee_id %d: %v", employeeID, err)
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		task := &models.Task{}
		err := rows.Scan(
			&task.ID,
			&task.Description,
			&task.EmployeeID,
			&task.ProjectID,
			&task.IsCompleted,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning task row: %v", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return tasks, nil
}
