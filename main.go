package main

import (
	"log"

	"github.com/Marc-Moonshot/temporal-guru/api"
	"github.com/Marc-Moonshot/temporal-guru/db"
	"github.com/gofiber/fiber/v3"
)

func main() {
	pool := db.Connect()
	defer pool.Close()
	app := fiber.New()

	api.RegisterRoutes(app, pool)

	log.Fatal(app.Listen(":8000"))

}
