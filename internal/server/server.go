package server

import (
	"net/http"

	"github.com/lunarKettle/task-management-platform-monolith/internal/server/middleware"
)

type HTTPServer struct {
	Address string
}

func NewServer(addr string) *HTTPServer {
	return &HTTPServer{
		Address: addr,
	}
}

func (s *HTTPServer) Start(handlers ...Handler) error {
	mux := http.NewServeMux()

	for _, handler := range handlers {
		handler.RegisterRoutes(mux, errorHandling)
	}

	contentTypeMux := middleware.ContentTypeMiddleware(mux)

	return http.ListenAndServe(s.Address, contentTypeMux)
}
