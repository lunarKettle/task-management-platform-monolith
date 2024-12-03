package server

import (
	"errors"
	"log"
	"net/http"

	"github.com/lunarKettle/task-management-platform-monolith/pkg/common"
)

type handler = func(http.ResponseWriter, *http.Request) error

func errorHandling(handler handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			switch {
			case errors.Is(err, common.ErrNotFound):
				log.Printf("Error: %v", err)
				http.Error(w, err.Error(), http.StatusNotFound)
			case errors.Is(err, common.ErrAlreadyExists):
				log.Printf("Error: %v", err)
				http.Error(w, err.Error(), http.StatusConflict)
			case errors.Is(err, common.ErrInvalidCredentials):
				log.Printf("Error: %v", err)
				http.Error(w, err.Error(), http.StatusUnauthorized)
			case errors.Is(err, common.ErrInvalidToken):
				log.Printf("Error: %v", err)
				http.Error(w, err.Error(), http.StatusUnauthorized)
			case errors.Is(err, common.ErrUnexpectedSigningMethod):
				log.Printf("Error: %v", err)
				http.Error(w, err.Error(), http.StatusUnauthorized)
			case errors.Is(err, common.ErrTokenNotValid):
				log.Printf("Error: %v", err)
				http.Error(w, err.Error(), http.StatusUnauthorized)
			case errors.Is(err, common.ErrForbidden):
				log.Printf("Error: %v", err)
				http.Error(w, err.Error(), http.StatusForbidden)
			default:
				log.Printf("Unexpected error: %v", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}
	})
}
