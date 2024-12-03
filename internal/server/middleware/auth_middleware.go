package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/lunarKettle/task-management-platform-monolith/pkg/common"
)

type tokenParser = func(string) (*common.Claims, error)

var noAuthPaths = map[string]struct{}{
	"/login":    {},
	"/register": {},
}

func AuthMiddleware(next http.Handler, tokenParser tokenParser) http.Handler {
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
			fmt.Print(err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), common.ContextKeyClaims, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
