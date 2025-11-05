package db

import (
	"encoding/json"

	"github.com/golang/protobuf/ptypes/timestamp"
)

type CacheStatus string

const (
	StatusValid   CacheStatus = "valid"
	StatusStale   CacheStatus = "stale"
	StatusError   CacheStatus = "error"
	StatusPending CacheStatus = "pending"
)

type CacheEntry struct {
	Endpoint     string              `json:"endpoint"`
	Query_params json.RawMessage     `json:"query_params"`
	Params_hash  string              `json:"params_hash"`
	Response     json.RawMessage     `json:"response"`
	Fetched_at   timestamp.Timestamp `json:"fetched_at"`
	Expires_at   timestamp.Timestamp `json:"expires_at"`
	Status       CacheStatus         `json:"cache_status"`
}
