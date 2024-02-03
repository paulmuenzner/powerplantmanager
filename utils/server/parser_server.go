package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	config "github.com/paulmuenzner/powerplantmanager/config"
	cookie "github.com/paulmuenzner/powerplantmanager/utils/cookies"
)

// ParserBodyRequest is a middleware for parsing the request body globally
func BodyRequestParser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check content type
		contentType := r.Header.Get("Content-Type")

		if strings.Contains(contentType, "application/json") {
			// Parse the request body as JSON
			var data map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				http.Error(w, "Invalid JSON in request body.", http.StatusBadRequest)
				return
			}

			// Attach JSON data to request context
			r = r.WithContext(context.WithValue(r.Context(), "requestBody", data))
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// CookieParser is a middleware that reads the user data from the cookie
// and attaches it to the request context for use in controllers
func CookieParser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read the cookie
		authData, found := cookie.GetCookie(r, config.AuthCookieName)

		if found {
			// Attach the decoded data to the request context
			const AuthCookieName string = config.AuthCookieName
			ctx := context.WithValue(r.Context(), AuthCookieName, authData)
			r = r.WithContext(ctx)
		}

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
