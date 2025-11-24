package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Marc-Moonshot/temporal-guru/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrResponseExpired = errors.New("response expired")

// gets cache data from postres
func Get(pool *pgxpool.Pool, endpoint string, queryHash string) (*types.CacheEntry, error) {

	query := `
    SELECT id, response, expires_at, status
    FROM "cacheentry"
    WHERE endpoint = $1 AND query_hash = $2
    `
	fmt.Printf("-----\n[CACHE]\nquery: %s\nendpoint: %s\nquery hash: %s\n", query, endpoint, queryHash)

	var response types.CacheEntry

	err := pool.QueryRow(context.Background(), query, endpoint, queryHash).Scan(
		&response.ID,
		&response.Response,
		&response.Expires_at,
		&response.Status,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			fmt.Println("[CACHE]No entry found.\n-----")
			return nil, pgx.ErrNoRows
		}
		return nil, fmt.Errorf("query failed: %w", err)
	}

	if response.Expires_at.Before(time.Now()) {
		fmt.Println("[CACHE] expires at: ", response.Expires_at)
		fmt.Println("[CACHE] now:", time.Now())
		fmt.Println("[CACHE] Entry expired.")
		return &response, ErrResponseExpired
	}
	return &response, nil

}
