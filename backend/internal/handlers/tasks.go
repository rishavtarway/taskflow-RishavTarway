package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/rishavtarway/taskflow/internal/errors"
	"github.com/rishavtarway/taskflow/internal/middleware"
	"github.com/rishavtarway/taskflow/internal/models"
)

type TaskHandler struct {
	taskModel    *models.TaskModel
	projectModel *models.ProjectModel
}

func NewTaskHandler(taskModel *models.TaskModel, projectModel *models.ProjectModel) *TaskHandler {
	return &TaskHandler{
		taskModel:    taskModel,
		projectModel: projectModel,
	}
}

type CreateTaskRequest struct {
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Priority    string     `json:"priority"`
	AssigneeID  *uuid.UUID `json:"assignee_id"`
	DueDate     *string    `json:"due_date"`
}

type UpdateTaskRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Status      *string    `json:"status"`
	Priority    *string    `json:"priority"`
	AssigneeID  *uuid.UUID `json:"assignee_id"`
	DueDate     *string    `json:"due_date"`
}

var validStatuses = map[string]bool{"todo": true, "in_progress": true, "done": true}
var validPriorities = map[string]bool{"low": true, "medium": true, "high": true}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	projectIDStr := r.Context().Value("projectID").(string)
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		errors.WriteError(w, errors.BadRequest("invalid project id"))
		return
	}

	project, err := h.projectModel.GetByID(r.Context(), projectID)
	if err != nil || project == nil {
		errors.WriteError(w, errors.NotFound())
		return
	}

	status := r.URL.Query().Get("status")
	assignee := r.URL.Query().Get("assignee")

	tasks, err := h.taskModel.GetByProjectIDWithFilters(r.Context(), projectID, status, assignee)
	if err != nil {
		errors.WriteError(w, errors.InternalServerError())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string][]models.Task{"tasks": tasks})
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	projectIDStr := r.Context().Value("projectID").(string)
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		errors.WriteError(w, errors.BadRequest("invalid project id"))
		return
	}

	project, err := h.projectModel.GetByID(r.Context(), projectID)
	if err != nil || project == nil {
		errors.WriteError(w, errors.NotFound())
		return
	}

	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteError(w, errors.BadRequest("invalid request body"))
		return
	}

	if req.Title == "" {
		errors.WriteErrorWithFields(w, errors.BadRequestValidation(map[string]string{"title": "is required"}))
		return
	}

	status := "todo"
	if req.Priority == "" {
		req.Priority = "medium"
	}

	userID := middleware.GetUserID(r.Context())
	task, err := h.taskModel.Create(r.Context(), req.Title, req.Description, status, req.Priority, projectID, req.AssigneeID, userID, req.DueDate)
	if err != nil {
		errors.WriteError(w, errors.InternalServerError())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	taskID, err := uuid.Parse(r.Context().Value("taskID").(string))
	if err != nil {
		errors.WriteError(w, errors.BadRequest("invalid task id"))
		return
	}

	task, err := h.taskModel.GetByID(r.Context(), taskID)
	if err != nil || task == nil {
		errors.WriteError(w, errors.NotFound())
		return
	}

	project, err := h.projectModel.GetByID(r.Context(), task.ProjectID)
	if err != nil || project == nil {
		errors.WriteError(w, errors.NotFound())
		return
	}

	userID := middleware.GetUserID(r.Context())
	if project.OwnerID != userID {
		hasCreatorAccess := task.CreatorID != nil && *task.CreatorID == userID
		if !hasCreatorAccess {
			errors.WriteError(w, errors.Forbidden())
			return
		}
	}

	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteError(w, errors.BadRequest("invalid request body"))
		return
	}

	if req.Status != nil && !validStatuses[*req.Status] {
		errors.WriteErrorWithFields(w, errors.BadRequestValidation(map[string]string{"status": "must be todo, in_progress, or done"}))
		return
	}

	if req.Priority != nil && !validPriorities[*req.Priority] {
		errors.WriteErrorWithFields(w, errors.BadRequestValidation(map[string]string{"priority": "must be low, medium, or high"}))
		return
	}

	updatedTask, err := h.taskModel.Update(r.Context(), taskID, req.Title, req.Description, req.Status, req.Priority, req.AssigneeID, req.DueDate)
	if err != nil {
		errors.WriteError(w, errors.InternalServerError())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedTask)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	taskID, err := uuid.Parse(r.Context().Value("taskID").(string))
	if err != nil {
		errors.WriteError(w, errors.BadRequest("invalid task id"))
		return
	}

	task, err := h.taskModel.GetByID(r.Context(), taskID)
	if err != nil || task == nil {
		errors.WriteError(w, errors.NotFound())
		return
	}

	project, err := h.projectModel.GetByID(r.Context(), task.ProjectID)
	if err != nil || project == nil {
		errors.WriteError(w, errors.NotFound())
		return
	}

	userID := middleware.GetUserID(r.Context())
	if project.OwnerID != userID {
		hasCreatorAccess := task.CreatorID != nil && *task.CreatorID == userID
		if !hasCreatorAccess {
			errors.WriteError(w, errors.Forbidden())
			return
		}
	}

	if err := h.taskModel.Delete(r.Context(), taskID); err != nil {
		errors.WriteError(w, errors.InternalServerError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
