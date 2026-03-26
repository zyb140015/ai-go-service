package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultHTTPAddr          = ":8080"
	defaultReadTimeout       = 5 * time.Second
	defaultReadHeaderTimeout = 2 * time.Second
	defaultWriteTimeout      = 10 * time.Second
	defaultIdleTimeout       = 30 * time.Second
	defaultShutdownTimeout   = 10 * time.Second
	defaultLogLevel          = "INFO"
	defaultAuthTokenTTL      = 24 * time.Hour
	defaultGoAdminTimeout    = 10 * time.Second
	defaultAppEnv            = "development"
)

// Config stores application configuration loaded from the environment.
type Config struct {
	HTTPAddr          string
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
	LogLevel          string
	DatabaseURL       string
	AuthTokenSecret   string
	AuthTokenTTL      time.Duration
	GoAdminBaseURL    string
	GoAdminTimeout    time.Duration
	AppEnv            string
	EnableDesktopSeed bool
}

// Load returns the application configuration using environment overrides when present.
func Load() (Config, error) {
	config := Config{
		HTTPAddr:          getString("HTTP_ADDR", defaultHTTPAddr),
		ReadTimeout:       getDuration("HTTP_READ_TIMEOUT", defaultReadTimeout),
		ReadHeaderTimeout: getDuration("HTTP_READ_HEADER_TIMEOUT", defaultReadHeaderTimeout),
		WriteTimeout:      getDuration("HTTP_WRITE_TIMEOUT", defaultWriteTimeout),
		IdleTimeout:       getDuration("HTTP_IDLE_TIMEOUT", defaultIdleTimeout),
		ShutdownTimeout:   getDuration("HTTP_SHUTDOWN_TIMEOUT", defaultShutdownTimeout),
		LogLevel:          strings.ToUpper(getString("LOG_LEVEL", defaultLogLevel)),
		DatabaseURL:       getString("DATABASE_URL", ""),
		AuthTokenSecret:   getString("AUTH_TOKEN_SECRET", ""),
		AuthTokenTTL:      getDuration("AUTH_TOKEN_TTL", defaultAuthTokenTTL),
		GoAdminBaseURL:    getString("GOADMIN_BASE_URL", ""),
		GoAdminTimeout:    getDuration("GOADMIN_TIMEOUT", defaultGoAdminTimeout),
		AppEnv:            strings.ToLower(getString("APP_ENV", defaultAppEnv)),
	}
	config.EnableDesktopSeed = getBool("ENABLE_DESKTOP_SEED", config.AppEnv != "production")

	if config.HTTPAddr == "" {
		return Config{}, fmt.Errorf("HTTP_ADDR must not be empty")
	}

	if config.ReadTimeout <= 0 {
		return Config{}, fmt.Errorf("HTTP_READ_TIMEOUT must be positive")
	}

	if config.ReadHeaderTimeout <= 0 {
		return Config{}, fmt.Errorf("HTTP_READ_HEADER_TIMEOUT must be positive")
	}

	if config.WriteTimeout <= 0 {
		return Config{}, fmt.Errorf("HTTP_WRITE_TIMEOUT must be positive")
	}

	if config.IdleTimeout <= 0 {
		return Config{}, fmt.Errorf("HTTP_IDLE_TIMEOUT must be positive")
	}

	if config.ShutdownTimeout <= 0 {
		return Config{}, fmt.Errorf("HTTP_SHUTDOWN_TIMEOUT must be positive")
	}

	if config.AuthTokenTTL <= 0 {
		return Config{}, fmt.Errorf("AUTH_TOKEN_TTL must be positive")
	}

	if config.GoAdminBaseURL == "" {
		return Config{}, fmt.Errorf("GOADMIN_BASE_URL must not be empty")
	}

	if config.GoAdminTimeout <= 0 {
		return Config{}, fmt.Errorf("GOADMIN_TIMEOUT must be positive")
	}

	return config, nil
}

func getString(key string, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return strings.TrimSpace(value)
}

func getDuration(key string, fallback time.Duration) time.Duration {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	duration, err := time.ParseDuration(strings.TrimSpace(value))
	if err == nil {
		return duration
	}

	seconds, err := strconv.Atoi(strings.TrimSpace(value))
	if err == nil {
		return time.Duration(seconds) * time.Second
	}

	return fallback
}

func getBool(key string, fallback bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}
