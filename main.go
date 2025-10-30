package main

import (
	"context"

	"github.com/Marc-Moonshot/temporal-guru/db"
)

func main() {
	conn := db.Connect()
	defer conn.Close(context.Background())
}
