package tasks

import (
	"context"
	"strings"
	"time"

	"github.com/caroline-oliveira/backend-tcc/internal/domain"
	"github.com/caroline-oliveira/backend-tcc/internal/ports"
)

type service struct {
	tasks ports.TaskRepository
	users ports.UserRepository
	auth  ports.AuthService
}

func NewService(tasks ports.TaskRepository, users ports.UserRepository, auth ports.AuthService) ports.TaskService {
	return &service{tasks: tasks, users: users, auth: auth}
}

func (s *service) List(ctx context.Context, accessToken string) ([]ports.TaskDTO, error) {
	actor, policy, err := s.resolveAccess(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	tasks, err := policy.List(ctx, actor)
	if err != nil {
		return nil, err
	}
	return toDTOs(tasks), nil
}

func (s *service) Summary(ctx context.Context, accessToken string) (ports.TaskSummary, error) {
	actor, policy, err := s.resolveAccess(ctx, accessToken)
	if err != nil {
		return ports.TaskSummary{}, err
	}

	tasks, err := policy.List(ctx, actor)
	if err != nil {
		return ports.TaskSummary{}, err
	}

	var summary ports.TaskSummary
	for _, task := range tasks {
		switch task.Status {
		case domain.TaskStatusPending:
			summary.Pending++
		case domain.TaskStatusInProgress:
			summary.InProgress++
		case domain.TaskStatusDone:
			summary.Done++
		}
	}
	return summary, nil
}

func (s *service) Get(ctx context.Context, accessToken, id string) (ports.TaskDTO, error) {
	actor, policy, err := s.resolveAccess(ctx, accessToken)
	if err != nil {
		return ports.TaskDTO{}, err
	}

	task, err := s.tasks.GetByID(ctx, id)
	if err != nil {
		return ports.TaskDTO{}, err
	}
	if !policy.CanView(actor, task) {
		return ports.TaskDTO{}, domain.ErrForbidden
	}
	return toDTO(task), nil
}

func (s *service) Create(ctx context.Context, accessToken string, input ports.TaskCreateInput) (ports.TaskDTO, error) {
	actor, policy, err := s.resolveAccess(ctx, accessToken)
	if err != nil {
		return ports.TaskDTO{}, err
	}
	if !policy.CanCreate(actor) {
		return ports.TaskDTO{}, domain.ErrForbidden
	}

	status := input.Status
	if status == "" {
		status = domain.TaskStatusPending
	}
	if err := validateTaskFields(input.Title, input.Subject, input.DepartmentID, input.DueDate, status); err != nil {
		return ports.TaskDTO{}, err
	}
	if _, err := s.users.GetDepartment(ctx, input.DepartmentID); err != nil {
		return ports.TaskDTO{}, domain.ErrValidation
	}

	created, err := s.tasks.Create(ctx, domain.Task{
		Title:        strings.TrimSpace(input.Title),
		Subject:      strings.TrimSpace(input.Subject),
		DueDate:      input.DueDate.UTC(),
		Status:       status,
		AuthorID:     actor.ID,
		DepartmentID: input.DepartmentID,
	})
	if err != nil {
		return ports.TaskDTO{}, err
	}
	return toDTO(created), nil
}

func (s *service) Update(ctx context.Context, accessToken, id string, input ports.TaskUpdateInput) (ports.TaskDTO, error) {
	actor, policy, err := s.resolveAccess(ctx, accessToken)
	if err != nil {
		return ports.TaskDTO{}, err
	}

	task, err := s.tasks.GetByID(ctx, id)
	if err != nil {
		return ports.TaskDTO{}, err
	}
	if !policy.CanUpdate(actor, task) {
		return ports.TaskDTO{}, domain.ErrForbidden
	}

	if err := validateTaskFields(input.Title, input.Subject, input.DepartmentID, input.DueDate, input.Status); err != nil {
		return ports.TaskDTO{}, err
	}
	if _, err := s.users.GetDepartment(ctx, input.DepartmentID); err != nil {
		return ports.TaskDTO{}, domain.ErrValidation
	}

	task.Title = strings.TrimSpace(input.Title)
	task.Subject = strings.TrimSpace(input.Subject)
	task.DueDate = input.DueDate.UTC()
	task.Status = input.Status
	task.DepartmentID = input.DepartmentID

	updated, err := s.tasks.Update(ctx, task)
	if err != nil {
		return ports.TaskDTO{}, err
	}
	return toDTO(updated), nil
}

func (s *service) UpdateStatus(ctx context.Context, accessToken, id string, status domain.TaskStatus) (ports.TaskDTO, error) {
	actor, policy, err := s.resolveAccess(ctx, accessToken)
	if err != nil {
		return ports.TaskDTO{}, err
	}
	if !domain.IsValidTaskStatus(status) {
		return ports.TaskDTO{}, domain.ErrValidation
	}

	task, err := s.tasks.GetByID(ctx, id)
	if err != nil {
		return ports.TaskDTO{}, err
	}
	if !policy.CanUpdateStatus(actor, task) {
		return ports.TaskDTO{}, domain.ErrForbidden
	}

	task.Status = status
	updated, err := s.tasks.Update(ctx, task)
	if err != nil {
		return ports.TaskDTO{}, err
	}
	return toDTO(updated), nil
}

func (s *service) Delete(ctx context.Context, accessToken, id string) error {
	actor, policy, err := s.resolveAccess(ctx, accessToken)
	if err != nil {
		return err
	}

	task, err := s.tasks.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if !policy.CanDelete(actor, task) {
		return domain.ErrForbidden
	}
	return s.tasks.Delete(ctx, id)
}

func (s *service) resolveAccess(ctx context.Context, accessToken string) (domain.User, taskAccessPolicy, error) {
	claims, err := s.auth.ValidateToken(ctx, accessToken)
	if err != nil {
		return domain.User{}, nil, err
	}

	user, err := s.users.GetByID(ctx, claims.SubjectID)
	if err != nil || !user.Active {
		return domain.User{}, nil, domain.ErrUnauthorized
	}

	policy, err := policyFor(user.Role, s.tasks)
	if err != nil {
		return domain.User{}, nil, err
	}
	return user, policy, nil
}

func validateTaskFields(title, subject, departmentID string, dueDate time.Time, status domain.TaskStatus) error {
	if strings.TrimSpace(title) == "" || strings.TrimSpace(subject) == "" || strings.TrimSpace(departmentID) == "" {
		return domain.ErrValidation
	}
	if dueDate.IsZero() {
		return domain.ErrValidation
	}
	if !domain.IsValidTaskStatus(status) {
		return domain.ErrValidation
	}
	return nil
}

func toDTO(task domain.Task) ports.TaskDTO {
	return ports.TaskDTO{
		ID:           task.ID,
		Title:        task.Title,
		Subject:      task.Subject,
		DueDate:      task.DueDate,
		Status:       task.Status,
		AuthorID:     task.AuthorID,
		DepartmentID: task.DepartmentID,
	}
}

func toDTOs(tasks []domain.Task) []ports.TaskDTO {
	out := make([]ports.TaskDTO, 0, len(tasks))
	for _, task := range tasks {
		out = append(out, toDTO(task))
	}
	return out
}
