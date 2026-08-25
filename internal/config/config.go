package config

import (
	"os"
	"strconv"
)

// Config 服务配置。
type Config struct {
	Addr             string
	BucketSize       int64
	Workers          int
	RetentionSeconds int64
	PruneTTL         int64
}

// Load 从环境变量加载配置。
func Load() Config {
	return Config{
		Addr:       getenv("TELEMETRY_ADDR", ":8080"),
		BucketSize: getenvInt("TELEMETRY_BUCKET_SECONDS", 60),
		Workers:    int(getenvInt("TELEMETRY_WORKERS", 4)),
		RetentionSeconds: getenvInt("TELEMETRY_RETENTION_SECONDS", 86400),
		PruneTTL:         getenvInt("TELEMETRY_PRUNE_TTL", 120),
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getenvInt(key string, fallback int64) int64 {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
			return parsed
		}
	}
	return fallback
}
