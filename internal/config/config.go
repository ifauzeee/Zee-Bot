package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppID         int
	AppHash       string
	SessionString string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	appIDStr := os.Getenv("APP_ID")
	if appIDStr == "" {
		appIDStr = os.Getenv("TELEGRAM_API_ID")
	}
	if appIDStr == "" {
		return nil, fmt.Errorf("APP_ID or TELEGRAM_API_ID is not set")
	}

	appID, err := strconv.Atoi(appIDStr)
	if err != nil {
		return nil, fmt.Errorf("APP_ID must be an integer: %w", err)
	}

	appHash := os.Getenv("APP_HASH")
	if appHash == "" {
		appHash = os.Getenv("TELEGRAM_API_HASH")
	}
	if appHash == "" {
		return nil, fmt.Errorf("APP_HASH or TELEGRAM_API_HASH is not set")
	}

	sessionString := os.Getenv("SESSION_STRING")

	return &Config{
		AppID:         appID,
		AppHash:       appHash,
		SessionString: sessionString,
	}, nil
}
