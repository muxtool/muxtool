package middleware

import (
	"context"
	"net/http"

	"go.neonxp.ru/objectid"
)

type ctxKeyRequestID int

const (
	RequestIDKey    ctxKeyRequestID = 0
	RequestIDHeader string          = "X-Request-ID"
)

func RequestID(next http.Handler) http.Handler {
	objectid.Seed()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(RequestIDHeader)
		if requestID == "" {
			requestID = objectid.New().String()
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), RequestIDKey, requestID)))
	})
}

func GetRequestID(r *http.Request) string {
	rid := r.Context().Value(RequestIDKey)
	if rid == nil {
		return ""
	}
	srid, ok := rid.(string)
	if !ok {
		return ""
	}

	return srid
}
