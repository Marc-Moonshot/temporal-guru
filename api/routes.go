package api

import "github.com/gofiber/fiber/v3"

// handles API calls routed from nginx
// calls cache/get
func RegisterRoutes(app *fiber.App) {

	app.Get("/nrw/yearly", func(c fiber.Ctx) error {
		return c.SendString("yearly")
	})

	app.Get("/nrw/monthly", func(c fiber.Ctx) error {
		return c.SendString("monthly")
	})

	app.Get("/nrw/daily", func(c fiber.Ctx) error {
		return c.SendString("daily")
	})

}
