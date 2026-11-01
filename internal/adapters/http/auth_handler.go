package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/caroline-oliveira/backend-tcc/internal/ports"
)

type AuthHandler struct {
	service ports.AuthService
}

func NewAuthHandler(service ports.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func extractBearerToken(header string) string {
	header = strings.TrimSpace(header)
	if header == "" {
		return ""
	}
	const prefix = "Bearer "
	if strings.HasPrefix(header, prefix) {
		return strings.TrimSpace(header[len(prefix):])
	}
	// Also accept lowercase "bearer "
	if len(header) > 7 && strings.EqualFold(header[:7], "bearer ") {
		return strings.TrimSpace(header[7:])
	}
	return header
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var u loginRequest
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		writeJSON(w, http.StatusBadRequest, errorBody{Error: err.Error()})
		return
	}
	if u.Email == "" || u.Password == "" {
		writeJSON(w, http.StatusBadRequest, errorBody{Error: "Email and password are required"})
		return
	}

	credentials := ports.Credentials{
		Email:    u.Email,
		Password: u.Password,
	}

	token, err := h.service.Login(r.Context(), credentials)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, ports.AuthToken{
		AccessToken:      token.AccessToken,
		RefreshToken:     token.RefreshToken,
		ExpiresIn:        token.ExpiresIn,
		ExpiresInRefresh: token.ExpiresInRefresh,
	})
}

func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	accessToken := extractBearerToken(r.Header.Get("Authorization"))
	if accessToken == "" {
		writeJSON(w, http.StatusUnauthorized, errorBody{Error: "Unauthorized"})
		return
	}

	token, err := h.service.ValidateToken(r.Context(), accessToken)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, token)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	accessToken := extractBearerToken(r.Header.Get("Authorization"))
	if accessToken == "" {
		writeJSON(w, http.StatusUnauthorized, errorBody{Error: "Unauthorized"})
		return
	}

	profile, err := h.service.Me(r.Context(), accessToken)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, profile)
}
