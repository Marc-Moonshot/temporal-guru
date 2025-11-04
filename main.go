package main

import (
	"context"
	"log"

	"github.com/Marc-Moonshot/temporal-guru/api"
	"github.com/Marc-Moonshot/temporal-guru/db"
	"github.com/gofiber/fiber/v3"
)

func main() {
	conn := db.Connect()
	defer conn.Close(context.Background())

	app := fiber.New()

	api.RegisterRoutes(app, conn)

	log.Fatal(app.Listen(":8000"))

}
