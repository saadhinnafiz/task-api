package main

import (
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

// temporary in-memory storage — Postgres replaces this in phase 2
var tasks = map[int]Task{
	1: {ID: 1, Title: "Set up database", Status: "pending", AssignedTo: 1},
	2: {ID: 2, Title: "Write auth middleware", Status: "in_progress", AssignedTo: 1},
}
var nextID = 3

func main() {
	app := fiber.New()
	app.Use(logger.New())

	api := app.Group("/api")

	api.Get("/tasks", func(c fiber.Ctx) error {
		list := make([]Task, 0, len(tasks))
		for _, t := range tasks {
			list = append(list, t)
		}
		return c.JSON(list)
	})

	// GET route for /tasks/:id

	api.Get("/tasks/:id", func (c fiber.Ctx) error  {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"invalid id"})
		}
		task, ok := tasks[id]
		if !ok {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "task not found"})
		}
		return c.JSON(task)
	})

	// POST route for /tasks

	api.Post("/tasks", func(c fiber.Ctx) error {
		var task Task
		
		// create an empty task box, and tell Fiber to read the incoming JSON text and put it inside that box.
		if err := c.Bind().JSON(&task); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}

		task.ID = nextID
		if task.Status == "" {
			task.Status = "pending"
		}

		tasks[nextID] = task
		nextID++

		return c.Status(fiber.StatusCreated).JSON(task)
	})


	// PUT route for /tasks/:id

	api.Put("/tasks/:id", func(c fiber.Ctx) error {
        // 1. Get the ID from the URL
        idStr := c.Params("id")
        id, err := strconv.Atoi(idStr)
        if err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
        }

        // 2. Check if the task already exists
        existingTask, exists := tasks[id]
        if !exists {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "task not found"})
        }

        // 3. Read the new data from the user
        var updatedTask Task
        if err := c.Bind().JSON(&updatedTask); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
        }

        // 4. Overwrite the old fields with the new fields
        existingTask.Title = updatedTask.Title
        existingTask.Description = updatedTask.Description
        existingTask.Status = updatedTask.Status
        existingTask.AssignedTo = updatedTask.AssignedTo

        // 5. Save it back to the map
        tasks[id] = existingTask

        return c.JSON(existingTask)
    })

	// TODO: DELETE /tasks/:id

	log.Fatal(app.Listen(":3000"))
}