package handler

import (
	"net/http"
)

type HTTPServer struct {
	Address string
}

func NewServer(addr string) *HTTPServer {
	return &HTTPServer{Address: addr}
}

func (s *HTTPServer) Start() error {
	mux := http.NewServeMux()
	eh := errorHandling
	mux.Handle("GET /projects", eh(s.getProjects))

	return http.ListenAndServe(s.Address, mux)
}
