package grpc_client

import (
	"api_gateway/internal/models"
	pb "api_gateway/proto"
	"context"
	"log"
	"time"

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
		log.Fatalf("error: %v", err)
	}
	project := models.Project{
		Id:          r.GetProjectId(),
		Name:        r.GetProjectName(),
		Description: r.GetProjectDescription(),
	}
	return project, err
}
