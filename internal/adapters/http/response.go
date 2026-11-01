package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/caroline-oliveira/backend-tcc/internal/domain"
)

type errorBody struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if payload == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		writeJSON(w, http.StatusNotFound, errorBody{Error: "not found"})
	case errors.Is(err, domain.ErrUnauthorized):
		writeJSON(w, http.StatusUnauthorized, errorBody{Error: "unauthorized"})
	case errors.Is(err, domain.ErrInvalidTokenSignature):
		writeJSON(w, http.StatusUnauthorized, errorBody{Error: "invalid token signature"})
	case errors.Is(err, domain.ErrInvalidTokenClaims):
		writeJSON(w, http.StatusUnauthorized, errorBody{Error: "invalid token claims"})
	case errors.Is(err, domain.ErrTokenExpired):
		writeJSON(w, http.StatusUnauthorized, errorBody{Error: "token expired"})
	default:
		writeJSON(w, http.StatusInternalServerError, errorBody{Error: "internal server error"})
	}
}
