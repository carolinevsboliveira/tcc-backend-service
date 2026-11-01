package domain

import "time"

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pendente"
	TaskStatusInProgress TaskStatus = "em_andamento"
	TaskStatusDone       TaskStatus = "concluida"
)

type Task struct {
	ID           string
	Title        string
	Subject      string
	DueDate      time.Time
	Status       TaskStatus
	AuthorID     string
	DepartmentID string
}

func IsValidTaskStatus(status TaskStatus) bool {
	switch status {
	case TaskStatusPending, TaskStatusInProgress, TaskStatusDone:
		return true
	default:
		return false
	}
}
