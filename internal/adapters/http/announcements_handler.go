package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/caroline-oliveira/backend-tcc/internal/ports"
)

type AnnouncementsHandler struct {
	service ports.AnnouncementService
}

func NewAnnouncementsHandler(service ports.AnnouncementService) *AnnouncementsHandler {
	return &AnnouncementsHandler{service: service}
}

type announcementWriteRequest struct {
	Title        string `json:"title"`
	Body         string `json:"body"`
	DepartmentID string `json:"department_id"`
}

type markReadRequest struct {
	Text string `json:"text"`
}

func (h *AnnouncementsHandler) List(w http.ResponseWriter, r *http.Request) {
	token, ok := requireBearer(w, r)
	if !ok {
		return
	}

	items, err := h.service.List(r.Context(), token)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *AnnouncementsHandler) Get(w http.ResponseWriter, r *http.Request) {
	token, ok := requireBearer(w, r)
	if !ok {
		return
	}

	item, err := h.service.Get(r.Context(), token, r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *AnnouncementsHandler) Create(w http.ResponseWriter, r *http.Request) {
	token, ok := requireBearer(w, r)
	if !ok {
		return
	}

	var body announcementWriteRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, errorBody{Error: err.Error()})
		return
	}

	item, err := h.service.Create(r.Context(), token, ports.AnnouncementCreateInput{
		Title:        body.Title,
		Body:         body.Body,
		DepartmentID: body.DepartmentID,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *AnnouncementsHandler) Update(w http.ResponseWriter, r *http.Request) {
	token, ok := requireBearer(w, r)
	if !ok {
		return
	}

	var body announcementWriteRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, errorBody{Error: err.Error()})
		return
	}

	item, err := h.service.Update(r.Context(), token, r.PathValue("id"), ports.AnnouncementUpdateInput{
		Title:        body.Title,
		Body:         body.Body,
		DepartmentID: body.DepartmentID,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *AnnouncementsHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

func (h *AnnouncementsHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	token, ok := requireBearer(w, r)
	if !ok {
		return
	}

	var body markReadRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil && err != io.EOF {
		writeJSON(w, http.StatusBadRequest, errorBody{Error: err.Error()})
		return
	}

	item, err := h.service.MarkRead(r.Context(), token, r.PathValue("id"), ports.MarkReadInput{
		Text: body.Text,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *AnnouncementsHandler) Home(w http.ResponseWriter, r *http.Request) {
	token, ok := requireBearer(w, r)
	if !ok {
		return
	}

	summary, err := h.service.Home(r.Context(), token)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, summary)
}
