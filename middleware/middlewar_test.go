package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_TokenValid(t *testing.T) {
	// créer une réponse fictive
	w := httptest.NewRecorder()

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "token1")

	handler := AuthenticationMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code, "Le code devrait être 200")
}

func TestAuthMiddleware_TokenInvalid(t *testing.T) {

	w := httptest.NewRecorder()

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "invalid_token")

	handler := AuthenticationMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusUnauthorized, w.Code, "Le code devrait être 401")
}
