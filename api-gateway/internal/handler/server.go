package handler

import (
	"net/http"

	"github.com/lunarKettle/task-management-platform/api-gateway/internal/grpc_client"
	"github.com/lunarKettle/task-management-platform/api-gateway/internal/middleware"
)

type HTTPServer struct {
	Address    string
	grpcClient *grpc_client.GRPCClient
}

func NewServer(addr string, c *grpc_client.GRPCClient) *HTTPServer {
	return &HTTPServer{Address: addr, grpcClient: c}
}

func (s *HTTPServer) Start() error {
	mux := http.NewServeMux()
	eh := errorHandling

	mux.Handle("POST /users/register", eh(s.registerUser))
	mux.Handle("POST /users/login", eh(s.authenticate))

	mux.Handle("GET /projects/{id}", eh(s.getProject))
	mux.Handle("POST /projects", eh(s.createProject))
	mux.Handle("PUT /projects", eh(s.updateProject))
	mux.Handle("DELETE /projects/{id}", eh(s.deleteProject))

	contentTypeMux := middleware.ContentTypeMiddleware(mux)

	return http.ListenAndServe(s.Address, contentTypeMux)
}
