package infrastructure

import (
	"database/sql"
	"fmt"

	"github.com/lunarKettle/task-management-platform-monolith/internal/project/models"
	"github.com/lunarKettle/task-management-platform-monolith/pkg/common"
)

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(database *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: database}
}

func (r *ProjectRepository) Create(project *models.Project) (uint32, error) {
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
		project.TeamId,
		project.Budget).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error inserting project: %v", err)
	}
	return id, nil
}

func (r *ProjectRepository) Update(project *models.Project) error {
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
		project.TeamId,
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

func (r *ProjectRepository) Delete(projectId uint32) error {
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

func (r *ProjectRepository) GetById(projectId uint32) (*models.Project, error) {
	query := "SELECT * FROM projects WHERE id = $1"
	project := &models.Project{}

	err := r.db.QueryRow(query, projectId).Scan(
		&project.Id,
		&project.Name,
		&project.Description,
		&project.StartDate,
		&project.PlannedEndDate,
		&project.ActualEndDate,
		&project.Status,
		&project.Priority,
		&project.TeamId,
		&project.Budget)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("project with id %d not found: %w", projectId, common.ErrNotFound)
		}
		return nil, err
	}

	return project, nil
}
