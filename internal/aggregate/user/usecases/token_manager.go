package usecases

import "github.com/lunarKettle/task-management-platform-monolith/pkg/common"

type TokenManager interface {
	GenerateToken(userID uint32, role string) (string, error)
	ValidateToken(token string) (*common.Claims, error)
}
