package infrastructure

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTGenerator struct {
	secretKey string
}

func NewJWTGenerator(secretKey string) *JWTGenerator {
	return &JWTGenerator{
		secretKey: secretKey,
	}
}

func (g *JWTGenerator) GenerateToken(userID uint32, role string) (string, error) {
	expirationTime := time.Now().Add(time.Hour * 1)

	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(g.secretKey))

	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return tokenString, nil
}
