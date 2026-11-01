package memory

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"golang.org/x/crypto/bcrypt"

	"github.com/caroline-oliveira/backend-tcc/internal/domain"
)

type UserRepository struct {
	mu      sync.RWMutex
	byID    map[string]domain.User
	byEmail map[string]string
}

type seedUser struct {
	ID           string
	Name         string
	Email        string
	Role         domain.Role
	DepartmentID string
	Active       bool
}

func NewUserRepository(seedPassword string) (*UserRepository, error) {
	if seedPassword == "" {
		return nil, fmt.Errorf("seed password is required")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(seedPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash seed password: %w", err)
	}
	passwordHash := string(hash)

	seeds := []seedUser{
		{ID: "stu-001", Name: "Ana Souza", Email: "ana.souza@example.com", Role: domain.RoleStudent, DepartmentID: "dept-001", Active: true},
		{ID: "stu-002", Name: "Bruno Lima", Email: "bruno.lima@example.com", Role: domain.RoleStudent, DepartmentID: "dept-002", Active: true},
		{ID: "stu-003", Name: "Carla Mendes", Email: "carla.mendes@example.com", Role: domain.RoleStudent, DepartmentID: "dept-003", Active: false},
		{ID: "prf-001", Name: "Dr. Helena Castro", Email: "helena.castro@example.com", Role: domain.RoleProfessor, DepartmentID: "dept-001", Active: true},
		{ID: "prf-002", Name: "Dr. Igor Nunes", Email: "igor.nunes@example.com", Role: domain.RoleProfessor, DepartmentID: "dept-004", Active: true},
		{ID: "prf-003", Name: "Dra. Juliana Prado", Email: "juliana.prado@example.com", Role: domain.RoleProfessor, DepartmentID: "dept-003", Active: false},
	}

	byID := make(map[string]domain.User, len(seeds))
	byEmail := make(map[string]string, len(seeds))
	for _, seed := range seeds {
		user := domain.User{
			ID:           seed.ID,
			Name:         seed.Name,
			Email:        seed.Email,
			PasswordHash: passwordHash,
			Role:         seed.Role,
			DepartmentID: seed.DepartmentID,
			Active:       seed.Active,
		}
		byID[user.ID] = user
		byEmail[user.Email] = user.ID
	}

	return &UserRepository{byID: byID, byEmail: byEmail}, nil
}


func (r *UserRepository) List(_ context.Context) ([]domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]domain.User, 0, len(r.byID))
	for _, user := range r.byID {
		out = append(out, user)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

func (r *UserRepository) GetByID(_ context.Context, id string) (domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.byID[id]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(_ context.Context, email string) (domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.byEmail[email]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	return r.byID[id], nil
}

func (r *UserRepository) ListByRole(_ context.Context, role domain.Role) ([]domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]domain.User, 0)
	for _, user := range r.byID {
		if user.Role == role {
			out = append(out, user)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}
