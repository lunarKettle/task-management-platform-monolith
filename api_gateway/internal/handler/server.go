package handler

import (
	"api_gateway/internal/grpc_client"
	"net/http"
)

type HTTPServer struct {
	Address string
	client  *grpc_client.GRPCClient
}

func NewServer(addr string, c *grpc_client.GRPCClient) *HTTPServer {
	return &HTTPServer{Address: addr, client: c}
}

func (s *HTTPServer) Start() error {
	mux := http.NewServeMux()
	eh := errorHandling
	mux.Handle("GET /projects/{id}", eh(s.getProject))

	return http.ListenAndServe(s.Address, mux)
}
