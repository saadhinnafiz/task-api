# Task API

A REST API for creating and assigning tasks, built while learning Go and Fiber.

## Status

Work in progress, built in phases:

- [x] **Phase 1** — Fiber routes with in-memory storage
- [ ] **Phase 2** — PostgreSQL
- [ ] **Phase 3** — JWT authentication

## Endpoints

| Method | Path                 | Description       | Done |
| ------ | -------------------- | ----------------- | ---- |
| GET    | `/api/tasks`         | List all tasks    | [x]  |
| GET    | `/api/tasks/:id`     | Get one task      | [ ]  |
| POST   | `/api/tasks`         | Create a task     | [ ]  |
| PUT    | `/api/tasks/:id`     | Update a task     | [ ]  |
| DELETE | `/api/tasks/:id`     | Delete a task     | [ ]  |
| POST   | `/api/auth/register` | Create an account | [ ]  |
| POST   | `/api/auth/login`    | Get a JWT         | [ ]  |

## Data model

tasks
id, title, description, status, assigned_to, created_at

users
id, name, email, password_hash, created_at

`status` is one of `pending`, `in_progress`, `done`.

## Tech stack

- Go
- [Fiber](https://gofiber.io) v3
- PostgreSQL _(phase 2)_

## Running locally

```bash
go mod download
go run main.go
```

Server starts on `http://localhost:3000`.

```bash
curl http://localhost:3000/api/tasks
```

## Note

Storage is currently an in-memory map, so data resets on restart.
PostgreSQL replaces this in phase 2.
