package api

import (
	"errors"
	"fmt"

	"github.com/Marc-Moonshot/temporal-guru/cache"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
)

// handles API calls routed from nginx
// calls cache/get
func RegisterRoutes(app *fiber.App, conn *pgx.Conn) {

	app.Get("/nrw/yearly", func(c fiber.Ctx) error {
		entry, err := cache.Get(conn, "nrw/yearly")

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
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
		entry, err := cache.Get(conn, "nrw/monthly")

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
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
		entry, err := cache.Get(conn, "nrw/daily")

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
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
