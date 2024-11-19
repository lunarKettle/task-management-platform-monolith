package dto

type CreateTeamRequestDTO struct {
	Name    string `json:"name"`
	Members []MemberDTO
}
