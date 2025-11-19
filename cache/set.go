package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Marc-Moonshot/temporal-guru/types"
	"github.com/Marc-Moonshot/temporal-guru/utils"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// sets cache data in postgres
func Set(pool *pgxpool.Pool, url string, params []string, response *types.Response, status types.CacheStatus) (pgconn.CommandTag, error) {

	query := `INSERT INTO "CacheEntry"
    (endpoint, query_params, query_hash, response, fetched_at, expires_at, status)
    VALUES($1, $2, $3, $4, $5, $6, $7)`

	fmt.Printf("-----\n[CACHE]\nquery: %s\nendpoint: %s\nparams: %s\nresponse: %v\n-----\n", query, url, params, response)

	now := time.Now().UTC()
	expires := now.Add(6 * time.Hour)

	entry := types.CacheEntry{
		Endpoint:   url,
		Fetched_at: now,
		Expires_at: expires,
		Status:     status,
		Query_hash: utils.HashParams(params),
	}

	entry.Query_params, _ = json.Marshal(params)
	entry.Response, _ = json.Marshal(response)

	commandTag, err := pool.Exec(context.Background(), query, entry.Endpoint, entry.Query_params, entry.Query_hash, entry.Response, entry.Fetched_at,
		entry.Expires_at, entry.Status)

	fmt.Printf("[CACHE] tag: %v\n", commandTag)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("[CACHE] error: %w", err)
	}
	return commandTag, nil
}
