package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/Marc-Moonshot/temporal-guru/types"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UpdateOne(id string, field string, value any, pool *pgxpool.Pool) (pgconn.CommandTag, error) {
	// allowedFields := map[string]bool{
	// 	"fetched_at": true,
	// 	"status":     true,
	// }
	//
	// if !allowedFields[field] {
	// 	return pgconn.CommandTag{}, fmt.Errorf("[CACHE] invalid field name: %s", field)
	// }

	query := fmt.Sprintf(`UPDATE "CacheEntry" SET %s=$1 WHERE id=$2`, field)

	fmt.Printf("-----\n[CACHE]\nquery: %s\n%s: %s\nID: %s\n-----\n", query, field, value, id)

	tag, err := pool.Exec(context.Background(), query, value, id)

	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("[CACHE] error: %w", err)
	}
	fmt.Printf("[CACHE] tag: %v\n", tag)

	return tag, nil
}

func UpdateResponse(id string, response types.Response, pool *pgxpool.Pool) (pgconn.CommandTag, error) {
	query := `UPDATE "CacheEntry" SET response=$1, expires_at=$2 WHERE id=$3`

	expires_at := time.Now().Add(time.Hour * 6)
	fmt.Printf("-----\n[CACHE]\nquery: %s\nresponse: %v\nexpires_at: %s\nID: %s\n-----\n", query, response, expires_at, id)

	tag, err := pool.Exec(context.Background(), query, response, expires_at, id)

	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("[CACHE] error: %w", err)
	}
	fmt.Printf("[CACHE] tag: %v\n", tag)

	return tag, nil
}
