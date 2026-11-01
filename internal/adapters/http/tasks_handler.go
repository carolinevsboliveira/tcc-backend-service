package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/caroline-oliveira/backend-tcc/internal/domain"
	"github.com/caroline-oliveira/backend-tcc/internal/ports"
)

type TasksHandler struct {
	service ports.TaskService
}

func NewTasksHandler(service ports.TaskService) *TasksHandler {
	return &TasksHandler{service: service}
}

type taskWriteRequest struct {
	Title        string            `json:"title"`
	Subject      string            `json:"subject"`
	DueDate      time.Time         `json:"due_date"`
	Status       domain.TaskStatus `json:"status"`
	DepartmentID string            `json:"department_id"`
}

type taskStatusRequest struct {
	Status domain.TaskStatus `json:"status"`
}

func (h *TasksHandler) List(w http.ResponseWriter, r *http.Request) {
	token, ok := requireBearer(w, r)
	if !ok {
		return
	}

	tasks, err := h.service.List(r.Context(), token)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, tasks)
}

func (h *TasksHandler) Summary(w http.ResponseWriter, r *http.Request) {
	token, ok := requireBearer(w, r)
	if !ok {
		return
	}

	summary, err := h.service.Summary(r.Context(), token)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, summary)
}

func (h *TasksHandler) Get(w http.ResponseWriter, r *http.Request) {
	token, ok := requireBearer(w, r)
	if !ok {
		return
	}

	task, err := h.service.Get(r.Context(), token, r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *TasksHandler) Create(w http.ResponseWriter, r *http.Request) {
	token, ok := requireBearer(w, r)
	if !ok {
		return
	}

	var body taskWriteRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, errorBody{Error: err.Error()})
		return
	}

	task, err := h.service.Create(r.Context(), token, ports.TaskCreateInput{
		Title:        body.Title,
		Subject:      body.Subject,
		DueDate:      body.DueDate,
		Status:       body.Status,
		DepartmentID: body.DepartmentID,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, task)
}

func (h *TasksHandler) Update(w http.ResponseWriter, r *http.Request) {
	token, ok := requireBearer(w, r)
	if !ok {
		return
	}

	var body taskWriteRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, errorBody{Error: err.Error()})
		return
	}

	task, err := h.service.Update(r.Context(), token, r.PathValue("id"), ports.TaskUpdateInput{
		Title:        body.Title,
		Subject:      body.Subject,
		DueDate:      body.DueDate,
		Status:       body.Status,
		DepartmentID: body.DepartmentID,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *TasksHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	token, ok := requireBearer(w, r)
	if !ok {
		return
	}

	var body taskStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, errorBody{Error: err.Error()})
		return
	}

	task, err := h.service.UpdateStatus(r.Context(), token, r.PathValue("id"), body.Status)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *TasksHandler) Delete(w http.ResponseWriter, r *http.Request) {
	token, ok := requireBearer(w, r)
	if !ok {
		return
	}

	if err := h.service.Delete(r.Context(), token, r.PathValue("id")); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusNoContent, nil)
}

func requireBearer(w http.ResponseWriter, r *http.Request) (string, bool) {
	token := extractBearerToken(r.Header.Get("Authorization"))
	if token == "" {
		writeJSON(w, http.StatusUnauthorized, errorBody{Error: "Unauthorized"})
		return "", false
	}
	return token, true
}
