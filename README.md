# TaskFlow

TaskFlow is a task management system built as a take-home assignment. Users can register, log in, create projects, and manage tasks with status and priority tracking.

Stack: Go (chi router), PostgreSQL 16, React 18 + TypeScript, shadcn/ui, Docker.

## Architecture Decisions

### Why did you structure things the way you did?

I split the Go backend into handlers/models/config rather than putting everything in main.go because handlers grow fast and mixing DB logic with HTTP logic makes both harder to test. The models layer handles all database queries, keeping handlers thin (~20 lines each) and focused only on parsing input and writing responses.

For the frontend, I used React Query for server state because this app's state is almost entirely server state - projects and tasks come from the API. Redux would add boilerplate with no benefit at this scale. I used Zustand only for the auth token which is genuinely client state.

### What tradeoffs did you make?

PATCH uses pointer fields (*string) to distinguish 'field not provided' from 'field set to empty string'. This adds verbosity in the handler (need to check for nil before including in query) but prevents accidental field zeroing when a client sends a partial update.

I chose chi over gin because its middleware composition is cleaner for this size of API - gin's context type adds overhead I didn't need here.

### What did you intentionally leave out and why?

No refresh tokens - the spec says 24h JWT which is sufficient for an internal tool. In production I'd use rotating refresh tokens stored in httpOnly cookies with a Redis revocation list.

No rate limiting - would add golang.org/x/time/rate per-IP middleware in production, but not needed for a take-home assignment.

No role-based permissions beyond owner/non-owner - the spec only requires owner-level project control, so I kept it simple.

The tasks table has a creator_id field to support the "project owner OR task creator" delete permission specified in the requirements.

## Running Locally

```bash
git clone https://github.com/rishavtarway/taskflow-RishavTarway.git
cd taskflow-RishavTarway
cp .env.example .env
docker compose up
```

App available at http://localhost:3000
API available at http://localhost:8080

## Running Migrations

Migrations run automatically when you run `docker compose up`. No manual steps needed.

## Test Credentials

Email:    test@example.com
Password: password123

## API Reference

### POST /auth/register
```bash
curl -s -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","password":"secret123"}'
```
Response 201:
```json
{
  "token": "eyJ...",
  "user": { "id": "uuid", "name": "John Doe", "email": "john@example.com" }
}
```

### POST /auth/login
```bash
curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```
Response 200:
```json
{
  "token": "eyJ...",
  "user": { "id": "uuid", "name": "Test User", "email": "test@example.com" }
}
```

### GET /projects
Requires `Authorization: Bearer <token>`
```bash
curl -s http://localhost:8080/projects \
  -H "Authorization: Bearer <token>"
```
Response 200:
```json
{ "projects": [...] }
```

### POST /projects
Requires `Authorization: Bearer <token>`
```bash
curl -s -X POST http://localhost:8080/projects \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"New Project","description":"Optional description"}'
```
Response 201: Returns created project object

### GET /projects/:id
Requires `Authorization: Bearer <token>`
Response 200: Returns project with tasks array

### PATCH /projects/:id
Requires `Authorization: Bearer <token>` (owner only)
Response 200: Returns updated project
Response 403: If not owner

### DELETE /projects/:id
Requires `Authorization: Bearer <token>` (owner only)
Response 204: No content

### GET /projects/:id/tasks
Requires `Authorization: Bearer <token>`
Query params: ?status=todo&assignee=uuid
```bash
curl -s "http://localhost:8080/projects/:id/tasks?status=todo" \
  -H "Authorization: Bearer <token>"
```
Response 200:
```json
{ "tasks": [...] }
```

### POST /projects/:id/tasks
Requires `Authorization: Bearer <token>`
```bash
curl -s -X POST http://localhost:8080/projects/:id/tasks \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"title":"New task","priority":"high"}'
```
Response 201: Returns created task

### PATCH /tasks/:id
Requires `Authorization: Bearer <token>`
```bash
curl -s -X PATCH http://localhost:8080/tasks/:id \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"status":"done"}'
```
Response 200: Returns updated task

### DELETE /tasks/:id
Requires `Authorization: Bearer <token>` (project owner or task creator)
Response 204: No content
Response 403: If not authorized

### Error Responses

```json
// 400 Validation error
{ "error": "validation failed", "fields": { "email": "is required" } }

// 401 Unauthorized (no/invalid token)
{ "error": "unauthorized" }

// 403 Forbidden (authenticated but not allowed)
{ "error": "forbidden" }

// 404 Not found
{ "error": "not found" }

// 409 Conflict (duplicate email)
{ "error": "email already exists" }
```

## What You'd Do With More Time

1. Replace offset-based pagination with cursor/keyset pagination — offset degrades at high page numbers with large tables. Use LIMIT with a last-seen-id for efficient navigation.

2. Add integration tests using testcontainers-go with a real Postgres instance, covering the auth flow (register → login → get JWT → call protected endpoint) and task CRUD operations.

3. Push real-time task updates via Server-Sent Events (SSE) — simpler than WebSockets for this use case, backend emits events on task mutations, frontend subscribes per-project.

4. Add rotating refresh tokens stored in httpOnly cookies with Redis-backed revocation — current 24h JWT is fine for an assignment but production needs rotation.

5. The optimistic update rollback currently uses React Query's onError context — I'd move to a proper mutation queue to handle concurrent edits to the same task gracefully.

6. Rate limit the /auth/* endpoints — 10 requests per minute per IP using golang.org/x/time/rate to prevent brute force attacks.