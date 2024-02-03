package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

// CheckAllowedMethodsMiddleware is a custom middleware to check allowed request methods
func CheckAllowedMethodsMiddleware(allowedMethods map[string]bool) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !allowedMethods[r.Method] {
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
