package tasks

import (
	"context"

	"github.com/caroline-oliveira/backend-tcc/internal/domain"
	"github.com/caroline-oliveira/backend-tcc/internal/ports"
)

type taskAccessPolicy interface {
	List(ctx context.Context, actor domain.User) ([]domain.Task, error)
	CanView(actor domain.User, task domain.Task) bool
	CanCreate(actor domain.User) bool
	CanUpdate(actor domain.User, task domain.Task) bool
	CanUpdateStatus(actor domain.User, task domain.Task) bool
	CanDelete(actor domain.User, task domain.Task) bool
}

func policyFor(role domain.Role, tasks ports.TaskRepository) (taskAccessPolicy, error) {
	switch role {
	case domain.RoleProfessor:
		return professorPolicy{tasks: tasks}, nil
	case domain.RoleStudent:
		return studentPolicy{tasks: tasks}, nil
	default:
		return nil, domain.ErrForbidden
	}
}

type professorPolicy struct {
	tasks ports.TaskRepository
}

func (p professorPolicy) List(ctx context.Context, actor domain.User) ([]domain.Task, error) {
	return p.tasks.ListByAuthor(ctx, actor.ID)
}

func (p professorPolicy) CanView(actor domain.User, task domain.Task) bool {
	return task.AuthorID == actor.ID
}

func (p professorPolicy) CanCreate(domain.User) bool { return true }

func (p professorPolicy) CanUpdate(actor domain.User, task domain.Task) bool {
	return task.AuthorID == actor.ID
}

func (p professorPolicy) CanUpdateStatus(actor domain.User, task domain.Task) bool {
	return task.AuthorID == actor.ID
}

func (p professorPolicy) CanDelete(actor domain.User, task domain.Task) bool {
	return task.AuthorID == actor.ID
}

type studentPolicy struct {
	tasks ports.TaskRepository
}

func (p studentPolicy) List(ctx context.Context, actor domain.User) ([]domain.Task, error) {
	return p.tasks.ListByDepartment(ctx, actor.DepartmentID)
}

func (p studentPolicy) CanView(actor domain.User, task domain.Task) bool {
	return task.DepartmentID == actor.DepartmentID
}

func (p studentPolicy) CanCreate(domain.User) bool { return false }

func (p studentPolicy) CanUpdate(domain.User, domain.Task) bool { return false }

func (p studentPolicy) CanUpdateStatus(actor domain.User, task domain.Task) bool {
	return task.DepartmentID == actor.DepartmentID
}

func (p studentPolicy) CanDelete(domain.User, domain.Task) bool { return false }
