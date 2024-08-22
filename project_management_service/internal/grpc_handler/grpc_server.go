package grpc_handler

import (
	"context"
	"log"
	"net"

	pb "project_management_service/proto"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	pb.ProjectManagementServiceServer
}

func NewGRPCServer() *GRPCServer {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	server := &GRPCServer{}
	pb.RegisterProjectManagementServiceServer(s, server)
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return server
}

func (s *GRPCServer) GetProject(ctx context.Context, request *pb.ProjectRequest) (*pb.ProjectResponse, error) {
	return &pb.ProjectResponse{
		ProjectId:          request.ProjectId,
		ProjectName:        "test name",
		ProjectDescription: "test descrition",
	}, nil
}
