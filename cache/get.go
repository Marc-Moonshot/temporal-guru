package cache

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// gets cache data from postres
func Get(conn *pgx.Conn, endpoint string) (*CacheEntry, error) {

	fmt.Printf("GET: %s\n", endpoint)
	query :=`SELECT * from "CacheEntry" WHERE endpoint = $1`

	var response CacheEntry

	err := conn.QueryRow(context.Background(), query, endpoint).Scan(
		&response.Endpoint,
		&response.Query_params,
		&response.Response,
		&response.Fetched_at,
		&response.Expires_at,
		&response.Status,
	)

if err != nil {
		// Check if no rows were found
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		// For all other errors, wrap them
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return &response, nil

}
