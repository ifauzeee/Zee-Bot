package database

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

//go:embed migrations/*.sql
var embedMigrations embed.FS

func Init() error {
	dbPath := "zee.db"
	if dir := filepath.Dir(dbPath); dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.Up(DB, "migrations"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func SetKV(key, value string) error {
	_, err := DB.Exec("INSERT OR REPLACE INTO kv_store (key, value) VALUES (?, ?)", key, value)
	return err
}

func GetKV(key string) (string, error) {
	var value string
	err := DB.QueryRow("SELECT value FROM kv_store WHERE key = ?", key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}

func DeleteKV(key string) error {
	_, err := DB.Exec("DELETE FROM kv_store WHERE key = ?", key)
	return err
}
