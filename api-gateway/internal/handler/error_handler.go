package handler

import (
	"net/http"
)

type Handler = func(http.ResponseWriter, *http.Request) error

func errorHandling(handler Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {

		}
	})
}
