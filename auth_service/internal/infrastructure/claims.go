package infrastructure

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID uint32 `json:"sub"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}
