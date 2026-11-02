package domain

import "time"

type Announcement struct {
	ID           string
	Title        string
	Body         string
	AuthorID     string
	DepartmentID string
	CreatedAt    time.Time
}

type AnnouncementRead struct {
	AnnouncementID string
	UserID         string
	Text           string
	ReadAt         time.Time
}
