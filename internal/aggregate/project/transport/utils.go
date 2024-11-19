package transport

import (
	"github.com/lunarKettle/task-management-platform-monolith/internal/aggregate/project/models"
	"github.com/lunarKettle/task-management-platform-monolith/internal/aggregate/project/transport/dto"
)

func memberModelToDTO(member models.Member) dto.MemberDTO {
	return dto.MemberDTO{
		ID:   member.ID,
		Role: member.Role,
	}
}
