package api

import (
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
		fmt.Println("[API] GET /nrw/yearly")

		year := c.Queries()["year"]
		device := c.Queries()["device"]

		if year == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Missing 'year' parameter",
			})
		}

		if device == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Missing 'device' parameter",
			})
		}

		entry, err := cache.Get(conn, "nrw/yearly")

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return handleSchedulerCall(c, "http://localhost/api/nrw/yearly", []string{"year=" + year, "device=" + device}, conn)
			}
			fmt.Printf("error: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
		return c.Status(fiber.StatusOK).JSON(entry)
	})

	app.Get("/nrw/monthly", func(c fiber.Ctx) error {
		fmt.Println("[API] GET /nrw/monthly")

		month := c.Queries()["month"]
		device := c.Queries()["device"]

		if month == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Missing 'month' parameter",
			})
		}

		if device == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Missing 'device' parameter",
			})
		}

		entry, err := cache.Get(conn, "nrw/monthly")

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return handleSchedulerCall(c, "http://localhost/api/nrw/monthly", []string{"month=" + month, "device=" + device}, conn)
			}
			fmt.Printf("error: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
		return c.Status(fiber.StatusOK).JSON(entry)
	})

	app.Get("/nrw/daily", func(c fiber.Ctx) error {

		fmt.Println("[API] GET /nrw/daily")

		month := c.Queries()["month"]
		device := c.Queries()["device"]

		if month == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Missing 'month' parameter",
			})
		}

		if device == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Missing 'device' parameter",
			})
		}

		entry, err := cache.Get(conn, "nrw/daily")

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return handleSchedulerCall(c, "http://localhost/api/nrw/daily", []string{"month=" + month, "device=" + device}, conn)
			}
			fmt.Printf("error: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
		return c.Status(fiber.StatusOK).JSON(entry)
	})

}

func handleSchedulerCall(c fiber.Ctx, url string, params []string, conn *pgx.Conn) error {
	fmt.Println("Invoking Scheduler.")
	data, apiErr := scheduler.Call(url, params, conn)

	if apiErr != nil {
		fmt.Printf("error: %v", apiErr)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}

	if data == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "no cached data found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"payload": data,
	})
}
