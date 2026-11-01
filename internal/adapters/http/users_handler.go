package http

import (
	"encoding/json"
	"net/http"

	"github.com/caroline-oliveira/backend-tcc/internal/ports"
)

type UserHandler struct {
	service ports.AuthService
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUserHandler(service ports.AuthService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var u loginRequest
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if u.Email == "" || u.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
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
