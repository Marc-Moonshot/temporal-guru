package api

import (
	// "encoding/json"
	"errors"
	"fmt"
	"time"
	"github.com/Marc-Moonshot/temporal-guru/cache"
	"github.com/Marc-Moonshot/temporal-guru/scheduler"
	"github.com/Marc-Moonshot/temporal-guru/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Handles API calls routed from nginx
func RegisterRoutes(app *fiber.App, pool *pgxpool.Pool) {

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

		fullParams := []string{`year=` + year, `device=` + device}

		paramsHash := utils.HashParams(fullParams)
		entry, err := cache.Get(pool, "/nrw/yearly", paramsHash)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				// Entry not found: create a new one and trigger fetch
				if _, err := cache.Set(pool, "/nrw/yearly", []string{"year=" + year, "device=" + device}, nil, "pending"); err != nil {
					fmt.Printf("[API] cache.Set failed: %v\n", err)
				}

				updatedEntry, err := cache.Get(pool, "/nrw/yearly", paramsHash)

				if err != nil {
					fmt.Printf("[API] cache.Get failed: %v\n", err)
				}

				scheduler.HandleSchedulerCall("/nrw/yearly", []string{"year=" + year, "device=" + device}, pool, &updatedEntry.ID)
				return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
					"message": "Request accepted. Data is being fetched and cached.",
				})
			}

			if errors.Is(err, cache.ErrResponseExpired) {
				// Expired entry exists
				if entry.Status != "pending" {
					idPtr := &entry.ID
					scheduler.HandleSchedulerCall("/nrw/yearly", []string{"year=" + year, "device=" + device}, pool, idPtr)
				}
				return c.Status(fiber.StatusOK).JSON(entry.Response)
			}

			fmt.Printf("error: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}

		fetchedAt := time.Now()
		if _, updateErr := cache.UpdateOne(entry.ID, "fetched_at", fetchedAt, pool); updateErr != nil {
			fmt.Printf("[API] update error: %v", updateErr)
		}

		if entry.Status == "pending" {
			return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
				"message": "Request accepted. Data is currently being fetched.",
			})
		}

		if entry.Status == "error" {
			updatedEntry, err := cache.Get(pool, "/nrw/yearly", paramsHash)

			if err != nil {
				fmt.Printf("[API] cache.Get failed: %v\n", err)
			}

			scheduler.HandleSchedulerCall("/nrw/yearly", []string{"year=" + year, "device=" + device}, pool, &updatedEntry.ID)
			return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
				"message": "Request accepted. Data is being revalidated.",
			})
		}

		return c.Status(fiber.StatusOK).JSON(entry.Response)
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
		fullParams := []string{`month=` + month, `device=` + device}

		paramsHash := utils.HashParams(fullParams)
		entry, err := cache.Get(pool, "/nrw/monthly", paramsHash)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				// Entry not found: create a new one and trigger fetch
				if _, err := cache.Set(pool, "/nrw/monthly", []string{"month=" + month, "device=" + device}, nil, "pending"); err != nil {
					fmt.Printf("[API] cache.Set failed: %v\n", err)
				}

				updatedEntry, err := cache.Get(pool, "/nrw/monthly", paramsHash)

				if err != nil {
					fmt.Printf("[API] cache.Get failed: %v\n", err)
				}

				scheduler.HandleSchedulerCall("/nrw/monthly", []string{"month=" + month, "device=" + device}, pool, &updatedEntry.ID)
				return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
					"message": "Request accepted. Data is being fetched and cached.",
				})
			}

			if errors.Is(err, cache.ErrResponseExpired) {
				// Expired entry exists
				if entry.Status != "pending" {
					idPtr := &entry.ID
					scheduler.HandleSchedulerCall("/nrw/monthly", []string{"month=" + month, "device=" + device}, pool, idPtr)
				}
				return c.Status(fiber.StatusOK).JSON(entry.Response)
			}

			fmt.Printf("error: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}

		fetchedAt := time.Now()
		if _, updateErr := cache.UpdateOne(entry.ID, "fetched_at", fetchedAt, pool); updateErr != nil {
			fmt.Printf("[API] update error: %v", updateErr)
		}

		if entry.Status == "pending" {
			return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
				"message": "Request accepted. Data is currently being fetched.",
			})
		}

		if entry.Status == "error" {
			updatedEntry, err := cache.Get(pool, "/nrw/monthly", paramsHash)

			if err != nil {
				fmt.Printf("[API] cache.Get failed: %v\n", err)
			}

			scheduler.HandleSchedulerCall("/nrw/monthly", []string{"month=" + month, "device=" + device}, pool, &updatedEntry.ID)
			return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
				"message": "Request accepted. Data is being revalidated.",
			})
		}

		return c.Status(fiber.StatusOK).JSON(entry.Response)
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

		fullParams := []string{`month=` + month, `device=` + device}

		paramsHash := utils.HashParams(fullParams)
		entry, err := cache.Get(pool, "/nrw/daily", paramsHash)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				// Entry not found: create a new one and trigger fetch
				if _, err := cache.Set(pool, "/nrw/daily", []string{"month=" + month, "device=" + device}, nil, "pending"); err != nil {
					fmt.Printf("[API] cache.Set failed: %v\n", err)
				}

				updatedEntry, err := cache.Get(pool, "/nrw/daily", paramsHash)

				if err != nil {
					fmt.Printf("[API] cache.Get failed: %v\n", err)
				}

				scheduler.HandleSchedulerCall("/nrw/daily", []string{"month=" + month, "device=" + device}, pool, &updatedEntry.ID)
				return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
					"message": "Request accepted. Data is being fetched and cached.",
				})
			}

			if errors.Is(err, cache.ErrResponseExpired) {
				// Expired entry exists
				if entry.Status != "pending" {
					idPtr := &entry.ID
					scheduler.HandleSchedulerCall("/nrw/daily", []string{"month=" + month, "device=" + device}, pool, idPtr)
				}
				return c.Status(fiber.StatusOK).JSON(entry.Response)
			}

			fmt.Printf("error: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}

		fetchedAt := time.Now()
		if _, updateErr := cache.UpdateOne(entry.ID, "fetched_at", fetchedAt, pool); updateErr != nil {
			fmt.Printf("[API] update error: %v", updateErr)
		}

		if entry.Status == "pending" {
			return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
				"message": "Request accepted. Data is currently being fetched.",
			})
		}

		if entry.Status == "error" {
			updatedEntry, err := cache.Get(pool, "/nrw/daily", paramsHash)

			if err != nil {
				fmt.Printf("[API] cache.Get failed: %v\n", err)
			}

			scheduler.HandleSchedulerCall("/nrw/daily", []string{"month=" + month, "device=" + device}, pool, &updatedEntry.ID)
			return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
				"message": "Request accepted. Data is being revalidated.",
			})
		}

		return c.Status(fiber.StatusOK).JSON(entry.Response)
	})

	app.Get("/status/*", func(c fiber.Ctx) error {

		endpoint := fmt.Sprintf("/%s", c.Params("*"))
		queryParams := c.Queries()

		queryArray := []string{}

		for k, v := range queryParams {
			queryArray = append(queryArray, fmt.Sprintf("%s=%s", k, v))
		}

		queryHash := utils.HashParams(queryArray)

		fmt.Printf("[API] GET /status/%s\n", endpoint)
		fmt.Printf("[API] queryHash: %s\n", queryHash)

		data, err := cache.Get(pool, endpoint, queryHash)

		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "entry not found"})
		}

		if err != nil && !errors.Is(err, cache.ErrResponseExpired) {
			fmt.Printf("[API] error: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}

		// if errors.Is(err, cache.ErrResponseExpired) {
		// 	return c.Status(fiber.StatusOK).JSON("expired")
		// }

		return c.Status(fiber.StatusOK).JSON(data.Status)
	})
}
