package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/lunarKettle/task-management-platform-monolith/pkg/common"
)

var noAuthPaths = map[string]struct{}{
	"/login":    {},
	"/register": {},
}

func authMiddleware(next http.Handler, tokenParser tokenParser) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := noAuthPaths[r.URL.Path]; ok {
			next.ServeHTTP(w, r)
			return
		}

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(tokenString, "Bearer ")

		claims, err := tokenParser(token)

		if err != nil {
			if errors.Is(err, common.ErrInvalidToken) {
				httpError := &HTTPError{
					Code:  http.StatusUnauthorized,
					Error: common.ErrInvalidToken.Error(),
				}
				WriteHTTPError(w, httpError)
			} else {
				httpError := &HTTPError{
					Code:  http.StatusUnauthorized,
					Error: "Unauthorized",
				}
				WriteHTTPError(w, httpError)
			}
			fmt.Print(err)
			return
		}

		ctx := context.WithValue(r.Context(), common.ContextKeyClaims, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
