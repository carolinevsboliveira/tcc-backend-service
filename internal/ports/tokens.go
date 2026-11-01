package ports

import "github.com/caroline-oliveira/backend-tcc/internal/domain"

type TokenClaims struct {
	SubjectID string      `json:"subject_id"`
	Email     string      `json:"email"`
	Role      domain.Role `json:"role"`
}

// TokenIssuer issues and parses JWT access/refresh pairs.
type JWTProvider interface {
	IssuePair(claims TokenClaims) (AuthToken, error)
	ParseAccessToken(token string) (TokenClaims, error)
}
