package handler

import (
	"log"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler = func(http.ResponseWriter, *http.Request) error

func errorHandling(handler Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			grpcErr, ok := status.FromError(err)
			if ok {
				log.Printf("gRPC error: %v", grpcErr.Message())

				switch grpcErr.Code() {
				case codes.NotFound:
					http.Error(w, grpcErr.Message(), http.StatusNotFound)
				case codes.Internal:
					http.Error(w, "internal server error", http.StatusInternalServerError)
				case codes.InvalidArgument:
					http.Error(w, "invalid argument", http.StatusBadRequest)
				default:
					http.Error(w, "unexpected gRPC error", http.StatusInternalServerError)
				}
				return
			}

			log.Printf("non-gRPC error: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	})
}
