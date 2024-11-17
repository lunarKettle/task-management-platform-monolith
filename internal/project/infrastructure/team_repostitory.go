package infrastructure

import (
	"database/sql"
	"fmt"

	"github.com/lunarKettle/task-management-platform-monolith/internal/project/models"
)

type TeamRepository struct {
	db *sql.DB
}

func NewTeamRepository(database *sql.DB) *TeamRepository {
	return &TeamRepository{db: database}
}

func (r *TeamRepository) Create(team models.Team) (uint32, error) {
	query := `INSERT INTO teams (name, created_at, updated_at)
		VALUES ($1, $2, $3)
		RETURNING id`

	var id uint32
	err := r.db.QueryRow(query, team.Name, team.CreatedAt, team.UpdatedAt).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error inserting project: %v", err)
	}
	return id, nil
}

func (r *TeamRepository) Update(team models.Team) error {
	query := `UPDATE teams
		SET name = $1, created_at = $2, updated_at = $3
		WHERE id = $4`

	_, err := r.db.Exec(query, team.Name, team.CreatedAt, team.UpdatedAt, team.ID)
	if err != nil {
		return fmt.Errorf("error updating project: %v", err)
	}
	return nil
}

func (r *TeamRepository) Delete(teamId uint32) error {
	query := `DELETE FROM teams WHERE id = $1`

	_, err := r.db.Exec(query, teamId)
	if err != nil {
		return fmt.Errorf("error deleting project: %v", err)
	}
	return nil
}

func (r *TeamRepository) GetById(teamId uint32) (models.Team, error) {
	query := "SELECT * FROM teams"
	var team models.Team
	err := r.db.QueryRow(query, teamId).Scan(
		&team.ID, &team.Name, &team.CreatedAt, &team.UpdatedAt)
	return team, err
}
