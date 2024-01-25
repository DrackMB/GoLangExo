package middleware

import (
	"log"
	"net/http"
)

var validTokens = map[string]bool{
	"token1": true,
	"token2": true,
}

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the request header
		token := r.Header.Get("Authorization")

		// Validate the token
		if !isValidToken(token) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Forward request to the next handler
		next.ServeHTTP(w, r)
	})
}

func isValidToken(token string) bool {
	// Check if the token is valid
	if _, ok := validTokens[token]; !ok {
		return false
	}

	return true
}
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Write log entry
		log.Println("Incoming request from", r.RemoteAddr, r.Method, r.URL.Path)

		// Validate the token
		if !isValidToken(r.Header.Get("Authorization")) {
			log.Println("Invalid token")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Forward request to the next handler
		next.ServeHTTP(w, r)
	})
}
