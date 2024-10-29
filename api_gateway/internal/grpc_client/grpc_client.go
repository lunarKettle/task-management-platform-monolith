package grpc_client

import (
	"log"

	pb "api_gateway/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	authClient    pb.AuthServiceClient
	projectClient pb.ProjectManagementServiceClient
	authConn      *grpc.ClientConn
	projectConn   *grpc.ClientConn
}

func NewGRPCClient() *GRPCClient {
	const (
		projectAddress = "localhost:50051"
		authAddress    = "localhost:50052"
	)

	authConn, err := grpc.NewClient(authAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	authClient := pb.NewAuthServiceClient(authConn)

	projectConn, err := grpc.NewClient(projectAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	projectClient := pb.NewProjectManagementServiceClient(projectConn)

	return &GRPCClient{
		authClient:    authClient,
		projectClient: projectClient,
		authConn:      authConn,
		projectConn:   projectConn,
	}
}

func (g *GRPCClient) Close() {
	if err := g.projectConn.Close(); err != nil {
		log.Fatalf("failed to close project management service connection: %v", err)
	}
	if err := g.authConn.Close(); err != nil {
		log.Fatalf("failed to close auth service connection: %v", err)
	}
}
