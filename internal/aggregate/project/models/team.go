package models

type Team struct {
	ID        uint32
	Name      string
	Members   []Member
	ManagerID uint32
}
