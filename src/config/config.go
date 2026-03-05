package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const defaultJWTSecret = "change-me-in-production"

type Config struct {
	Port             string
	Env              string
	LogLevel         string
	RequestTimeout   time.Duration
	JWTSecret        string
	CORSAllowOrigins string // comma-separated; empty means allow all origins
}

func Load() *Config {
	// .env is optional; ignore error when file is absent.
	_ = godotenv.Load()

	cfg := &Config{
		Port:             getEnv("PORT", "3000"),
		Env:              getEnv("ENV", "development"),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
		RequestTimeout:   parseDuration("REQUEST_TIMEOUT", 30*time.Second),
		JWTSecret:        getEnv("JWT_SECRET", defaultJWTSecret),
		CORSAllowOrigins: getEnv("CORS_ALLOW_ORIGINS", ""),
	}

	if err := cfg.validate(); err != nil {
		panic(fmt.Sprintf("invalid configuration: %v", err))
	}

	return cfg
}

func (c *Config) validate() error {
	if c.Env == "production" && c.JWTSecret == defaultJWTSecret {
		return fmt.Errorf("JWT_SECRET must be changed from the default value in production")
	}
	return nil
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func parseDuration(key string, defaultVal time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			log.Printf("WARN: invalid value for %s (%q): %v; using default %s", key, v, err, defaultVal)
			return defaultVal
		}
		return d
	}
	return defaultVal
}
