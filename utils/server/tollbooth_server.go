package server

import (
	"net/http"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gorilla/mux"
)

// TollboothMiddleware is a custom middleware that wraps Gorilla Mux router with tollbooth rate limiting
func TollboothMiddleware(lim *limiter.Limiter) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			httpError := tollbooth.LimitByRequest(lim, w, r)
			if httpError != nil {
				http.Error(w, httpError.Message, httpError.StatusCode)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
