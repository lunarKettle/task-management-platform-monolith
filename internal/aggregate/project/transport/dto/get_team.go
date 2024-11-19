package dto

type GetTeamResponseDTO struct {
	ID      uint32      `json:"id"`
	Name    string      `json:"name"`
	Members []MemberDTO `json:"members"`
}
