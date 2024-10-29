package handler

import (
	"api_gateway/internal/models/dto"
	pb "api_gateway/proto"
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *HTTPServer) registerUser(w http.ResponseWriter, r *http.Request) error {
	var regUserReq dto.RegisterUserRequestDTO
	err := json.NewDecoder(r.Body).Decode(&regUserReq)
	if err != nil {
		return fmt.Errorf("error while decoding request body: %w", err)
	}
	defer r.Body.Close()

	grpcRequest := &pb.RegisterRequest{
		Email:    regUserReq.Email,
		Password: regUserReq.Password,
		Username: regUserReq.Username,
		Role:     regUserReq.Role,
	}

	grpcResponse, err := s.grpcClient.Register(grpcRequest)

	if err != nil {
		return fmt.Errorf("failed to register user: %w", err)

	}

	reqUserResp := dto.RegisterUserResponseDTO{
		AccessToken: grpcResponse.Token,
	}
	if err := json.NewEncoder(w).Encode(reqUserResp); err != nil {
		return fmt.Errorf("failed to encode response to JSON: %w", err)
	}
	return err
}

func (s *HTTPServer) authenticate(w http.ResponseWriter, r *http.Request) error {
	var authUserReq dto.LoginUserRequestDTO
	err := json.NewDecoder(r.Body).Decode(&authUserReq)
	if err != nil {
		return fmt.Errorf("error while decoding request body: %w", err)
	}
	defer r.Body.Close()

	grpcRequest := &pb.AuthRequest{
		Username: authUserReq.Username,
		Password: authUserReq.Password,
	}

	grpcResponse, err := s.grpcClient.Authenticate(grpcRequest)

	reqUserResp := dto.LoginUserResponseDTO{
		AccessToken: grpcResponse.Token,
	}
	if err := json.NewEncoder(w).Encode(reqUserResp); err != nil {
		return fmt.Errorf("failed to encode response to JSON: %w", err)
	}
	return err
}
