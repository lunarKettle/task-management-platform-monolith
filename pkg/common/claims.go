package common

import "github.com/golang-jwt/jwt/v5"

const (
	ContextKeyClaims = "userClaims"
)

type Claims struct {
	UserID uint32 `json:"sub"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}
