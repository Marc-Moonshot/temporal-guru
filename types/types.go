package types

import (
	"encoding/json"
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

type StationReading struct {
	BilledCompleted float64 `json:"billed_completed"`
	BilledEst       float64 `json:"billed_est"`
	DailyFlow       float64 `json:"daily_flow"`
	Date            string  `json:"date"`
	DeviceCode      string  `json:"device_code"`
	NrwM3           float64 `json:"nrw_m3"`
	NrwPercent      float64 `json:"nrw_percent"`
}

type Response []map[string]StationReading


