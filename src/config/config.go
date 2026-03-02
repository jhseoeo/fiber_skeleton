package config

import (
	"os"
	"time"
)

type Config struct {
	Port           string
	Env            string
	LogLevel       string
	RequestTimeout time.Duration
	JWTSecret      string
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "3000"),
		Env:            getEnv("ENV", "development"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		RequestTimeout: parseDuration("REQUEST_TIMEOUT", 30*time.Second),
		JWTSecret:      getEnv("JWT_SECRET", "change-me-in-production"),
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func parseDuration(key string, defaultVal time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return defaultVal
}
