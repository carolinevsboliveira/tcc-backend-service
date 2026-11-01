package memory

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/caroline-oliveira/backend-tcc/internal/domain"
)

type TaskRepository struct {
	mu     sync.RWMutex
	byID   map[string]domain.Task
	nextID int
}

func NewTaskRepository() *TaskRepository {
	repo := &TaskRepository{
		byID:   make(map[string]domain.Task),
		nextID: 1,
	}
	for _, seed := range defaultSeedTasks() {
		id := fmt.Sprintf("task-%03d", repo.nextID)
		repo.nextID++
		seed.ID = id
		repo.byID[id] = seed
	}
	return repo
}

func defaultSeedTasks() []domain.Task {
	now := time.Now().UTC()
	return []domain.Task{
		{
			Title:        "Lista de derivadas",
			Subject:      "Cálculo I",
			DueDate:      now.Add(7 * 24 * time.Hour),
			Status:       domain.TaskStatusPending,
			AuthorID:     "prf-001",
			DepartmentID: "dept-001",
		},
		{
			Title:        "Prova parcial",
			Subject:      "Álgebra Linear",
			DueDate:      now.Add(14 * 24 * time.Hour),
			Status:       domain.TaskStatusInProgress,
			AuthorID:     "prf-001",
			DepartmentID: "dept-001",
		},
		{
			Title:        "Leitura capítulo 3",
			Subject:      "Programação Web",
			DueDate:      now.Add(5 * 24 * time.Hour),
			Status:       domain.TaskStatusPending,
			AuthorID:     "prf-002",
			DepartmentID: "dept-004",
		},
		{
			Title:        "Projeto final",
			Subject:      "Engenharia de Software",
			DueDate:      now.Add(21 * 24 * time.Hour),
			Status:       domain.TaskStatusDone,
			AuthorID:     "prf-002",
			DepartmentID: "dept-004",
		},
		{
			Title:        "Relatório de laboratório",
			Subject:      "Biologia Celular",
			DueDate:      now.Add(3 * 24 * time.Hour),
			Status:       domain.TaskStatusPending,
			AuthorID:     "prf-003",
			DepartmentID: "dept-003",
		},
		{
			Title:        "Quiz de anatomia",
			Subject:      "Anatomia",
			DueDate:      now.Add(10 * 24 * time.Hour),
			Status:       domain.TaskStatusInProgress,
			AuthorID:     "prf-003",
			DepartmentID: "dept-003",
		},
	}
}

func (r *TaskRepository) Create(_ context.Context, task domain.Task) (domain.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := fmt.Sprintf("task-%03d", r.nextID)
	r.nextID++
	task.ID = id
	r.byID[id] = task
	return task, nil
}

func (r *TaskRepository) GetByID(_ context.Context, id string) (domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, ok := r.byID[id]
	if !ok {
		return domain.Task{}, domain.ErrNotFound
	}
	return task, nil
}

func (r *TaskRepository) Update(_ context.Context, task domain.Task) (domain.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byID[task.ID]; !ok {
		return domain.Task{}, domain.ErrNotFound
	}
	r.byID[task.ID] = task
	return task, nil
}

func (r *TaskRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byID[id]; !ok {
		return domain.ErrNotFound
	}
	delete(r.byID, id)
	return nil
}

func (r *TaskRepository) ListByAuthor(_ context.Context, authorID string) ([]domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]domain.Task, 0)
	for _, task := range r.byID {
		if task.AuthorID == authorID {
			out = append(out, task)
		}
	}
	sortTasks(out)
	return out, nil
}

func (r *TaskRepository) ListByDepartment(_ context.Context, departmentID string) ([]domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]domain.Task, 0)
	for _, task := range r.byID {
		if task.DepartmentID == departmentID {
			out = append(out, task)
		}
	}
	sortTasks(out)
	return out, nil
}

func sortTasks(tasks []domain.Task) {
	sort.Slice(tasks, func(i, j int) bool { return tasks[i].ID < tasks[j].ID })
}
