package http

import (
	"net/http"
	"time"
)

type Dependencies struct {
	Auth          *AuthHandler
	Tasks         *TasksHandler
	Announcements *AnnouncementsHandler
}

// NewRouter builds the HTTP mux for the API.
func NewRouter(deps Dependencies) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
			"time":   time.Now().UTC().Format(time.RFC3339),
		})
	})

	mux.HandleFunc("POST /auth/login", deps.Auth.Login)
	mux.HandleFunc("GET /auth/validate", deps.Auth.ValidateToken)
	mux.HandleFunc("GET /auth/me", deps.Auth.Me)

	mux.HandleFunc("GET /tasks", deps.Tasks.List)
	mux.HandleFunc("GET /tasks/summary", deps.Tasks.Summary)
	mux.HandleFunc("POST /tasks", deps.Tasks.Create)
	mux.HandleFunc("GET /tasks/{id}", deps.Tasks.Get)
	mux.HandleFunc("PUT /tasks/{id}", deps.Tasks.Update)
	mux.HandleFunc("PATCH /tasks/{id}/status", deps.Tasks.UpdateStatus)
	mux.HandleFunc("DELETE /tasks/{id}", deps.Tasks.Delete)

	mux.HandleFunc("GET /announcements", deps.Announcements.List)
	mux.HandleFunc("POST /announcements", deps.Announcements.Create)
	mux.HandleFunc("GET /announcements/{id}", deps.Announcements.Get)
	mux.HandleFunc("PUT /announcements/{id}", deps.Announcements.Update)
	mux.HandleFunc("DELETE /announcements/{id}", deps.Announcements.Delete)
	mux.HandleFunc("POST /announcements/{id}/read", deps.Announcements.MarkRead)
	mux.HandleFunc("GET /home", deps.Announcements.Home)

	return mux
}
