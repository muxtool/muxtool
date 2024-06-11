package middleware

import (
	"net/http"
	"runtime/debug"

	"log/slog"
)

func Recover(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				err := recover()
				if err == nil {
					return
				}
				debug.PrintStack()
				requestID := GetRequestID(r)
				logger.ErrorContext(
					r.Context(),
					"panic",
					slog.Any("panic", err),
					slog.String("proto", r.Proto),
					slog.String("method", r.Method),
					slog.String("request_uri", r.RequestURI),
					slog.String("request_id", requestID),
				)
			}()

			next.ServeHTTP(w, r)
		})
	}
}
