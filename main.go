package main

import (
	"context"
	"log"

	"github.com/Marc-Moonshot/temporal-guru/db"
	"github.com/gofiber/fiber/v3"
)

func main() {
	conn := db.Connect()
	defer conn.Close(context.Background())

	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	log.Fatal(app.Listen(":8000"))
}
