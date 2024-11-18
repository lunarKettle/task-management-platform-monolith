package transport

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/lunarKettle/task-management-platform-monolith/internal/aggregate/user/models/dto"
)

func extractBasicAuth(r *http.Request) (*dto.LoginUserRequestDTO, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("authorization header missing")
	}

	if len(authHeader) < 7 || authHeader[:6] != "Basic " {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	encodedCredentials := authHeader[6:]
	decodedCredentials, err := base64.StdEncoding.DecodeString(encodedCredentials)
	if err != nil {
		return nil, fmt.Errorf("failed to decode authorization header: %w", err)
	}

	credentials := strings.SplitN(string(decodedCredentials), ":", 2)
	if len(credentials) != 2 {
		return nil, fmt.Errorf("invalid authorization credentials")
	}

	return &dto.LoginUserRequestDTO{
		Username: credentials[0],
		Password: credentials[1],
	}, nil
}
