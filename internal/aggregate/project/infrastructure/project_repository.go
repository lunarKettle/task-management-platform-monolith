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

func (r *ProjectRepository) GetProjectById(projectId uint32) (*models.Project, error) {
	query := `
	SELECT 
    	p.id, 
		p.name, 
		p.description, 
		p.start_date, 
		p.planned_end_date, 
		p.actual_end_date, 
		p.status, p.priority,  
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
		&project.Team.ID,
		&project.Team.Name,
		&project.Team.ManagerID,
		&project.Budget)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("project with id %d not found: %w", projectId, common.ErrNotFound)
		}
		return nil, err
	}

	return project, nil
}

func (r *ProjectRepository) CreateTeam(team *models.Team) (uint32, error) {
	query := `INSERT INTO teams (name, manager_id)
		VALUES ($1, $2)
		RETURNING id`

	var id uint32
	err := r.db.QueryRow(query, team.Name, team.ManagerID).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error inserting team: %v", err)
	}
	return id, nil
}

func (r *ProjectRepository) UpdateTeam(team *models.Team) error {
	query := `UPDATE teams
		SET name = $1, manager_id = $2
		WHERE id = $3`

	_, err := r.db.Exec(query, team.Name, team.ManagerID, team.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("team with id %d not found: %w", team.ID, common.ErrNotFound)
		}
		return fmt.Errorf("error updating project: %v", err)
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
	var teamID uint32
	query := "SELECT team_id FROM users WHERE id = $1"

	err := r.db.QueryRow(query, userID).Scan(&teamID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no user found with id %d: %w", userID, common.ErrNotFound)
		}
		return 0, err
	}

	return teamID, nil
}

func (r *ProjectRepository) CreateTask(task *models.Task) (uint32, error) {
	query := `INSERT INTO tasks (description, employee_id, project_id)
		VALUES ($1, $2, $3)
		RETURNING id`

	var id uint32
	err := r.db.QueryRow(query,
		task.Description,
		task.EmployeeID,
		task.ProjectID).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error inserting task: %v", err)
	}
	return id, nil
}

func (r *ProjectRepository) UpdateTask(task *models.Task) error {
	query := `UPDATE tasks
		SET description = $1, employee_id = $2, project_id = $3
		WHERE id = $4`

	_, err := r.db.Exec(query,
		task.Description,
		task.EmployeeID,
		task.ProjectID,
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
		t.project_id
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
		&task.ProjectID)

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
		t.project_id
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
