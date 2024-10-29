package grpc_client

import (
	"context"
	"fmt"
	"time"

	"github.com/lunarKettle/task-management-platform/api-gateway/internal/models"
	pb "github.com/lunarKettle/task-management-platform/api-gateway/proto"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (g *GRPCClient) GetProject(id uint32) (models.Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := g.projectClient.GetProject(ctx, &pb.ProjectRequest{ProjectId: id})
	if err != nil {
		return models.Project{}, err
	}

	if r == nil {
		return models.Project{}, fmt.Errorf("received nil response for project with id %d", id)
	}

	project := models.Project{
		Id:             r.Project.GetProjectId(),
		Name:           r.Project.GetProjectName(),
		Description:    r.Project.GetProjectDescription(),
		StartDate:      r.Project.GetStartDate().AsTime(),
		PlannedEndDate: r.Project.GetPlannedEndDate().AsTime(),
		ActualEndDate:  r.Project.GetActualEndDate().AsTime(),
		Status:         r.Project.GetStatus(),
		Priority:       r.Project.GetPriority(),
		TeamId:         r.Project.GetTeamId(),
		Budget:         r.Project.GetBudget(),
	}

	return project, nil
}

func (g *GRPCClient) CreateProject(project models.Project) (uint32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := g.projectClient.CreateProject(ctx, &pb.CreateProjectRequest{
		ProjectName:        project.Name,
		ProjectDescription: project.Description,
		StartDate:          timestamppb.New(project.StartDate),
		PlannedEndDate:     timestamppb.New(project.PlannedEndDate),
		ActualEndDate:      timestamppb.New(project.ActualEndDate),
		Status:             project.Status,
		Priority:           project.Priority,
		TeamId:             project.TeamId,
		Budget:             project.Budget,
	})
	if err != nil {
		return 0, err
	}
	return r.GetProjectId(), nil
}

func (g *GRPCClient) UpdateProject(project models.Project) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := g.projectClient.UpdateProject(ctx, &pb.UpdateProjectRequest{
		Project: &pb.Project{
			ProjectId:          project.Id,
			ProjectName:        project.Name,
			ProjectDescription: project.Description,
			StartDate:          timestamppb.New(project.StartDate),
			PlannedEndDate:     timestamppb.New(project.PlannedEndDate),
			ActualEndDate:      timestamppb.New(project.ActualEndDate),
			Status:             project.Status,
			Priority:           project.Priority,
			TeamId:             project.TeamId,
			Budget:             project.Budget,
		}})
	if err != nil {
		return err
	}
	return nil
}

func (g *GRPCClient) DeleteProject(id uint32) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := g.projectClient.DeleteProject(ctx, &pb.ProjectRequest{ProjectId: id})
	if err != nil {
		return err
	}
	return nil
}
