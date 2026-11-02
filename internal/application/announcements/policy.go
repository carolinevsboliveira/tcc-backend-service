package announcements

import (
	"context"

	"github.com/caroline-oliveira/backend-tcc/internal/domain"
	"github.com/caroline-oliveira/backend-tcc/internal/ports"
)

type announcementAccessPolicy interface {
	List(ctx context.Context, actor domain.User) ([]domain.Announcement, error)
	CanView(actor domain.User, announcement domain.Announcement) bool
	CanCreate(actor domain.User) bool
	CanUpdate(actor domain.User, announcement domain.Announcement) bool
	CanDelete(actor domain.User, announcement domain.Announcement) bool
	CanMarkRead(actor domain.User, announcement domain.Announcement) bool
}

func policyFor(role domain.Role, announcements ports.AnnouncementRepository) (announcementAccessPolicy, error) {
	switch role {
	case domain.RoleProfessor:
		return professorPolicy{announcements: announcements}, nil
	case domain.RoleStudent:
		return studentPolicy{announcements: announcements}, nil
	default:
		return nil, domain.ErrForbidden
	}
}

type professorPolicy struct {
	announcements ports.AnnouncementRepository
}

func (p professorPolicy) List(ctx context.Context, actor domain.User) ([]domain.Announcement, error) {
	return p.announcements.ListByAuthor(ctx, actor.ID)
}

func (p professorPolicy) CanView(actor domain.User, announcement domain.Announcement) bool {
	return announcement.AuthorID == actor.ID
}

func (p professorPolicy) CanCreate(domain.User) bool { return true }

func (p professorPolicy) CanUpdate(actor domain.User, announcement domain.Announcement) bool {
	return announcement.AuthorID == actor.ID
}

func (p professorPolicy) CanDelete(actor domain.User, announcement domain.Announcement) bool {
	return announcement.AuthorID == actor.ID
}

func (p professorPolicy) CanMarkRead(domain.User, domain.Announcement) bool { return false }

type studentPolicy struct {
	announcements ports.AnnouncementRepository
}

func (p studentPolicy) List(ctx context.Context, actor domain.User) ([]domain.Announcement, error) {
	return p.announcements.ListByDepartment(ctx, actor.DepartmentID)
}

func (p studentPolicy) CanView(actor domain.User, announcement domain.Announcement) bool {
	return announcement.DepartmentID == actor.DepartmentID
}

func (p studentPolicy) CanCreate(domain.User) bool { return false }

func (p studentPolicy) CanUpdate(domain.User, domain.Announcement) bool { return false }

func (p studentPolicy) CanDelete(domain.User, domain.Announcement) bool { return false }

func (p studentPolicy) CanMarkRead(actor domain.User, announcement domain.Announcement) bool {
	return announcement.DepartmentID == actor.DepartmentID
}
