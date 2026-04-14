-- Remove seed data
DELETE FROM tasks WHERE title IN ('Design the homepage', 'Set up CI pipeline', 'Write documentation');
DELETE FROM projects WHERE name = 'Demo Project';
DELETE FROM users WHERE email = 'test@example.com';