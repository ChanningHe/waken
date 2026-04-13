package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port          string
	AuthToken     string
	BroadcastAddr string
	WOLPort       int
	DBPath        string
}

func Load() Config {
	return Config{
		Port:          envOrDefault("WOL_PORT", "19527"),
		AuthToken:     os.Getenv("WOL_AUTH_TOKEN"),
		BroadcastAddr: envOrDefault("WOL_BROADCAST_ADDR", "255.255.255.255"),
		WOLPort:       envOrDefaultInt("WOL_WOL_PORT", 9),
		DBPath:        envOrDefault("WOL_DB_PATH", "/app/waken/config/wol.db"),
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envOrDefaultInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}
