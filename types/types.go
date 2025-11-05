package types

import (
	"encoding/json"
	"fmt"
	"time"
)

type CacheStatus string

const (
	StatusValid   CacheStatus = "valid"
	StatusStale   CacheStatus = "stale"
	StatusError   CacheStatus = "error"
	StatusPending CacheStatus = "pending"
)

type CacheEntry struct {
	Endpoint     string          `json:"endpoint"`
	Query_params json.RawMessage `json:"query_params"`
	Query_hash   string          `json:"query_hash"`
	Response     json.RawMessage `json:"response"`
	Fetched_at   time.Time       `json:"fetched_at"`
	Expires_at   time.Time       `json:"expires_at"`
	Status       CacheStatus     `json:"cache_status"`
}

type DailyReading struct {
	BilledCompleted float64 `json:"billed_completed"`
	BilledEst       float64 `json:"billed_est"`
	DailyFlow       float64 `json:"daily_flow"`
	Date            string  `json:"date"`
	DeviceCode      string  `json:"device_code"`
	NrwM3           float64 `json:"nrw_m3"`
	NrwPercent      float64 `json:"nrw_percent"`
}

type MonthlyReading struct {
	BilledCompleted string  `json:"billed_completed"`
	BilledQty       float64 `json:"billed_qty"`
	DeviceCode      string  `json:"device_code"`
	NrwM3           float64 `json:"nrw_m3"`
	NrwPercent      float64 `json:"nrw_percent"`
	TotalFlow       float64 `json:"total_flow"`
}

type Response struct {
	DailyData   []map[string]DailyReading
	MonthlyData map[string]MonthlyReading
}

func (r *Response) UnmarshalJSON(data []byte) error {
	var arr []map[string]DailyReading
	if err := json.Unmarshal(data, &arr); err == nil {
		r.DailyData = arr
		return nil
	}

	// Fallback to map format
	var obj map[string]MonthlyReading
	if err := json.Unmarshal(data, &obj); err == nil {
		r.MonthlyData = obj
		return nil
	}

	// Neither matched
	return fmt.Errorf("Response: unsupported JSON structure")
}
