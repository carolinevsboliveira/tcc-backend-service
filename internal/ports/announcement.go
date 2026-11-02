package ports

import (
	"context"
	"time"

	"github.com/caroline-oliveira/backend-tcc/internal/domain"
)

type AnnouncementCreateInput struct {
	Title        string
	Body         string
	DepartmentID string
}

type AnnouncementUpdateInput struct {
	Title        string
	Body         string
	DepartmentID string
}

type MarkReadInput struct {
	Text string
}

type HomeSummary struct {
	Unread int `json:"unread"`
}

type AnnouncementDTO struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Body         string    `json:"body"`
	CreatedAt    time.Time `json:"created_at"`
	AuthorID     string    `json:"author_id"`
	DepartmentID string    `json:"department_id"`
	Read         bool      `json:"read"`
}

type AnnouncementRepository interface {
	Create(ctx context.Context, announcement domain.Announcement) (domain.Announcement, error)
	GetByID(ctx context.Context, id string) (domain.Announcement, error)
	Update(ctx context.Context, announcement domain.Announcement) (domain.Announcement, error)
	Delete(ctx context.Context, id string) error
	ListByAuthor(ctx context.Context, authorID string) ([]domain.Announcement, error)
	ListByDepartment(ctx context.Context, departmentID string) ([]domain.Announcement, error)
	MarkRead(ctx context.Context, read domain.AnnouncementRead) error
	IsRead(ctx context.Context, announcementID, userID string) (bool, error)
}

type AnnouncementService interface {
	List(ctx context.Context, accessToken string) ([]AnnouncementDTO, error)
	Get(ctx context.Context, accessToken, id string) (AnnouncementDTO, error)
	Create(ctx context.Context, accessToken string, input AnnouncementCreateInput) (AnnouncementDTO, error)
	Update(ctx context.Context, accessToken, id string, input AnnouncementUpdateInput) (AnnouncementDTO, error)
	Delete(ctx context.Context, accessToken, id string) error
	MarkRead(ctx context.Context, accessToken, id string, input MarkReadInput) (AnnouncementDTO, error)
	Home(ctx context.Context, accessToken string) (HomeSummary, error)
}
