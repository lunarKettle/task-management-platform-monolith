package usecases

import "github.com/lunarKettle/task-management-platform-monolith/internal/aggregate/project/models"

type Member struct {
	id   uint32
	name string
	role string
}

func NewMember(id uint32, name string, role string) *Member {
	return &Member{
		id:   id,
		name: name,
		role: role,
	}
}

func (m *Member) ToModel() models.Member {
	return models.Member{
		ID:   m.id,
		Name: m.name,
		Role: m.role,
	}
}
