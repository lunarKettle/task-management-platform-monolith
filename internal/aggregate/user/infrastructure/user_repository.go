package infrastructure

import (
	"database/sql"

	"github.com/lunarKettle/task-management-platform-monolith/internal/aggregate/user/models"
	"github.com/lunarKettle/task-management-platform-monolith/pkg/common"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(database *sql.DB) *UserRepository {
	return &UserRepository{db: database}
}

// Create создаёт нового пользователя
func (r *UserRepository) Create(user *models.User) (uint32, error) {
	query := `INSERT INTO users (username, email, password_hash, role) 
              VALUES ($1, $2, $3, $4) RETURNING id`

	var userID uint32
	err := r.db.QueryRow(query, user.Username, user.Email, user.PasswordHash, user.Role).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

// GetById получает пользователя по его ID.
func (r *UserRepository) GetById(id uint32) (*models.User, error) {
	query := `SELECT id, username, email, password_hash, role 
              FROM users WHERE id = $1`

	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, common.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

// GetByUsername получает пользователя по его Username.
func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	query := `SELECT id, username, email, password_hash, role 
              FROM users WHERE username = $1`

	user := &models.User{}
	err := r.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, common.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

// DeleteById удаляет пользователя по его ID.
func (r *UserRepository) DeleteById(id uint32) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return common.ErrNotFound
	}

	return nil
}

// Update обновляет данные пользователя.
func (r *UserRepository) Update(user *models.User) error {
	query := `UPDATE users 
              SET username = $1, email = $2, password_hash = $3, role = $4, 
              WHERE id = $5`

	result, err := r.db.Exec(query, user.Username, user.Email, user.PasswordHash, user.Role, user.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return common.ErrNotFound
	}

	return nil
}
