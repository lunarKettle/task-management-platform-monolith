package grpc_client

import (
	"api_gateway/internal/models"
	"context"
	"fmt"
	"log"
	"time"

	pb "api_gateway/proto"

	"google.golang.org/protobuf/types/known/timestamppb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	client pb.ProjectManagementServiceClient
	conn   *grpc.ClientConn
}

func NewGRPCClient() *GRPCClient {
	const (
		address = "localhost:50051"
	)
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := pb.NewProjectManagementServiceClient(conn)
	return &GRPCClient{c, conn}
}

func (g *GRPCClient) Close() {
	if err := g.conn.Close(); err != nil {
		log.Fatalf("failed to close connection: %v", err)
	}
}

func (g *GRPCClient) GetProject(id uint32) (models.Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := g.client.GetProject(ctx, &pb.ProjectRequest{ProjectId: id})
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
	r, err := g.client.CreateProject(ctx, &pb.CreateProjectRequest{
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
		log.Fatalf("error: %v", err)
	}
	return r.GetProjectId(), err
}
