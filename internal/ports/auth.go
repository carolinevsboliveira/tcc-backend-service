package ports

import (
	"context"

	"github.com/caroline-oliveira/backend-tcc/internal/domain"
)

type Credentials struct {
	Email    string
	Password string
}

type AuthToken struct {
	AccessToken      string
	RefreshToken     string
	ExpiresIn        int64
	ExpiresInRefresh int64
}

type UserProfile struct {
	Name       string             `json:"name"`
	Role       domain.Role        `json:"role"`
	Department domain.Department  `json:"department"`
}

type AuthService interface {
	Login(ctx context.Context, credentials Credentials) (AuthToken, error)
	ValidateToken(ctx context.Context, accessToken string) (TokenClaims, error)
	Me(ctx context.Context, accessToken string) (UserProfile, error)
}
