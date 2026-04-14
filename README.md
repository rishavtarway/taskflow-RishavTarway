# TaskFlow - Task Management System

A full-stack task management system built with Go and React. Users can register, login, create projects, add tasks, and manage them through a clean UI.

## Overview

**Tech Stack:**
- **Backend**: Go with Chi router, PostgreSQL 16, JWT authentication
- **Frontend**: React 18 with TypeScript, React Query, Zustand, shadcn/ui
- **Database**: PostgreSQL with golang-migrate for migrations
- **Infrastructure**: Docker + Docker Compose v2

## Architecture Decisions

**Why Chi router?**
- Lightweight, fast, and idiomatic Go HTTP router
- No dependencies beyond standard library (uses Go 1.22's built-in routing)
- Clean middleware pattern with context propagation

**Why React Query over Redux?**
- Automatic request deduplication and caching
- Built-in loading/error states
- Optimistic updates made simple
- Less boilerplate than Redux

**Why golang-migrate over ORM auto-migrate?**
- Explicit, version-controlled schema changes
- Up and down migrations for safety
- No magic schema inference that can cause issues

**Tradeoffs made:**
- No refresh tokens (single JWT with 24h expiry)
- No role-based access control (owner-only project operations)
- No WebSocket for real-time updates (polling via React Query)
- Basic client-side task filtering (status/assignee)

## Running Locally

```bash
# Clone the repository
git clone https://github.com/rishavtarway/taskflow-rishavtarway.git
cd taskflow-rishavtarway

# Copy environment file
cp .env.example .env

# Start all services (PostgreSQL, migrations, API, frontend)
docker compose up --build

# Access the application
# Frontend: http://localhost:3000
# API: http://localhost:8080
```

## Running Migrations

Migrations run automatically on `docker compose up` via the migrate service. 

To run manually:
```bash
docker compose run --rm migrate
```

## Test Credentials

```
Email:    test@example.com
Password: password123
```

## API Reference

### Authentication

**POST /auth/register**
```json
// Request
{ "name": "John Doe", "email": "john@example.com", "password": "secret123" }

// Response 201
{ "token": "eyJ...", "user": { "id": "...", "name": "John Doe", "email": "john@example.com" } }
```

**POST /auth/login**
```json
// Request
{ "email": "test@example.com", "password": "password123" }

// Response 200
{ "token": "eyJ...", "user": { ... } }
```

### Projects

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /projects | List user's projects |
| POST | /projects | Create new project |
| GET | /projects/:id | Get project with tasks |
| PATCH | /projects/:id | Update project (owner only) |
| DELETE | /projects/:id | Delete project (owner only) |

### Tasks

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /projects/:id/tasks | List tasks (?status= & ?assignee=) |
| POST | /projects/:id/tasks | Create task |
| PATCH | /tasks/:id | Update task |
| DELETE | /tasks/:id | Delete task |

### Error Responses

```json
// 400 Validation
{ "error": "validation failed", "fields": { "email": "is required" } }

// 401 Unauthorized
{ "error": "unauthorized" }

// 403 Forbidden
{ "error": "forbidden" }

// 404 Not Found
{ "error": "not found" }
```

## What You'd Do With More Time

1. **Real-time updates**: Implement WebSocket for live task updates across clients
2. **Pagination**: Add cursor-based pagination to list endpoints
3. **Refresh tokens**: Add JWT refresh mechanism for longer sessions
4. **Dark mode**: Add theme toggle with localStorage persistence
5. **Tests**: Write integration tests for auth and task endpoints
6. **Drag-and-drop**: Add @dnd-kit for task reordering
7. **Project stats**: Add `/projects/:id/stats` endpoint with task analytics
8. **Rate limiting**: Add request rate limiting for API endpoints