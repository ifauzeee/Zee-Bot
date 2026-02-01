-- +goose Up
CREATE TABLE IF NOT EXISTS kv_store (
    key TEXT PRIMARY KEY,
    value TEXT
);

-- +goose Down
DROP TABLE IF EXISTS kv_store;
