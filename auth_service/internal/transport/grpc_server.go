package transport

import (
	"fmt"
	"net"

	"auth_service/internal/usecases"
	pb "auth_service/proto"
	"log"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	pb.AuthServiceServer
	usecases *usecases.AuthUseCases
}

func NewGRPCServer(usecases *usecases.AuthUseCases) *GRPCServer {
	server := &GRPCServer{
		usecases: usecases,
	}
	return server
}

func (s *GRPCServer) Start(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	defer lis.Close()

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, s)

	log.Printf("server listening at %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}
	return nil
}
