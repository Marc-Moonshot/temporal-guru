package types

import (
	"encoding/json"
	"fmt"
	"time"
)

type CacheStatus string

const (
	StatusValid      CacheStatus = "valid"
	StatusStale      CacheStatus = "stale"
	StatusError      CacheStatus = "error"
	StatusProcessing CacheStatus = "processing"
)

type CacheEntry struct {
	ID           string          `json:"id" db:"id"`
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

type YearlyReading struct {
	BilledCompleted string  `json:"billed_completed"`
	BilledQty       float64 `json:"billed_qty"`
	DeviceCode      string  `json:"device_code"`
	Month           string  `json:"month"`
	NrwM3           float64 `json:"nrw_m3"`
	NrwPercent      float64 `json:"nrw_percent"`
	TotalFlow       float64 `json:"total_flow"`
}

type Response struct {
	DailyData   []map[string]DailyReading           `json:"-"`
	MonthlyData map[string]MonthlyReading           `json:"-"`
	YearlyData  map[string]map[string]YearlyReading `json:"-"`
}

func (r *Response) UnmarshalJSON(data []byte) error {
	// Try to parse as array first (DailyData)
	var arr []map[string]DailyReading
	if err := json.Unmarshal(data, &arr); err == nil {
		r.DailyData = arr
		return nil
	}

	// Try to parse as nested map (YearlyData)
	var yearlyObj map[string]map[string]YearlyReading
	if err := json.Unmarshal(data, &yearlyObj); err == nil {
		// Verify it's actually yearly by checking if values are maps
		for _, v := range yearlyObj {
			if len(v) > 0 {
				// If we have nested maps, it's yearly data
				r.YearlyData = yearlyObj
				return nil
			}
		}
	}

	// Finally try MonthlyData (single level map)
	var monthlyObj map[string]MonthlyReading
	if err := json.Unmarshal(data, &monthlyObj); err == nil {
		r.MonthlyData = monthlyObj
		return nil
	}

	return fmt.Errorf("Response: unsupported JSON structure")
}

func (r Response) MarshalJSON() ([]byte, error) {
	if r.DailyData != nil {
		return json.Marshal(r.DailyData)
	}
	if r.MonthlyData != nil {
		return json.Marshal(r.MonthlyData)
	}
	if r.YearlyData != nil {
		return json.Marshal(r.YearlyData)
	}
	return nil, fmt.Errorf("Response: no data to marshal")
}
