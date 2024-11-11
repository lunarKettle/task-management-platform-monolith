package server

import (
	"log"
	"net/http"

	"github.com/lunarKettle/task-management-platform-monolith/pkg/common"
)

type handler = func(http.ResponseWriter, *http.Request) error

func errorHandling(handler handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			switch err {
			case common.ErrNotFound:
				log.Printf("Error: %v", err)
				http.Error(w, err.Error(), http.StatusNotFound)
			case common.ErrAlreadyExists:
				log.Printf("Error: %v", err)
				http.Error(w, err.Error(), http.StatusConflict)
			case common.ErrInvalidCredentials:
				log.Printf("Error: %v", err)
				http.Error(w, err.Error(), http.StatusUnauthorized)
			case common.ErrInvalidToken:
				log.Printf("Error: %v", err)
				http.Error(w, err.Error(), http.StatusUnauthorized)
			case common.ErrUnexpectedSigningMethod:
				log.Printf("Error: %v", err)
				http.Error(w, err.Error(), http.StatusUnauthorized)
			case common.ErrTokenNotValid:
				log.Printf("Error: %v", err)
				http.Error(w, err.Error(), http.StatusUnauthorized)
			default:
				log.Printf("Unexpected error: %v", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}
	})
}
