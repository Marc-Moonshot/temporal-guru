package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Marc-Moonshot/temporal-guru/types"
	"github.com/jackc/pgx/v5"
)

var ErrResponseExpired = errors.New("response expired")

// gets cache data from postres
func Get(conn *pgx.Conn, endpoint string, queryHash string) (*types.CacheEntry, error) {

	query := `
    SELECT id, response, expires_at, status
    FROM "CacheEntry"
    WHERE endpoint = $1 AND query_hash = $2
    `
	fmt.Printf("-----\nquery: %s\nendpoint: %s\nquery hash: %s\n-----\n", query, endpoint, queryHash)

	var response types.CacheEntry

	err := conn.QueryRow(context.Background(), query, endpoint, queryHash).Scan(
		&response.ID,
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

	if response.Expires_at.Before(time.Now()) {
		return &response, ErrResponseExpired
	}
	return &response, nil

}
