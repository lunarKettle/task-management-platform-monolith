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
