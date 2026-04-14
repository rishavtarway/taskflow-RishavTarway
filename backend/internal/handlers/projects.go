package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/rishavtarway/taskflow/internal/errors"
	"github.com/rishavtarway/taskflow/internal/middleware"
	"github.com/rishavtarway/taskflow/internal/models"
)

type ProjectHandler struct {
	projectModel *models.ProjectModel
}

func NewProjectHandler(projectModel *models.ProjectModel) *ProjectHandler {
	return &ProjectHandler{projectModel: projectModel}
}

type CreateProjectRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type UpdateProjectRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

func (h *ProjectHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.GetUserID(r.Context())
	projects, err := h.projectModel.ListForUser(r.Context(), userID)
	if err != nil {
		errors.WriteError(w, errors.InternalServerError())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string][]models.Project{"projects": projects})
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteError(w, errors.BadRequest("invalid request body"))
		return
	}

	if req.Name == "" {
		errors.WriteErrorWithFields(w, errors.BadRequestValidation(map[string]string{"name": "is required"}))
		return
	}

	userID := middleware.GetUserID(r.Context())
	project, err := h.projectModel.Create(r.Context(), req.Name, req.Description, userID)
	if err != nil {
		errors.WriteError(w, errors.InternalServerError())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(project)
}

func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		errors.WriteError(w, errors.BadRequest("invalid project id"))
		return
	}

	project, err := h.projectModel.GetByIDWithTasks(r.Context(), projectID)
	if err != nil {
		errors.WriteError(w, errors.InternalServerError())
		return
	}
	if project == nil {
		errors.WriteError(w, errors.NotFound())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(project)
}

func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		errors.WriteError(w, errors.BadRequest("invalid project id"))
		return
	}

	project, err := h.projectModel.GetByID(r.Context(), projectID)
	if err != nil {
		errors.WriteError(w, errors.InternalServerError())
		return
	}
	if project == nil {
		errors.WriteError(w, errors.NotFound())
		return
	}

	userID := middleware.GetUserID(r.Context())
	if project.OwnerID != userID {
		errors.WriteError(w, errors.Forbidden())
		return
	}

	var req UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteError(w, errors.BadRequest("invalid request body"))
		return
	}

	updatedProject, err := h.projectModel.Update(r.Context(), projectID, req.Name, req.Description)
	if err != nil {
		errors.WriteError(w, errors.InternalServerError())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedProject)
}

func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		errors.WriteError(w, errors.BadRequest("invalid project id"))
		return
	}

	project, err := h.projectModel.GetByID(r.Context(), projectID)
	if err != nil {
		errors.WriteError(w, errors.InternalServerError())
		return
	}
	if project == nil {
		errors.WriteError(w, errors.NotFound())
		return
	}

	userID := middleware.GetUserID(r.Context())
	if project.OwnerID != userID {
		errors.WriteError(w, errors.Forbidden())
		return
	}

	if err := h.projectModel.Delete(r.Context(), projectID); err != nil {
		errors.WriteError(w, errors.InternalServerError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
