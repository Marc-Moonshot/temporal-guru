package scheduler

import (
	"fmt"

	"github.com/Marc-Moonshot/temporal-guru/cache"
	"github.com/jackc/pgx/v5/pgxpool"
)

// opens a go routine that calls the python API and inserts or updates a response
func HandleSchedulerCall(url string, params []string, pool *pgxpool.Pool, id *string) {
	fmt.Println("Invoking Scheduler.")

	go func() {
		// TEST: set row to pending here
		if id != nil {
			if _, updateStatusPendingErr := cache.UpdateOne(*id, "status", "pending", pool); updateStatusPendingErr != nil {
				fmt.Printf("[SCHEDULER] update status to 'pending' failed: %v\n", updateStatusPendingErr)
			}
		}

		data, err := Call(url, params)
		if err != nil {
			if id != nil {
				// TEST: set row to error
				if _, updateStatusErrorErr := cache.UpdateOne(*id, "status", "error", pool); updateStatusErrorErr != nil {
					fmt.Printf("[SCHEDULER] update status to 'error' failed: %v\n", updateStatusErrorErr)
				}
			}

			fmt.Printf("[SCHEDULER] scheduler call failed: %v\n", err)
			return
		}

		fmt.Printf("[SCHEDULER] async call succeeded, caching result.\n")

		if id != nil {
			if _, err := cache.UpdateResponse(*id, data, pool); err != nil {
				// TEST: set row to error
				if _, updateStatusErrorErr := cache.UpdateOne(*id, "status", "error", pool); updateStatusErrorErr != nil {
					fmt.Printf("[SCHEDULER] update status to 'error' failed: %v\n", updateStatusErrorErr)
				}

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
