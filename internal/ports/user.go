package ports

import (
	"context"

	"github.com/caroline-oliveira/backend-tcc/internal/domain"
)

type UserRepository interface {
	List(ctx context.Context) ([]domain.User, error)
	GetByID(ctx context.Context, id string) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	ListByRole(ctx context.Context, role domain.Role) ([]domain.User, error)
}
