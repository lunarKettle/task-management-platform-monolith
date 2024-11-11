package server

import "net/http"

type Handler interface {
	RegisterRoutes(mux *http.ServeMux, eh func(handler) http.Handler)
}
