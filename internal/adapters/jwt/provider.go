package jwt

import (
	"time"

	"github.com/caroline-oliveira/backend-tcc/internal/domain"
	"github.com/caroline-oliveira/backend-tcc/internal/ports"
	"github.com/golang-jwt/jwt/v5"
)

type TokenType string
type TokenTTL struct {
	Access  time.Duration
	Refresh time.Duration
}

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)
const (
	AccessTTL  time.Duration = time.Hour * 1
	RefreshTTL time.Duration = time.Hour * 24
)

type JWTProvider struct {
	secret     []byte
	now        func() time.Time
	accessTTL  time.Duration
	refreshTTL time.Duration
}

type Claims struct {
	Email string      `json:"email"`
	Role  domain.Role `json:"role"`
	Typ   TokenType   `json:"typ"`
	jwt.RegisteredClaims
}

func NewJWTProvider(secret []byte, now func() time.Time) (*JWTProvider, error) {
	if len(secret) == 0 {
		return nil, domain.ErrInvalidSecret
	}

	return &JWTProvider{secret: secret, now: now, accessTTL: AccessTTL, refreshTTL: RefreshTTL}, nil
}

func (p *JWTProvider) IssuePair(claims ports.TokenClaims) (ports.AuthToken, error) {
	now := p.now()
	accessToken, err := p.sign(claims, AccessToken, now, p.accessTTL)
	if err != nil {
		return ports.AuthToken{}, err
	}
	refreshToken, err := p.sign(claims, RefreshToken, now, p.refreshTTL)
	if err != nil {
		return ports.AuthToken{}, err
	}
	return ports.AuthToken{AccessToken: accessToken, RefreshToken: refreshToken, ExpiresIn: now.Add(p.accessTTL).Unix(), ExpiresInRefresh: now.Add(p.refreshTTL).Unix()}, nil
}

func (p *JWTProvider) ParseAccessToken(token string) (ports.TokenClaims, error) {
	parsed, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (any, error) {
		return p.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return ports.TokenClaims{}, domain.ErrInvalidTokenSignature
	}

	claims, ok := parsed.Claims.(*Claims)
	if !ok || !parsed.Valid {
		return ports.TokenClaims{}, domain.ErrInvalidTokenClaims
	}

	if claims.Typ != AccessToken {
		return ports.TokenClaims{}, domain.ErrInvalidTokenClaims
	}

	if claims.ExpiresAt.Before(p.now()) {
		return ports.TokenClaims{}, domain.ErrTokenExpired
	}

	return ports.TokenClaims{
		Email:     claims.Email,
		Role:      claims.Role,
		SubjectID: claims.Subject,
	}, nil
}

func (p *JWTProvider) sign(claims ports.TokenClaims, tokenType TokenType, now time.Time, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		Email: claims.Email,
		Role:  claims.Role,
		Typ:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   claims.SubjectID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	})

	signedToken, err := token.SignedString(p.secret)
	if err != nil {
		return "", domain.ErrInvalidTokenSignature
	}
	return signedToken, nil
}
