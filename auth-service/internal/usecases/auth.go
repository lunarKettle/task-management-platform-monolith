package usecases

import (
	"errors"
	"fmt"

	"github.com/lunarKettle/task-management-platform/auth-service/internal/common"
	"github.com/lunarKettle/task-management-platform/auth-service/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type AuthUseCases struct {
	repo           UserRepository
	tokenGenerator TokenGenerator
}

func NewAuthUseCases(repo UserRepository, tokenGenerator TokenGenerator) *AuthUseCases {
	return &AuthUseCases{
		repo:           repo,
		tokenGenerator: tokenGenerator,
	}
}

type CreateUserCommand struct {
	username string
	email    string
	password string
	role     string
}

func NewCreateUserCommand(
	username string,
	email string,
	password string,
	role string,
) *CreateUserCommand {
	return &CreateUserCommand{
		username: username,
		email:    email,
		password: password,
		role:     role,
	}
}

func (a *AuthUseCases) CreateUser(cmd *CreateUserCommand) (string, error) {
	_, err := a.repo.GetByUsername(cmd.username)

	switch {
	case errors.Is(err, common.ErrNotFound):
	case err == nil:
		return "", fmt.Errorf("%w: user with username %q already exists", common.ErrAlreadyExists, cmd.username)
	default:
		return "", fmt.Errorf("failed to get user by username %q: %w", cmd.username, err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(cmd.password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	user := models.User{
		Username:     cmd.username,
		Email:        cmd.email,
		PasswordHash: passwordHash,
		Role:         cmd.role,
	}

	userId, err := a.repo.Create(user)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	token, err := a.tokenGenerator.GenerateToken(userId, cmd.role)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}
