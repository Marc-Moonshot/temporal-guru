package cache

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func UpdateOne(id string, field string, value any, conn *pgx.Conn) (pgx.Row, error) {
	query := `UPDATE public."CacheEntry" SET fetched_at=$1 WHERE id=$2`

	fmt.Printf("-----\nquery: %s\n%s:%s\nID: %s\n-----\n", query, field, value, id)

	var row pgx.Row
	tag, err := conn.Exec(context.Background(), query, value, id)

	fmt.Printf("[CACHE] response: %v", tag)
	if err != nil {
		return nil, fmt.Errorf("[CACHE] error: %w", err)
	}
	return row, nil
}
