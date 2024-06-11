package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func Use(handler http.Handler, middlewares ...Middleware) http.Handler {
	for _, h := range middlewares {
		handler = h(handler)
	}

	return handler
}
