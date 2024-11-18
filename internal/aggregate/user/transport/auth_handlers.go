package transport

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lunarKettle/task-management-platform-monolith/internal/aggregate/user/models/dto"
	"github.com/lunarKettle/task-management-platform-monolith/internal/aggregate/user/usecases"
)

type AuthHandlers struct {
	usecases *usecases.AuthUseCases
}

func NewAuthHandlers(usecases *usecases.AuthUseCases) *AuthHandlers {
	return &AuthHandlers{
		usecases: usecases,
	}
}

type handler = func(w http.ResponseWriter, r *http.Request) error

func (h *AuthHandlers) RegisterRoutes(mux *http.ServeMux, eh func(handler) http.Handler) {
	mux.Handle("POST /users/register", eh(h.registerUser))
	mux.Handle("POST /users/login", eh(h.authenticate))
}

func (h *AuthHandlers) registerUser(w http.ResponseWriter, r *http.Request) error {
	var regUserReq dto.RegisterUserRequestDTO
	err := json.NewDecoder(r.Body).Decode(&regUserReq)
	if err != nil {
		return fmt.Errorf("error while decoding request body: %w", err)
	}
	defer r.Body.Close()

	cmd := usecases.NewCreateUserCommand(regUserReq.Username, regUserReq.Email, regUserReq.Password, regUserReq.Role)
	token, err := h.usecases.CreateUser(cmd)
	if err != nil {
		return fmt.Errorf("failed to register user: %v", err)
	}

	reqUserResp := dto.RegisterUserResponseDTO{
		AccessToken: token,
	}
	if err := json.NewEncoder(w).Encode(reqUserResp); err != nil {
		return fmt.Errorf("failed to encode response to JSON: %w", err)
	}
	return err
}

func (h *AuthHandlers) authenticate(w http.ResponseWriter, r *http.Request) error {
	authUserReq, err := extractBasicAuth(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return fmt.Errorf("basic auth header is missing or malformed: %w", err)
	}

	token, err := h.usecases.AuthenticateUser(authUserReq.Username, authUserReq.Password)

	if err != nil {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return fmt.Errorf("failed to authenticate user: %w", err)
	}

	reqUserResp := dto.LoginUserResponseDTO{
		AccessToken: token,
	}
	if err := json.NewEncoder(w).Encode(reqUserResp); err != nil {
		return fmt.Errorf("failed to encode response to JSON: %w", err)
	}
	return err
}
