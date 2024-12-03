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

func (m *JWTManager) GenerateToken(userID uint32, role string) (string, error) {
	expirationTime := time.Now().Add(time.Hour * 1)

	claims := &common.Claims{
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

func (m *JWTManager) ValidateToken(token string) (*common.Claims, error) {
	claims := &common.Claims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, common.ErrUnexpectedSigningMethod
		}
		return []byte(m.secretKey), nil
	})

	if err != nil {
		return nil, common.ErrInvalidToken
	}

	if !tkn.Valid {
		return nil, common.ErrTokenNotValid
	}

	return claims, nil
}
