package scheduler

import (
	"fmt"

	"github.com/Marc-Moonshot/temporal-guru/cache"
	"github.com/jackc/pgx/v5/pgxpool"
)

// opens a go routine that calls the python API and inserts or updates a response
// TODO: wrap all this into a single transaction.
func HandleSchedulerCall(url string, params []string, pool *pgxpool.Pool, exists bool, id *string) {
	fmt.Println("Invoking Scheduler.")

	if exists && id != nil {
		if _, updateStatusStaleErr := cache.UpdateOne(*id, "status", "stale", pool); updateStatusStaleErr != nil {
			fmt.Printf("[SCHEDULER] update status to 'stale' failed: %v\n", updateStatusStaleErr)
		}
	}

	go func() {
		// TODO: TEST: set document to pending here
		if exists && id != nil {
			if _, updateStatusPendingErr := cache.UpdateOne(*id, "status", "pending", pool); updateStatusPendingErr != nil {
				fmt.Printf("[SCHEDULER] update status to 'stale' failed: %v\n", updateStatusPendingErr)
			}
		}
		data, err := Call(url, params)
		if err != nil {
			fmt.Printf("[SCHEDULER] scheduler call failed: %v\n", err)
			return
		}

		fmt.Printf("[SCHEDULER] async call succeeded, caching result.\n")

		if !exists {
			if _, err := cache.Set(pool, url, params, data); err != nil {
				fmt.Printf("[SCHEDULER] cache.Set failed: %v\n", err)
				return
			}
		} else if id != nil {

			if _, err := cache.UpdateResponse(*id, data, pool); err != nil {
				fmt.Printf("[SCHEDULER] Update response failed: %v\n", err)
				return
			}

			if _, updateStatusValidErr := cache.UpdateOne(*id, "status", "valid", pool); updateStatusValidErr != nil {
				fmt.Printf("[SCHEDULER] update status to 'valid' failed: %v\n", updateStatusValidErr)
			}
		}

		fmt.Printf("[SCHEDULER] result cached successfully.\n")
	}()
}
