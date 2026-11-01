package http

import (
	"net/http"
	"time"
)

type Dependencies struct {
	Auth *AuthHandler
}

// NewRouter builds the HTTP mux for the API.
func NewRouter(deps Dependencies) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
			"time":   time.Now().UTC().Format(time.RFC3339),
		})
	})

	mux.HandleFunc("POST /auth/login", deps.Auth.Login)
	mux.HandleFunc("GET /auth/validate", deps.Auth.ValidateToken)
	mux.HandleFunc("GET /auth/me", deps.Auth.Me)

	return mux
}
