package ports

import "github.com/caroline-oliveira/backend-tcc/internal/domain"

type TokenClaims struct {
	SubjectID string
	Email     string
	Role      domain.Role
}

// TokenIssuer issues and parses JWT access/refresh pairs.
type JWTProvider interface {
	IssuePair(claims TokenClaims) (AuthToken, error)
	ParseAccessToken(token string) (TokenClaims, error)
}
