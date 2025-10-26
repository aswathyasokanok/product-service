package config

import (
	"os"
	"strconv"
	"time"
)

// config holds application configuration
type Config struct {
	Workers   int
	QueueSize int
	Port      string

	// High throughput configuration
	BatchSize          int
	BatchFlushInterval time.Duration

	// Error handling configuration
	MaxRetryAttempts        int
	InitialRetryDelay       time.Duration
	MaxRetryDelay           time.Duration
	CircuitBreakerThreshold int
	CircuitBreakerTimeout   time.Duration

	// Memory management
	MaxMemoryUsage   int64
	CleanupThreshold float64
	GCInterval       time.Duration
}

// load the config from the environment variables
func LoadConfig() *Config {
	return &Config{
		Workers:   getEnvInt("WORKERS", 3),
		QueueSize: getEnvInt("QUEUE_SIZE", 1000),
		Port:      getEnv("PORT", "8080"),

		// High throughput configuration
		BatchSize:          getEnvInt("BATCH_SIZE", 100),
		BatchFlushInterval: getEnvDuration("BATCH_FLUSH_INTERVAL", 1*time.Second),

		// Error handling configuration
		MaxRetryAttempts:        getEnvInt("MAX_RETRY_ATTEMPTS", 3),
		InitialRetryDelay:       getEnvDuration("INITIAL_RETRY_DELAY", 100*time.Millisecond),
		MaxRetryDelay:           getEnvDuration("MAX_RETRY_DELAY", 30*time.Second),
		CircuitBreakerThreshold: getEnvInt("CIRCUIT_BREAKER_THRESHOLD", 5),
		CircuitBreakerTimeout:   getEnvDuration("CIRCUIT_BREAKER_TIMEOUT", 60*time.Second),

		// Memory management
		MaxMemoryUsage:   getEnvInt64("MAX_MEMORY_USAGE", 1024*1024*1024), // 1GB
		CleanupThreshold: getEnvFloat64("CLEANUP_THRESHOLD", 0.8),
		GCInterval:       getEnvDuration("GC_INTERVAL", 30*time.Second),
	}
}

// helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvFloat64(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
