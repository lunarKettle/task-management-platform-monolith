package usecases

type TokenGenerator interface {
	GenerateToken(userID uint32, role string) (string, error)
}
