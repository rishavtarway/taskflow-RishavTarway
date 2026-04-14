package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	OwnerID     uuid.UUID `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
	Tasks       []Task    `json:"tasks,omitempty"`
}

type ProjectModel struct {
	db *sql.DB
}

func NewProjectModel(db *sql.DB) *ProjectModel {
	return &ProjectModel{db: db}
}

func (m *ProjectModel) Create(ctx context.Context, name string, description *string, ownerID uuid.UUID) (*Project, error) {
	var project Project
	err := m.db.QueryRowContext(ctx,
		`INSERT INTO projects (name, description, owner_id) VALUES ($1, $2, $3) 
		 RETURNING id, name, description, owner_id, created_at`,
		name, description, ownerID,
	).Scan(&project.ID, &project.Name, &project.Description, &project.OwnerID, &project.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (m *ProjectModel) GetByID(ctx context.Context, id uuid.UUID) (*Project, error) {
	var project Project
	var description sql.NullString
	err := m.db.QueryRowContext(ctx,
		`SELECT id, name, description, owner_id, created_at FROM projects WHERE id = $1`,
		id,
	).Scan(&project.ID, &project.Name, &description, &project.OwnerID, &project.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if description.Valid {
		project.Description = &description.String
	}
	return &project, nil
}

func (m *ProjectModel) GetByIDWithTasks(ctx context.Context, id uuid.UUID) (*Project, error) {
	project, err := m.GetByID(ctx, id)
	if err != nil || project == nil {
		return project, err
	}

	tasks, err := NewTaskModel(m.db).GetByProjectID(ctx, id)
	if err != nil {
		return nil, err
	}
	project.Tasks = tasks

	return project, nil
}

func (m *ProjectModel) ListForUser(ctx context.Context, userID uuid.UUID) ([]Project, error) {
	rows, err := m.db.QueryContext(ctx,
		`SELECT DISTINCT p.id, p.name, p.description, p.owner_id, p.created_at 
		 FROM projects p
		 LEFT JOIN tasks t ON t.project_id = p.id
		 WHERE p.owner_id = $1 OR t.assignee_id = $1
		 ORDER BY p.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		var description sql.NullString
		if err := rows.Scan(&p.ID, &p.Name, &description, &p.OwnerID, &p.CreatedAt); err != nil {
			return nil, err
		}
		if description.Valid {
			p.Description = &description.String
		}
		projects = append(projects, p)
	}
	return projects, nil
}

func (m *ProjectModel) Update(ctx context.Context, id uuid.UUID, name, description *string) (*Project, error) {
	var project Project
	var desc sql.NullString
	err := m.db.QueryRowContext(ctx,
		`UPDATE projects SET name = COALESCE($1, name), description = COALESCE($2, description) 
		 WHERE id = $3 RETURNING id, name, description, owner_id, created_at`,
		name, description, id,
	).Scan(&project.ID, &project.Name, &desc, &project.OwnerID, &project.CreatedAt)
	if err != nil {
		return nil, err
	}
	if desc.Valid {
		project.Description = &desc.String
	}
	return &project, nil
}

func (m *ProjectModel) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := m.db.ExecContext(ctx, "DELETE FROM projects WHERE id = $1", id)
	return err
}
