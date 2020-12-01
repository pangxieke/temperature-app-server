package middlewares

import (
	"net/http"

	"github.com/rs/cors"
)

func Handler(next http.Handler) http.Handler {
	handlers := []func(http.Handler) http.Handler{
		Log,
		cors.AllowAll().Handler,
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler := next
		for i := len(handlers) - 1; i >= 0; i-- {
			handler = handlers[i](handler)
		}
		handler.ServeHTTP(w, r)
	})
}
