package auth

import (
	"context"
	"strings"

	"github.com/caroline-oliveira/backend-tcc/internal/domain"
	"github.com/caroline-oliveira/backend-tcc/internal/ports"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepository ports.UserRepository
	tokenService   ports.JWTProvider
}

func NewAuthService(userRepository ports.UserRepository, tokenService ports.JWTProvider) ports.AuthService {
	return &authService{userRepository: userRepository, tokenService: tokenService}
}

func (s *authService) Login(ctx context.Context, credentials ports.Credentials) (ports.AuthToken, error) {
	user, err := s.validateUserCredentials(ctx, credentials)
	if err != nil {
		return ports.AuthToken{}, err
	}

	token, err := s.tokenService.IssuePair(ports.TokenClaims{
		SubjectID: user.ID,
		Email:     user.Email,
		Role:      user.Role,
	})
	if err != nil {
		return ports.AuthToken{}, err
	}

	return token, nil
}

func (s *authService) ValidateToken(ctx context.Context, accessToken string) (ports.TokenClaims, error) {
	token, err := s.tokenService.ParseAccessToken(strings.TrimSpace(accessToken))
	if err != nil {
		return ports.TokenClaims{}, err
	}

	return token, nil
}

func (s *authService) Me(ctx context.Context, accessToken string) (ports.UserProfile, error) {
	claims, err := s.ValidateToken(ctx, accessToken)
	if err != nil {
		return ports.UserProfile{}, err
	}

	user, err := s.userRepository.GetByID(ctx, claims.SubjectID)
	if err != nil || !user.Active {
		return ports.UserProfile{}, domain.ErrUnauthorized
	}

	return ports.UserProfile{Name: user.Name, Role: user.Role, DepartmentID: user.DepartmentID}, nil
}

func (s *authService) validateUserCredentials(ctx context.Context, credentials ports.Credentials) (domain.User, error) {
	user, err := s.userRepository.GetByEmail(ctx, credentials.Email)
	if err != nil {
		return domain.User{}, domain.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password)); err != nil {
		return domain.User{}, domain.ErrUnauthorized
	}

	if !user.Active {
		return domain.User{}, domain.ErrUnauthorized
	}

	return user, nil
}

func (s *authService) ValidateAccessToken(_ context.Context, token string) error {
	if _, err := s.tokenService.ParseAccessToken(strings.TrimSpace(token)); err != nil {
		return domain.ErrUnauthorized
	}
	return nil
}
