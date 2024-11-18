package infrastructure

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lunarKettle/task-management-platform-monolith/pkg/common"
)

type JWTManager struct {
	secretKey string
}

func NewJWTManager(secretKey string) *JWTManager {
	return &JWTManager{
		secretKey: secretKey,
	}
}

type Claims struct {
	UserID uint32 `json:"sub"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (m *JWTManager) GenerateToken(userID uint32, role string) (string, error) {
	expirationTime := time.Now().Add(time.Hour * 1)

	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(m.secretKey))

	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return tokenString, nil
}

func (m *JWTManager) ValidateToken(token string) (uint32, string, error) {
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, common.ErrUnexpectedSigningMethod
		}
		return []byte(m.secretKey), nil
	})

	if err != nil {
		return 0, "", common.ErrInvalidToken
	}

	if !tkn.Valid {
		return 0, "", common.ErrTokenNotValid
	}

	return claims.UserID, claims.Role, nil
}
