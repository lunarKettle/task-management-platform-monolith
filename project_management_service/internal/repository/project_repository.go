package repository

import (
	"fmt"
	"project_management_service/internal/models"
)

type ProjectRepository struct {
	db *Database
}

func NewProjectRepository(database *Database) ProjectRepository {
	return ProjectRepository{db: database}
}

func (r *ProjectRepository) AddProject(project models.Project) (uint32, error) {
	query := `INSERT INTO projects (name, description, start_date, planned_end_date, actual_end_date, status, priority, team_id, budget)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	var id uint32
	err := r.db.connection.QueryRow(query, project.Name, project.Description, project.StartDate, project.PlannedEndDate, project.ActualEndDate, project.Status, project.Priority, project.TeamId, project.Budget).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error inserting project: %v", err)
	}
	return id, nil
}

func (r *ProjectRepository) GetProjectById(projectId uint32) (models.Project, error) {
	query := "SELECT * FROM projects"
	var project models.Project
	//err := r.db.connection.Get(&project, query)
	err := r.db.connection.QueryRow(query, projectId).Scan(&project.Id, &project.Name,
		&project.StartDate, &project.PlannedEndDate, &project.ActualEndDate,
		&project.Status, &project.Priority, &project.TeamId, &project.Budget)
	return project, err
}
