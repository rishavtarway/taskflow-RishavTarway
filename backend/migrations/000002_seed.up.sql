-- Seed data for testing
-- Password: password123 (bcrypt cost 12)

INSERT INTO users (id, name, email, password) VALUES
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Test User', 'test@example.com', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYL8J8U7K4i');

INSERT INTO projects (id, name, description, owner_id) VALUES
('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'Demo Project', 'A sample project with tasks', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11');

INSERT INTO tasks (title, status, priority, project_id, assignee_id, creator_id) VALUES
('Design the homepage', 'todo', 'high', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'),
('Set up CI pipeline', 'in_progress', 'medium', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'),
('Write documentation', 'done', 'low', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11');