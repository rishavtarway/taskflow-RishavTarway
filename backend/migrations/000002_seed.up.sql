-- Seed data for testing
-- Password: password123 (bcrypt cost 12)

INSERT INTO users (id, name, email, password) VALUES
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Test User', 'test@example.com', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYL8J8U7K4i');

INSERT INTO projects (id, name, description, owner_id) VALUES
('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'Demo Project', 'A demo project for testing', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11');

INSERT INTO tasks (title, description, status, priority, project_id, assignee_id, due_date) VALUES
('Setup project structure', 'Initialize the project with proper folders', 'done', 'high', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2026-04-15'),
('Design database schema', 'Create tables and indexes for users, projects, tasks', 'in_progress', 'high', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2026-04-20'),
('Build authentication API', 'Implement register and login endpoints with JWT', 'todo', 'medium', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2026-04-25');