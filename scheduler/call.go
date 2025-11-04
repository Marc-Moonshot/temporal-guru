package scheduler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type Result struct {
	Data any
	Err error
}

func Call(BaseUrl string, params []string, conn *pgx.Conn) (any, error) {
	const retries = 3
	const timeout = 120 * time.Second
	const backoff = 2 * time.Second

	client := http.Client{Timeout: timeout}

	var lastErr error

	if BaseUrl == "" {
		return nil, fmt.Errorf("no provided URL.")
	}

	u, err := url.Parse(BaseUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	q := u.Query()
	for _, p := range params {
		kv := strings.SplitN(p, "=", 2)
		if len(kv) == 2 {
			q.Add(kv[0], kv[1])
		}
	}
	u.RawQuery = q.Encode()

	fullUrl := u.String()

	for i := range retries {
		req, err := http.NewRequest(http.MethodGet, fullUrl, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		fmt.Printf("[SCHEDULER] calling: %s\n", fullUrl)
		res, err := client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("attempt %d failed: %w", i+1, err)
			time.Sleep(backoff * time.Duration(i+1))
			continue
		}

		defer res.Body.Close()

		if res.StatusCode >= 500 {
			lastErr = fmt.Errorf("server error (status %d) on attempt %d", res.StatusCode, i+1)
			time.Sleep(backoff * time.Duration(i+1))
			continue
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}

		var data any
		if err := json.Unmarshal(body, &data); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %w", err)
		}

		return data, nil
	}

	return nil, fmt.Errorf("API call failed after %d retries: %w", retries, lastErr)
}

// Calls an API asynchronously
func CallAsync(BaseUrl string, params []string, conn *pgx.Conn) <-chan Result {
	resultChan := make(chan Result, 1)
	
	go func(){
		data, err := Call(BaseUrl, params, conn)
		resultChan <- Result{Data: data, Err: err}
		close(resultChan)
	}()

	return resultChan
}
