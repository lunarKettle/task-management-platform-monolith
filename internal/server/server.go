package server

import (
	"net/http"

	"github.com/lunarKettle/task-management-platform-monolith/internal/server/middleware"
	"github.com/lunarKettle/task-management-platform-monolith/pkg/common"
)

type tokenParser = func(string) (*common.Claims, error)

type HTTPServer struct {
	address     string
	tokenParser tokenParser
}

func NewServer(addr string, tokenParser tokenParser) *HTTPServer {
	return &HTTPServer{
		address:     addr,
		tokenParser: tokenParser,
	}
}

func (s *HTTPServer) Address() string {
	return s.address
}

func (s *HTTPServer) Start(handlers ...Handler) error {
	mux := http.NewServeMux()

	for _, handler := range handlers {
		handler.RegisterRoutes(mux, errorHandling)
	}

	authMux := middleware.AuthMiddleware(mux, s.tokenParser)

	authAndLoggingMux := middleware.LoggingMiddleware(authMux)

	finalMux := middleware.CORSMiddleware(authAndLoggingMux)

	return http.ListenAndServe(s.address, finalMux)
}
