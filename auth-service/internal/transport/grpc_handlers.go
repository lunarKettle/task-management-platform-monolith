package transport

import (
	"context"
	"errors"

	"github.com/lunarKettle/task-management-platform/auth-service/internal/common"
	"github.com/lunarKettle/task-management-platform/auth-service/internal/usecases"
	"github.com/lunarKettle/task-management-platform/auth-service/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *GRPCServer) Register(_ context.Context, r *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	cmd := usecases.NewCreateUserCommand(r.Username, r.Email, r.Password, r.Role)
	token, err := s.usecases.CreateUser(cmd)

	if err != nil {
		if errors.Is(err, common.ErrAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}

	return &proto.RegisterResponse{
		Token: token,
	}, nil
}

func (s *GRPCServer) Authenticate(_ context.Context, r *proto.AuthRequest) (*proto.AuthResponse, error) {
	return &proto.AuthResponse{}, nil
}
