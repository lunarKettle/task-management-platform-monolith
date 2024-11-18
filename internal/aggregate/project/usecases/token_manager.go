package usecases

type TokenManager interface {
	GenerateToken(userID uint32, role string) (string, error)
	ValidateToken(token string) (uint32, string, error)
}
