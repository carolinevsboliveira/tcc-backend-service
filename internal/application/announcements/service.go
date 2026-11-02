package announcements

import (
	"context"
	"strings"
	"time"

	"github.com/caroline-oliveira/backend-tcc/internal/domain"
	"github.com/caroline-oliveira/backend-tcc/internal/ports"
)

type service struct {
	announcements ports.AnnouncementRepository
	users         ports.UserRepository
	auth          ports.AuthService
}

func NewService(announcements ports.AnnouncementRepository, users ports.UserRepository, auth ports.AuthService) ports.AnnouncementService {
	return &service{announcements: announcements, users: users, auth: auth}
}

func (s *service) List(ctx context.Context, accessToken string) ([]ports.AnnouncementDTO, error) {
	actor, policy, err := s.resolveAccess(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	items, err := policy.List(ctx, actor)
	if err != nil {
		return nil, err
	}
	return s.toDTOs(ctx, actor.ID, items)
}

func (s *service) Get(ctx context.Context, accessToken, id string) (ports.AnnouncementDTO, error) {
	actor, policy, err := s.resolveAccess(ctx, accessToken)
	if err != nil {
		return ports.AnnouncementDTO{}, err
	}

	announcement, err := s.announcements.GetByID(ctx, id)
	if err != nil {
		return ports.AnnouncementDTO{}, err
	}
	if !policy.CanView(actor, announcement) {
		return ports.AnnouncementDTO{}, domain.ErrForbidden
	}
	return s.toDTO(ctx, actor.ID, announcement)
}

func (s *service) Create(ctx context.Context, accessToken string, input ports.AnnouncementCreateInput) (ports.AnnouncementDTO, error) {
	actor, policy, err := s.resolveAccess(ctx, accessToken)
	if err != nil {
		return ports.AnnouncementDTO{}, err
	}
	if !policy.CanCreate(actor) {
		return ports.AnnouncementDTO{}, domain.ErrForbidden
	}
	if err := validateAnnouncementFields(input.Title, input.Body, input.DepartmentID); err != nil {
		return ports.AnnouncementDTO{}, err
	}
	if _, err := s.users.GetDepartment(ctx, input.DepartmentID); err != nil {
		return ports.AnnouncementDTO{}, domain.ErrValidation
	}

	created, err := s.announcements.Create(ctx, domain.Announcement{
		Title:        strings.TrimSpace(input.Title),
		Body:         strings.TrimSpace(input.Body),
		AuthorID:     actor.ID,
		DepartmentID: input.DepartmentID,
		CreatedAt:    time.Now().UTC(),
	})
	if err != nil {
		return ports.AnnouncementDTO{}, err
	}
	return s.toDTO(ctx, actor.ID, created)
}

func (s *service) Update(ctx context.Context, accessToken, id string, input ports.AnnouncementUpdateInput) (ports.AnnouncementDTO, error) {
	actor, policy, err := s.resolveAccess(ctx, accessToken)
	if err != nil {
		return ports.AnnouncementDTO{}, err
	}

	announcement, err := s.announcements.GetByID(ctx, id)
	if err != nil {
		return ports.AnnouncementDTO{}, err
	}
	if !policy.CanUpdate(actor, announcement) {
		return ports.AnnouncementDTO{}, domain.ErrForbidden
	}
	if err := validateAnnouncementFields(input.Title, input.Body, input.DepartmentID); err != nil {
		return ports.AnnouncementDTO{}, err
	}
	if _, err := s.users.GetDepartment(ctx, input.DepartmentID); err != nil {
		return ports.AnnouncementDTO{}, domain.ErrValidation
	}

	announcement.Title = strings.TrimSpace(input.Title)
	announcement.Body = strings.TrimSpace(input.Body)
	announcement.DepartmentID = input.DepartmentID

	updated, err := s.announcements.Update(ctx, announcement)
	if err != nil {
		return ports.AnnouncementDTO{}, err
	}
	return s.toDTO(ctx, actor.ID, updated)
}

func (s *service) Delete(ctx context.Context, accessToken, id string) error {
	actor, policy, err := s.resolveAccess(ctx, accessToken)
	if err != nil {
		return err
	}

	announcement, err := s.announcements.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if !policy.CanDelete(actor, announcement) {
		return domain.ErrForbidden
	}
	return s.announcements.Delete(ctx, id)
}

func (s *service) MarkRead(ctx context.Context, accessToken, id string, input ports.MarkReadInput) (ports.AnnouncementDTO, error) {
	actor, policy, err := s.resolveAccess(ctx, accessToken)
	if err != nil {
		return ports.AnnouncementDTO{}, err
	}

	announcement, err := s.announcements.GetByID(ctx, id)
	if err != nil {
		return ports.AnnouncementDTO{}, err
	}
	if !policy.CanMarkRead(actor, announcement) {
		return ports.AnnouncementDTO{}, domain.ErrForbidden
	}

	if err := s.announcements.MarkRead(ctx, domain.AnnouncementRead{
		AnnouncementID: id,
		UserID:         actor.ID,
		Text:           strings.TrimSpace(input.Text),
		ReadAt:         time.Now().UTC(),
	}); err != nil {
		return ports.AnnouncementDTO{}, err
	}
	return s.toDTO(ctx, actor.ID, announcement)
}

func (s *service) Home(ctx context.Context, accessToken string) (ports.HomeSummary, error) {
	actor, policy, err := s.resolveAccess(ctx, accessToken)
	if err != nil {
		return ports.HomeSummary{}, err
	}

	items, err := policy.List(ctx, actor)
	if err != nil {
		return ports.HomeSummary{}, err
	}

	unread := 0
	for _, item := range items {
		read, err := s.announcements.IsRead(ctx, item.ID, actor.ID)
		if err != nil {
			return ports.HomeSummary{}, err
		}
		if !read {
			unread++
		}
	}
	return ports.HomeSummary{Unread: unread}, nil
}

func (s *service) resolveAccess(ctx context.Context, accessToken string) (domain.User, announcementAccessPolicy, error) {
	claims, err := s.auth.ValidateToken(ctx, accessToken)
	if err != nil {
		return domain.User{}, nil, err
	}

	user, err := s.users.GetByID(ctx, claims.SubjectID)
	if err != nil || !user.Active {
		return domain.User{}, nil, domain.ErrUnauthorized
	}

	policy, err := policyFor(user.Role, s.announcements)
	if err != nil {
		return domain.User{}, nil, err
	}
	return user, policy, nil
}

func validateAnnouncementFields(title, body, departmentID string) error {
	if strings.TrimSpace(title) == "" || strings.TrimSpace(body) == "" || strings.TrimSpace(departmentID) == "" {
		return domain.ErrValidation
	}
	return nil
}

func (s *service) toDTO(ctx context.Context, userID string, announcement domain.Announcement) (ports.AnnouncementDTO, error) {
	read, err := s.announcements.IsRead(ctx, announcement.ID, userID)
	if err != nil {
		return ports.AnnouncementDTO{}, err
	}
	return ports.AnnouncementDTO{
		ID:           announcement.ID,
		Title:        announcement.Title,
		Body:         announcement.Body,
		CreatedAt:    announcement.CreatedAt,
		AuthorID:     announcement.AuthorID,
		DepartmentID: announcement.DepartmentID,
		Read:         read,
	}, nil
}

func (s *service) toDTOs(ctx context.Context, userID string, items []domain.Announcement) ([]ports.AnnouncementDTO, error) {
	out := make([]ports.AnnouncementDTO, 0, len(items))
	for _, item := range items {
		dto, err := s.toDTO(ctx, userID, item)
		if err != nil {
			return nil, err
		}
		out = append(out, dto)
	}
	return out, nil
}
