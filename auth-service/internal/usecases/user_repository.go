package usecases

import (
	"github.com/lunarKettle/task-management-platform/auth-service/internal/models"
)

type UserRepository interface {
	Create(user *models.User) (uint32, error)
	GetById(id uint32) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	DeleteById(id uint32) error
	Update(user *models.User) error
}
