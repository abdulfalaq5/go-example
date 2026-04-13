package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	AppEnv   string
	AppPort  string
	LogLevel        string
	KeycloakJWKSURL string
	DB              DBConfig
	Storage         StorageConfig
}

// StorageConfig holds configuration for the file storage service.
type StorageConfig struct {
	Type           string
	LocalUploadDir string
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
	MinioUseSSL    bool
}

// DBConfig holds configuration for all database connection pools.
type DBConfig struct {
	MainDSN  string
	AuditDSN string

	// Pool settings (shared between Main and Audit)
	MaxConns            int32
	MinConns            int32
	MaxConnLifetime     time.Duration
	MaxConnIdleTime     time.Duration
	HealthCheckPeriod   time.Duration
}

// Load reads the .env file (if present) and populates a Config struct.
// Missing optional keys fall back to sensible defaults.
func Load() (*Config, error) {
	// Load .env – ignore error when the file simply doesn't exist (e.g. production).
	_ = godotenv.Load()

	cfg := &Config{
		AppEnv:          getEnv("APP_ENV", "production"),
		AppPort:         getEnv("APP_PORT", "8080"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		KeycloakJWKSURL: getEnv("KEYCLOAK_JWKS_URL", ""),
		DB: DBConfig{
			MainDSN:           getEnv("DB_MAIN_DSN", ""),
			AuditDSN:          getEnv("DB_AUDIT_DSN", ""),
			MaxConns:          int32(getEnvInt("DB_MAX_CONNS", 25)),
			MinConns:          int32(getEnvInt("DB_MIN_CONNS", 5)),
			MaxConnLifetime:   getEnvDuration("DB_MAX_CONN_LIFETIME", 30*time.Minute),
			MaxConnIdleTime:   getEnvDuration("DB_MAX_CONN_IDLE_TIME", 15*time.Minute),
			HealthCheckPeriod: getEnvDuration("DB_HEALTH_CHECK_PERIOD", 1*time.Minute),
		},
		Storage: StorageConfig{
			Type:           getEnv("STORAGE_TYPE", "local"),
			LocalUploadDir: getEnv("LOCAL_UPLOAD_DIR", "./uploads"),
			MinioEndpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
			MinioAccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			MinioSecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
			MinioBucket:    getEnv("MINIO_BUCKET", "uploads"),
			MinioUseSSL:    getEnvBool("MINIO_USE_SSL", false),
		},
	}

	if cfg.AppPort == "" {
		return nil, fmt.Errorf("APP_PORT must not be empty")
	}
	if cfg.DB.MainDSN == "" {
		return nil, fmt.Errorf("DB_MAIN_DSN must not be empty")
	}
	if cfg.DB.AuditDSN == "" {
		return nil, fmt.Errorf("DB_AUDIT_DSN must not be empty")
	}

	return cfg, nil
}

// getEnv returns the value for key, or fallback when the key is unset / empty.
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// getEnvInt parses an integer env var, returning fallback on parse error.
func getEnvInt(key string, fallback int) int {
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

// getEnvDuration parses a duration env var (e.g. "30m"), returning fallback
// on parse error.
func getEnvDuration(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}

// getEnvBool parses a boolean env var (e.g. "true"), returning fallback
// on parse error.
func getEnvBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}
