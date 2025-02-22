package middleware

import (
	"log/slog"
	"net/http"
)

func SetHeader(key, value string) Middleware {
	return func(logger *slog.Logger, next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(key, value)
			next.ServeHTTP(w, r)
		})
	}
}
