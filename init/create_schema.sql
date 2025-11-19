CREATE SCHEMA IF NOT EXISTS cache_data;

CREATE TYPE cache_data.status AS ENUM ('valid', 'stale', 'error', 'pending');

CREATE TABLE cache_data.cacheEntry (
    id BIGSERIAL  PRIMARY KEY,  
    endpoint      TEXT NOT NULL,
    query_params  JSONB,
    query_hash    TEXT PRIMARY KEY,
    response      JSONB,
    fetched_at    TIMESTAMPTZ DEFAULT NOW(),
    expires_at    TIMESTAMPTZ,
    status        cache_data.status
);
