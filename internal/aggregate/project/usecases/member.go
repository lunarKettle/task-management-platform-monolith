package usecases

import "github.com/lunarKettle/task-management-platform-monolith/internal/aggregate/project/models"

type Member struct {
	id   uint32
	role string
}

func NewMember(id uint32, role string) *Member {
	return &Member{
		id:   id,
		role: role,
	}
}

func (m *Member) ToModel() models.Member {
	return models.Member{
		ID:   m.id,
		Role: m.role,
	}
}
