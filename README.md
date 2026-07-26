# Task API

A REST API for creating and assigning tasks, built while learning Go and Fiber.

## Status

Work in progress, built in phases:

- [x] **Phase 1** — Fiber routes with in-memory storage
- [x] **Phase 2** — PostgreSQL integration
- [ ] **Phase 3** — JWT authentication

## Endpoints

| Method | Path                 | Description       | Done |
| ------ | -------------------- | ----------------- | ---- |
| GET    | `/api/tasks`         | List all tasks    | [x]  |
| GET    | `/api/tasks/:id`     | Get one task      | [x]  |
| POST   | `/api/tasks`         | Create a task     | [x]  |
| PUT    | `/api/tasks/:id`     | Update a task     | [x]  |
| DELETE | `/api/tasks/:id`     | Delete a task     | [x]  |
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
- PostgreSQL (via `pgx/v5`)
- godotenv

## Running locally

Create a `.env` file in the project root:

DATABASE_URL=postgres://your_username@localhost:5432/taskapi

Then run:

```bash
go mod download
go run .
```

Server starts on `http://localhost:3000`.

```bash
curl http://localhost:3000/api/tasks
```

## Note

Data is persisted in PostgreSQL.
