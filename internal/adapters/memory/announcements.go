package memory

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/caroline-oliveira/backend-tcc/internal/domain"
)

type AnnouncementRepository struct {
	mu     sync.RWMutex
	byID   map[string]domain.Announcement
	reads  map[string]domain.AnnouncementRead
	nextID int
}

func NewAnnouncementRepository() *AnnouncementRepository {
	repo := &AnnouncementRepository{
		byID:   make(map[string]domain.Announcement),
		reads:  make(map[string]domain.AnnouncementRead),
		nextID: 1,
	}
	for _, seed := range defaultSeedAnnouncements() {
		id := fmt.Sprintf("ann-%03d", repo.nextID)
		repo.nextID++
		seed.ID = id
		repo.byID[id] = seed
	}
	return repo
}

func defaultSeedAnnouncements() []domain.Announcement {
	now := time.Now().UTC()
	return []domain.Announcement{
		{
			Title:        "Aula cancelada",
			Body:         "Não haverá aula amanhã.",
			AuthorID:     "prf-001",
			DepartmentID: "dept-001",
			CreatedAt:    now.Add(-48 * time.Hour),
		},
		{
			Title:        "Material disponível",
			Body:         "Slides da última aula estão no portal.",
			AuthorID:     "prf-001",
			DepartmentID: "dept-001",
			CreatedAt:    now.Add(-24 * time.Hour),
		},
		{
			Title:        "Entrega antecipada",
			Body:         "O prazo do projeto foi antecipado em dois dias.",
			AuthorID:     "prf-002",
			DepartmentID: "dept-004",
			CreatedAt:    now.Add(-12 * time.Hour),
		},
	}
}

func readKey(announcementID, userID string) string {
	return announcementID + "|" + userID
}

func (r *AnnouncementRepository) Create(_ context.Context, announcement domain.Announcement) (domain.Announcement, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := fmt.Sprintf("ann-%03d", r.nextID)
	r.nextID++
	announcement.ID = id
	if announcement.CreatedAt.IsZero() {
		announcement.CreatedAt = time.Now().UTC()
	}
	r.byID[id] = announcement
	return announcement, nil
}

func (r *AnnouncementRepository) GetByID(_ context.Context, id string) (domain.Announcement, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	announcement, ok := r.byID[id]
	if !ok {
		return domain.Announcement{}, domain.ErrNotFound
	}
	return announcement, nil
}

func (r *AnnouncementRepository) Update(_ context.Context, announcement domain.Announcement) (domain.Announcement, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byID[announcement.ID]; !ok {
		return domain.Announcement{}, domain.ErrNotFound
	}
	r.byID[announcement.ID] = announcement
	return announcement, nil
}

func (r *AnnouncementRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byID[id]; !ok {
		return domain.ErrNotFound
	}
	delete(r.byID, id)

	prefix := id + "|"
	for key := range r.reads {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			delete(r.reads, key)
		}
	}
	return nil
}

func (r *AnnouncementRepository) ListByAuthor(_ context.Context, authorID string) ([]domain.Announcement, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]domain.Announcement, 0)
	for _, announcement := range r.byID {
		if announcement.AuthorID == authorID {
			out = append(out, announcement)
		}
	}
	sortAnnouncements(out)
	return out, nil
}

func (r *AnnouncementRepository) ListByDepartment(_ context.Context, departmentID string) ([]domain.Announcement, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]domain.Announcement, 0)
	for _, announcement := range r.byID {
		if announcement.DepartmentID == departmentID {
			out = append(out, announcement)
		}
	}
	sortAnnouncements(out)
	return out, nil
}

func (r *AnnouncementRepository) MarkRead(_ context.Context, read domain.AnnouncementRead) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byID[read.AnnouncementID]; !ok {
		return domain.ErrNotFound
	}
	if read.ReadAt.IsZero() {
		read.ReadAt = time.Now().UTC()
	}
	r.reads[readKey(read.AnnouncementID, read.UserID)] = read
	return nil
}

func (r *AnnouncementRepository) IsRead(_ context.Context, announcementID, userID string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, ok := r.reads[readKey(announcementID, userID)]
	return ok, nil
}

func sortAnnouncements(items []domain.Announcement) {
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
}
