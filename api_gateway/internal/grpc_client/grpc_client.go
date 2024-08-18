package grpc_client

import (
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

func (g *GRPCClient) GetProject() (*pb.ProjectResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := g.client.GetProject(ctx, &pb.ProjectRequest{})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return r, err
}
