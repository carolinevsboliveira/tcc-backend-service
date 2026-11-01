package domain

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
)

var (
	ErrInvalidSecret         = errors.New("invalid secret")
	ErrInvalidTokenSignature = errors.New("invalid token signature")
	ErrInvalidTokenClaims    = errors.New("invalid token claims")
	ErrTokenExpired          = errors.New("token expired")
)
