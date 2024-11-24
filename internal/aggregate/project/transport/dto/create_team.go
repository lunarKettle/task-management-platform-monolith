package dto

type CreateTeamRequestDTO struct {
	Name      string      `json:"name"`
	Members   []MemberDTO `json:"members"`
	ManagerID uint32      `json:"managerId"`
}
