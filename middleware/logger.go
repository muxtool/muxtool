package middleware

import (
	"net/http"

	"log/slog"
)

func Logger(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			requestID := GetRequestID(r)
			args := []any{
				slog.String("proto", r.Proto),
				slog.String("method", r.Method),
				slog.String("request_uri", r.RequestURI),
			}
			if requestID != "" {
				args = append(args, slog.String("request_id", requestID))
			}
			logger.InfoContext(
				r.Context(),
				"request",
				args...,
			)
		})
	}
}
