package models

type User struct {
	ID           uint32
	Username     string
	Email        string
	PasswordHash []byte
	Role         string
}
