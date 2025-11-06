package scheduler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Marc-Moonshot/temporal-guru/types"
)

type Result struct {
	Data types.Response
	Err  error
}

func Call(BaseUrl string, params []string) (types.Response, error) {
	const retries = 3
	const timeout = 300 * time.Second
	const backoff = 2 * time.Second

	client := http.Client{Timeout: timeout}

	var lastErr error

	if BaseUrl == "" {
		return types.Response{}, fmt.Errorf("no provided URL.")
	}

	u, err := url.Parse(os.Getenv("PY_API_URL") + BaseUrl)
	if err != nil {
		return types.Response{}, fmt.Errorf("invalid URL: %w", err)
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
			return types.Response{}, fmt.Errorf("failed to create request: %w", err)
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
			return types.Response{}, fmt.Errorf("failed to read response: %w", err)
		}

		var data types.Response
		if err := json.Unmarshal(body, &data); err != nil {
			return types.Response{}, fmt.Errorf("failed to parse JSON: %w", err)

		}

		return data, nil
	}

	return types.Response{}, fmt.Errorf("API call failed after %d retries: %w", retries, lastErr)
}

