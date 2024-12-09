package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/lunarKettle/task-management-platform-monolith/pkg/common"
)

type handler = func(http.ResponseWriter, *http.Request) error

func errorHandling(handler handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			log.Printf("Error: %v", err)
			var (
				code         int
				errorMessage string
			)

			switch {
			case errors.Is(err, common.ErrNotFound):
				errorMessage = common.ErrNotFound.Error()
				code = http.StatusNotFound

			case errors.Is(err, common.ErrAlreadyExists):
				errorMessage = common.ErrAlreadyExists.Error()
				code = http.StatusConflict

			case errors.Is(err, common.ErrInvalidCredentials):
				errorMessage = common.ErrInvalidCredentials.Error()
				code = http.StatusUnauthorized

			case errors.Is(err, common.ErrInvalidToken):
				errorMessage = common.ErrInvalidToken.Error()
				code = http.StatusUnauthorized

			case errors.Is(err, common.ErrUnexpectedSigningMethod):
				errorMessage = common.ErrUnexpectedSigningMethod.Error()
				code = http.StatusUnauthorized

			case errors.Is(err, common.ErrTokenNotValid):
				errorMessage = common.ErrTokenNotValid.Error()
				code = http.StatusUnauthorized

			case errors.Is(err, common.ErrForbidden):
				errorMessage = common.ErrForbidden.Error()
				code = http.StatusForbidden

			default:
				errorMessage = err.Error()
				code = http.StatusInternalServerError
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}

			log.Printf("Error: %v", err)

			httpError := &HTTPError{
				Code:  code,
				Error: errorMessage,
			}
			WriteHTTPError(w, httpError)
		}
	})
}

type HTTPError struct {
	Code        int    `json:"code"`
	Error       string `json:"error"`
	Description string `json:"description,omitempty"`
}

func WriteHTTPError(w http.ResponseWriter, err *HTTPError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)
	json.NewEncoder(w).Encode(err)
}
