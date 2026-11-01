package ports

import (
	"context"
	"time"

	"github.com/caroline-oliveira/backend-tcc/internal/domain"
)

type TaskCreateInput struct {
	Title        string
	Subject      string
	DueDate      time.Time
	Status       domain.TaskStatus
	DepartmentID string
}

type TaskUpdateInput struct {
	Title        string
	Subject      string
	DueDate      time.Time
	Status       domain.TaskStatus
	DepartmentID string
}

type TaskSummary struct {
	Pending    int `json:"pending"`
	InProgress int `json:"in_progress"`
	Done       int `json:"done"`
}

type TaskDTO struct {
	ID           string            `json:"id"`
	Title        string            `json:"title"`
	Subject      string            `json:"subject"`
	DueDate      time.Time         `json:"due_date"`
	Status       domain.TaskStatus `json:"status"`
	AuthorID     string            `json:"author_id"`
	DepartmentID string            `json:"department_id"`
}

type TaskRepository interface {
	Create(ctx context.Context, task domain.Task) (domain.Task, error)
	GetByID(ctx context.Context, id string) (domain.Task, error)
	Update(ctx context.Context, task domain.Task) (domain.Task, error)
	Delete(ctx context.Context, id string) error
	ListByAuthor(ctx context.Context, authorID string) ([]domain.Task, error)
	ListByDepartment(ctx context.Context, departmentID string) ([]domain.Task, error)
}

type TaskService interface {
	List(ctx context.Context, accessToken string) ([]TaskDTO, error)
	Summary(ctx context.Context, accessToken string) (TaskSummary, error)
	Get(ctx context.Context, accessToken, id string) (TaskDTO, error)
	Create(ctx context.Context, accessToken string, input TaskCreateInput) (TaskDTO, error)
	Update(ctx context.Context, accessToken, id string, input TaskUpdateInput) (TaskDTO, error)
	UpdateStatus(ctx context.Context, accessToken, id string, status domain.TaskStatus) (TaskDTO, error)
	Delete(ctx context.Context, accessToken, id string) error
}
