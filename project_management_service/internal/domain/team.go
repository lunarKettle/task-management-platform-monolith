package domain

import "time"

type Team struct {
	ID        uint32
	Name      string
	Members   []Member
	CreatedAt time.Time
	UpdatedAt time.Time
}
