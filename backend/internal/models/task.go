package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	ProjectID   uuid.UUID  `json:"project_id"`
	AssigneeID  *uuid.UUID `json:"assignee_id"`
	CreatorID   *uuid.UUID `json:"creator_id"`
	DueDate     *string    `json:"due_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type TaskModel struct {
	db *sql.DB
}

func NewTaskModel(db *sql.DB) *TaskModel {
	return &TaskModel{db: db}
}

func (m *TaskModel) Create(ctx context.Context, title string, description *string, status, priority string, projectID uuid.UUID, assigneeID *uuid.UUID, creatorID uuid.UUID, dueDate *string) (*Task, error) {
	var task Task
	var desc, assignee, creator, due sql.NullString

	err := m.db.QueryRowContext(ctx,
		`INSERT INTO tasks (title, description, status, priority, project_id, assignee_id, creator_id, due_date) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		 RETURNING id, title, description, status, priority, project_id, assignee_id, creator_id, due_date, created_at, updated_at`,
		title, description, status, priority, projectID, assigneeID, creatorID, dueDate,
	).Scan(&task.ID, &task.Title, &desc, &task.Status, &task.Priority, &task.ProjectID, &assignee, &creator, &due, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if desc.Valid {
		task.Description = &desc.String
	}
	if assignee.Valid {
		task.AssigneeID = new(uuid.UUID)
		*task.AssigneeID, _ = uuid.Parse(assignee.String)
	}
	if creator.Valid {
		task.CreatorID = new(uuid.UUID)
		*task.CreatorID, _ = uuid.Parse(creator.String)
	}
	if due.Valid {
		task.DueDate = &due.String
	}
	return &task, nil
}

func (m *TaskModel) GetByID(ctx context.Context, id uuid.UUID) (*Task, error) {
	var task Task
	var desc, assignee, creator, due sql.NullString

	err := m.db.QueryRowContext(ctx,
		`SELECT id, title, description, status, priority, project_id, assignee_id, creator_id, due_date, created_at, updated_at 
		 FROM tasks WHERE id = $1`,
		id,
	).Scan(&task.ID, &task.Title, &desc, &task.Status, &task.Priority, &task.ProjectID, &assignee, &creator, &due, &task.CreatedAt, &task.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if desc.Valid {
		task.Description = &desc.String
	}
	if assignee.Valid {
		task.AssigneeID = new(uuid.UUID)
		*task.AssigneeID, _ = uuid.Parse(assignee.String)
	}
	if creator.Valid {
		task.CreatorID = new(uuid.UUID)
		*task.CreatorID, _ = uuid.Parse(creator.String)
	}
	if due.Valid {
		task.DueDate = &due.String
	}
	return &task, nil
}

func (m *TaskModel) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]Task, error) {
	rows, err := m.db.QueryContext(ctx,
		`SELECT id, title, description, status, priority, project_id, assignee_id, creator_id, due_date, created_at, updated_at 
		 FROM tasks WHERE project_id = $1 ORDER BY created_at DESC`,
		projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		var desc, assignee, creator, due sql.NullString
		if err := rows.Scan(&t.ID, &t.Title, &desc, &t.Status, &t.Priority, &t.ProjectID, &assignee, &creator, &due, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		if desc.Valid {
			t.Description = &desc.String
		}
		if assignee.Valid {
			t.AssigneeID = new(uuid.UUID)
			*t.AssigneeID, _ = uuid.Parse(assignee.String)
		}
		if creator.Valid {
			t.CreatorID = new(uuid.UUID)
			*t.CreatorID, _ = uuid.Parse(creator.String)
		}
		if due.Valid {
			t.DueDate = &due.String
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (m *TaskModel) GetByProjectIDWithFilters(ctx context.Context, projectID uuid.UUID, status, assigneeID string) ([]Task, error) {
	query := `SELECT id, title, description, status, priority, project_id, assignee_id, creator_id, due_date, created_at, updated_at 
		 FROM tasks WHERE project_id = $1`
	args := []interface{}{projectID}
	argNum := 1

	if status != "" {
		argNum++
		query += " AND status = $2"
		args = append(args, status)
	}
	if assigneeID != "" {
		argNum++
		if status != "" {
			query += " AND assignee_id = $3"
		} else {
			query += " AND assignee_id = $2"
		}
		args = append(args, assigneeID)
	}

	query += " ORDER BY created_at DESC"

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		var desc, assignee, creator, due sql.NullString
		if err := rows.Scan(&t.ID, &t.Title, &desc, &t.Status, &t.Priority, &t.ProjectID, &assignee, &creator, &due, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		if desc.Valid {
			t.Description = &desc.String
		}
		if assignee.Valid {
			t.AssigneeID = new(uuid.UUID)
			*t.AssigneeID, _ = uuid.Parse(assignee.String)
		}
		if creator.Valid {
			t.CreatorID = new(uuid.UUID)
			*t.CreatorID, _ = uuid.Parse(creator.String)
		}
		if due.Valid {
			t.DueDate = &due.String
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (m *TaskModel) Update(ctx context.Context, id uuid.UUID, title, description, status, priority *string, assigneeID *uuid.UUID, dueDate *string) (*Task, error) {
	var task Task
	var desc, assignee, creator, due sql.NullString

	err := m.db.QueryRowContext(ctx,
		`UPDATE tasks SET 
		 title = COALESCE(NULLIF($1, ''), title),
		 description = $2,
		 status = COALESCE(NULLIF($3, ''), status),
		 priority = COALESCE(NULLIF($4, ''), priority),
		 assignee_id = $5,
		 due_date = $6,
		 updated_at = NOW()
		 WHERE id = $7
		 RETURNING id, title, description, status, priority, project_id, assignee_id, creator_id, due_date, created_at, updated_at`,
		title, description, status, priority, assigneeID, dueDate, id,
	).Scan(&task.ID, &task.Title, &desc, &task.Status, &task.Priority, &task.ProjectID, &assignee, &creator, &due, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if desc.Valid {
		task.Description = &desc.String
	}
	if assignee.Valid {
		task.AssigneeID = new(uuid.UUID)
		*task.AssigneeID, _ = uuid.Parse(assignee.String)
	}
	if creator.Valid {
		task.CreatorID = new(uuid.UUID)
		*task.CreatorID, _ = uuid.Parse(creator.String)
	}
	if due.Valid {
		task.DueDate = &due.String
	}
	return &task, nil
}

func (m *TaskModel) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := m.db.ExecContext(ctx, "DELETE FROM tasks WHERE id = $1", id)
	return err
}
