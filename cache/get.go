package cache

import (
	"context"
	"errors"
	"fmt"

	"github.com/Marc-Moonshot/temporal-guru/types"
	"github.com/jackc/pgx/v5"
)

// gets cache data from postres
func Get(conn *pgx.Conn, endpoint string, queryHash string) (*types.CacheEntry, error) {

	query := `
    SELECT response, expires_at, status
    FROM "CacheEntry"
    WHERE endpoint = $1 AND query_hash = $2
    `
	fmt.Printf("-----\nquery: %s\nendpoint: %s\nquery hash: %s\n-----\n", query, endpoint, queryHash)

	var response types.CacheEntry

	err := conn.QueryRow(context.Background(), query, endpoint, queryHash).Scan(
		&response.Response,
		&response.Expires_at,
		&response.Status,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return &response, nil

}
