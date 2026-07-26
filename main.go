package main

import (
	"context"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	AssignedTo  int    `json:"assigned_to"`
}

func main() {
	connectDB()
	defer db.Close()

	app := fiber.New()
	app.Use(logger.New())

	api := app.Group("/api")

	// GET /tasks — Fetch all from Postgres
	api.Get("/tasks", func(c fiber.Ctx) error {
		rows, err := db.Query(context.Background(), "SELECT id, title, description, status, assigned_to FROM tasks")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		defer rows.Close()

		var list []Task
		for rows.Next() {
			var t Task
			if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.AssignedTo); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
			list = append(list, t)
		}

		if list == nil {
			list = []Task{} // return empty array instead of null if table is empty
		}

		return c.JSON(list)
	})

	// GET /tasks/:id — Fetch a single task by ID
	api.Get("/tasks/:id", func(c fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		var t Task
		query := "SELECT id, title, description, status, assigned_to FROM tasks WHERE id = $1"
		err = db.QueryRow(context.Background(), query, id).Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.AssignedTo)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "task not found"})
		}

		return c.JSON(t)
	})

	// POST /tasks — Insert into Postgres
	api.Post("/tasks", func(c fiber.Ctx) error {
		var task Task
		if err := c.Bind().JSON(&task); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}

		if task.Status == "" {
			task.Status = "pending"
		}

		query := `INSERT INTO tasks (title, description, status, assigned_to) 
		          VALUES ($1, $2, $3, $4) 
		          RETURNING id, title, description, status, assigned_to`

		var created Task
		err := db.QueryRow(context.Background(), query, task.Title, task.Description, task.Status, task.AssignedTo).
			Scan(&created.ID, &created.Title, &created.Description, &created.Status, &created.AssignedTo)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(created)
	})

	// PUT /tasks/:id — Update in Postgres
	api.Put("/tasks/:id", func(c fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		var updated Task
		if err := c.Bind().JSON(&updated); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}

		query := `UPDATE tasks SET title = $1, description = $2, status = $3, assigned_to = $4 
		          WHERE id = $5 
		          RETURNING id, title, description, status, assigned_to`

		var saved Task
		err = db.QueryRow(context.Background(), query, updated.Title, updated.Description, updated.Status, updated.AssignedTo, id).
			Scan(&saved.ID, &saved.Title, &saved.Description, &saved.Status, &saved.AssignedTo)

		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "task not found or update failed"})
		}

		return c.JSON(saved)
	})

	// DELETE /tasks/:id — Delete from Postgres
	api.Delete("/tasks/:id", func(c fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		result, err := db.Exec(context.Background(), "DELETE FROM tasks WHERE id = $1", id)
		if err != nil || result.RowsAffected() == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "task not found"})
		}

		return c.SendStatus(fiber.StatusNoContent)
	})

	log.Fatal(app.Listen(":3000"))
}