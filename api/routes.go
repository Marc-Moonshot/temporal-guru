package api

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Marc-Moonshot/temporal-guru/cache"
	"github.com/Marc-Moonshot/temporal-guru/scheduler"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
)

// Handles API calls routed from nginx
func RegisterRoutes(app *fiber.App, conn *pgx.Conn) {

	app.Get("/nrw/yearly", func(c fiber.Ctx) error {
		fmt.Println("GET /nrw/yearly")
		entry, err := cache.Get(conn, "nrw/yearly")

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				data, err := scheduler.Call("http://localhost/api/nrw/yearly", []string{"month=2025-11", "device=40961"}, conn)

				if err != nil {
					fmt.Printf("error: %v", err)
				}

				formatted, _ := json.MarshalIndent(data, "", "  ")
				fmt.Printf("Data:\n%s\n", string(formatted))
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "no cached data found",
				})
			}
			fmt.Printf("error: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
		return c.Status(fiber.StatusOK).JSON(entry)
	})

	app.Get("/nrw/monthly", func(c fiber.Ctx) error {
		fmt.Println("GET /nrw/monthly")
		entry, err := cache.Get(conn, "nrw/monthly")

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				data, err := scheduler.Call("http://localhost/api/nrw/monthly", []string{"month=2025-11", "device=40961"}, conn)

				if err != nil {
					fmt.Printf("error: %v", err)
				}

				formatted, _ := json.MarshalIndent(data, "", "  ")
				fmt.Printf("Data:\n%s\n", string(formatted))
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "no cached data found",
				})
			}
			fmt.Printf("error: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
		return c.Status(fiber.StatusOK).JSON(entry)
	})

	app.Get("/nrw/daily", func(c fiber.Ctx) error {
		fmt.Println("GET /nrw/daily")
		entry, err := cache.Get(conn, "nrw/daily")

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				data, err := scheduler.Call("http://localhost/api/nrw/daily", []string{"month=2025-11", "device=40961"}, conn)

				if err != nil {
					fmt.Printf("error: %v", err)
				}

				formatted, _ := json.MarshalIndent(data, "", "  ")
				fmt.Printf("Data:\n%s\n", string(formatted))
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "no cached data found",
				})
			}
			fmt.Printf("error: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
		return c.Status(fiber.StatusOK).JSON(entry)
	})

}
