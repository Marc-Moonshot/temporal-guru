package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Marc-Moonshot/temporal-guru/types"
	"github.com/Marc-Moonshot/temporal-guru/utils"
	"github.com/jackc/pgx/v5"
)

// sets cache data in postgres
func Set(conn *pgx.Conn, url string, params []string, response types.Response) error {

	query := `INSERT INTO public."CacheEntry"
    (endpoint, query_params, query_hash, response, fetched_at, expires_at, status)
    VALUES($1, $2, $3, $4, $5, $6, $7)`

	fmt.Printf("-----\nquery: %s\nendpoint: %s\nparams: %s\nresponse: %v\n-----\n", query, url, params, response)

	now := time.Now().UTC()
	expires := now.Add(6 * time.Hour)

	entry := types.CacheEntry{
		Endpoint:   url,
		Fetched_at: now,
		Expires_at: expires,
		Status:     "valid",
		Query_hash: utils.HashParams(params),
	}

	entry.Query_params, _ = json.Marshal(params)
	entry.Response, _ = json.Marshal(response)

	responseStatus, err := conn.Exec(context.Background(), query, entry.Endpoint, entry.Query_params, entry.Query_hash, entry.Response, entry.Fetched_at,
		entry.Expires_at, entry.Status)

	fmt.Printf("Response Status: %v", responseStatus)
	if err != nil {
		return fmt.Errorf("[CACHE] error: %w", err)
	}
	return nil
}
