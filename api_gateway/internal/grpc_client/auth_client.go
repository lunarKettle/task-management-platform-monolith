package grpc_client

import (
	"api_gateway/proto"
	"context"
	"time"
)

func (g *GRPCClient) Authenticate(r *proto.AuthRequest) (*proto.AuthResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := g.authClient.Authenticate(ctx, r)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (g *GRPCClient) Register(r *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := g.authClient.Register(ctx, r)
	if err != nil {
		return nil, err
	}
	return response, nil
}
