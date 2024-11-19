package usecases

import "github.com/lunarKettle/task-management-platform-monolith/internal/aggregate/project/models"

func mapMembersToModels(members []Member) []models.Member {
	memberModels := make([]models.Member, len(members))

	for i, v := range members {
		memberModels[i] = v.ToModel()
	}

	return memberModels
}
