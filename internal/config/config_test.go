package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig_DefaultValues(t *testing.T) {
	// Clear environment variables
	os.Clearenv()

	config := LoadConfig()

	if config.Workers != 3 {
		t.Errorf("Expected Workers 3, got %d", config.Workers)
	}
	if config.QueueSize != 1000 {
		t.Errorf("Expected QueueSize 1000, got %d", config.QueueSize)
	}
	if config.Port != "8080" {
		t.Errorf("Expected Port '8080', got '%s'", config.Port)
	}
	if config.BatchSize != 100 {
		t.Errorf("Expected BatchSize 100, got %d", config.BatchSize)
	}
	if config.BatchFlushInterval != 1*time.Second {
		t.Errorf("Expected BatchFlushInterval 1s, got %v", config.BatchFlushInterval)
	}
	if config.MaxRetryAttempts != 3 {
		t.Errorf("Expected MaxRetryAttempts 3, got %d", config.MaxRetryAttempts)
	}
	if config.InitialRetryDelay != 100*time.Millisecond {
		t.Errorf("Expected InitialRetryDelay 100ms, got %v", config.InitialRetryDelay)
	}
	if config.MaxRetryDelay != 30*time.Second {
		t.Errorf("Expected MaxRetryDelay 30s, got %v", config.MaxRetryDelay)
	}
	if config.CircuitBreakerThreshold != 5 {
		t.Errorf("Expected CircuitBreakerThreshold 5, got %d", config.CircuitBreakerThreshold)
	}
	if config.CircuitBreakerTimeout != 60*time.Second {
		t.Errorf("Expected CircuitBreakerTimeout 60s, got %v", config.CircuitBreakerTimeout)
	}
	if config.MaxMemoryUsage != 1024*1024*1024 {
		t.Errorf("Expected MaxMemoryUsage 1GB, got %d", config.MaxMemoryUsage)
	}
	if config.CleanupThreshold != 0.8 {
		t.Errorf("Expected CleanupThreshold 0.8, got %f", config.CleanupThreshold)
	}
	if config.GCInterval != 30*time.Second {
		t.Errorf("Expected GCInterval 30s, got %v", config.GCInterval)
	}
}

func TestLoadConfig_EnvironmentVariables(t *testing.T) {
	// Set environment variables
	os.Setenv("WORKERS", "5")
	os.Setenv("QUEUE_SIZE", "2000")
	os.Setenv("PORT", "9090")
	os.Setenv("BATCH_SIZE", "200")
	os.Setenv("BATCH_FLUSH_INTERVAL", "2s")
	os.Setenv("MAX_RETRY_ATTEMPTS", "5")
	os.Setenv("INITIAL_RETRY_DELAY", "200ms")
	os.Setenv("MAX_RETRY_DELAY", "60s")
	os.Setenv("CIRCUIT_BREAKER_THRESHOLD", "10")
	os.Setenv("CIRCUIT_BREAKER_TIMEOUT", "120s")
	os.Setenv("MAX_MEMORY_USAGE", "2147483648") // 2GB
	os.Setenv("CLEANUP_THRESHOLD", "0.9")
	os.Setenv("GC_INTERVAL", "60s")

	config := LoadConfig()

	if config.Workers != 5 {
		t.Errorf("Expected Workers 5, got %d", config.Workers)
	}
	if config.QueueSize != 2000 {
		t.Errorf("Expected QueueSize 2000, got %d", config.QueueSize)
	}
	if config.Port != "9090" {
		t.Errorf("Expected Port '9090', got '%s'", config.Port)
	}
	if config.BatchSize != 200 {
		t.Errorf("Expected BatchSize 200, got %d", config.BatchSize)
	}
	if config.BatchFlushInterval != 2*time.Second {
		t.Errorf("Expected BatchFlushInterval 2s, got %v", config.BatchFlushInterval)
	}
	if config.MaxRetryAttempts != 5 {
		t.Errorf("Expected MaxRetryAttempts 5, got %d", config.MaxRetryAttempts)
	}
	if config.InitialRetryDelay != 200*time.Millisecond {
		t.Errorf("Expected InitialRetryDelay 200ms, got %v", config.InitialRetryDelay)
	}
	if config.MaxRetryDelay != 60*time.Second {
		t.Errorf("Expected MaxRetryDelay 60s, got %v", config.MaxRetryDelay)
	}
	if config.CircuitBreakerThreshold != 10 {
		t.Errorf("Expected CircuitBreakerThreshold 10, got %d", config.CircuitBreakerThreshold)
	}
	if config.CircuitBreakerTimeout != 120*time.Second {
		t.Errorf("Expected CircuitBreakerTimeout 120s, got %v", config.CircuitBreakerTimeout)
	}
	if config.MaxMemoryUsage != 2147483648 {
		t.Errorf("Expected MaxMemoryUsage 2GB, got %d", config.MaxMemoryUsage)
	}
	if config.CleanupThreshold != 0.9 {
		t.Errorf("Expected CleanupThreshold 0.9, got %f", config.CleanupThreshold)
	}
	if config.GCInterval != 60*time.Second {
		t.Errorf("Expected GCInterval 60s, got %v", config.GCInterval)
	}

	// Clean up
	os.Clearenv()
}

func TestLoadConfig_InvalidValues(t *testing.T) {
	// Set invalid environment variables
	os.Setenv("WORKERS", "invalid")
	os.Setenv("QUEUE_SIZE", "not-a-number")
	os.Setenv("PORT", "8080") // Valid
	os.Setenv("BATCH_SIZE", "abc")
	os.Setenv("BATCH_FLUSH_INTERVAL", "invalid-duration")
	os.Setenv("MAX_RETRY_ATTEMPTS", "xyz")
	os.Setenv("INITIAL_RETRY_DELAY", "invalid-duration")
	os.Setenv("MAX_RETRY_DELAY", "invalid-duration")
	os.Setenv("CIRCUIT_BREAKER_THRESHOLD", "invalid")
	os.Setenv("CIRCUIT_BREAKER_TIMEOUT", "invalid-duration")
	os.Setenv("MAX_MEMORY_USAGE", "invalid")
	os.Setenv("CLEANUP_THRESHOLD", "invalid")
	os.Setenv("GC_INTERVAL", "invalid-duration")

	config := LoadConfig()

	// Should fall back to default values
	if config.Workers != 3 {
		t.Errorf("Expected default Workers 3, got %d", config.Workers)
	}
	if config.QueueSize != 1000 {
		t.Errorf("Expected default QueueSize 1000, got %d", config.QueueSize)
	}
	if config.Port != "8080" {
		t.Errorf("Expected Port '8080', got '%s'", config.Port)
	}
	if config.BatchSize != 100 {
		t.Errorf("Expected default BatchSize 100, got %d", config.BatchSize)
	}
	if config.BatchFlushInterval != 1*time.Second {
		t.Errorf("Expected default BatchFlushInterval 1s, got %v", config.BatchFlushInterval)
	}
	if config.MaxRetryAttempts != 3 {
		t.Errorf("Expected default MaxRetryAttempts 3, got %d", config.MaxRetryAttempts)
	}
	if config.InitialRetryDelay != 100*time.Millisecond {
		t.Errorf("Expected default InitialRetryDelay 100ms, got %v", config.InitialRetryDelay)
	}
	if config.MaxRetryDelay != 30*time.Second {
		t.Errorf("Expected default MaxRetryDelay 30s, got %v", config.MaxRetryDelay)
	}
	if config.CircuitBreakerThreshold != 5 {
		t.Errorf("Expected default CircuitBreakerThreshold 5, got %d", config.CircuitBreakerThreshold)
	}
	if config.CircuitBreakerTimeout != 60*time.Second {
		t.Errorf("Expected default CircuitBreakerTimeout 60s, got %v", config.CircuitBreakerTimeout)
	}
	if config.MaxMemoryUsage != 1024*1024*1024 {
		t.Errorf("Expected default MaxMemoryUsage 1GB, got %d", config.MaxMemoryUsage)
	}
	if config.CleanupThreshold != 0.8 {
		t.Errorf("Expected default CleanupThreshold 0.8, got %f", config.CleanupThreshold)
	}
	if config.GCInterval != 30*time.Second {
		t.Errorf("Expected default GCInterval 30s, got %v", config.GCInterval)
	}

	// Clean up
	os.Clearenv()
}

func TestGetEnv(t *testing.T) {
	// Test with environment variable set
	os.Setenv("TEST_VAR", "test_value")
	result := getEnv("TEST_VAR", "default")
	if result != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", result)
	}

	// Test with environment variable not set
	os.Unsetenv("TEST_VAR")
	result = getEnv("TEST_VAR", "default")
	if result != "default" {
		t.Errorf("Expected 'default', got '%s'", result)
	}

	// Clean up
	os.Clearenv()
}

func TestGetEnvInt(t *testing.T) {
	// Test with valid integer
	os.Setenv("TEST_INT", "42")
	result := getEnvInt("TEST_INT", 10)
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}

	// Test with invalid integer
	os.Setenv("TEST_INT", "invalid")
	result = getEnvInt("TEST_INT", 10)
	if result != 10 {
		t.Errorf("Expected 10, got %d", result)
	}

	// Test with environment variable not set
	os.Unsetenv("TEST_INT")
	result = getEnvInt("TEST_INT", 10)
	if result != 10 {
		t.Errorf("Expected 10, got %d", result)
	}

	// Clean up
	os.Clearenv()
}

func TestGetEnvInt64(t *testing.T) {
	// Test with valid int64
	os.Setenv("TEST_INT64", "9223372036854775807")
	result := getEnvInt64("TEST_INT64", 100)
	if result != 9223372036854775807 {
		t.Errorf("Expected 9223372036854775807, got %d", result)
	}

	// Test with invalid int64
	os.Setenv("TEST_INT64", "invalid")
	result = getEnvInt64("TEST_INT64", 100)
	if result != 100 {
		t.Errorf("Expected 100, got %d", result)
	}

	// Clean up
	os.Clearenv()
}

func TestGetEnvFloat64(t *testing.T) {
	// Test with valid float64
	os.Setenv("TEST_FLOAT", "3.14")
	result := getEnvFloat64("TEST_FLOAT", 1.0)
	if result != 3.14 {
		t.Errorf("Expected 3.14, got %f", result)
	}

	// Test with invalid float64
	os.Setenv("TEST_FLOAT", "invalid")
	result = getEnvFloat64("TEST_FLOAT", 1.0)
	if result != 1.0 {
		t.Errorf("Expected 1.0, got %f", result)
	}

	// Clean up
	os.Clearenv()
}

func TestGetEnvDuration(t *testing.T) {
	// Test with valid duration
	os.Setenv("TEST_DURATION", "5s")
	result := getEnvDuration("TEST_DURATION", 1*time.Second)
	if result != 5*time.Second {
		t.Errorf("Expected 5s, got %v", result)
	}

	// Test with invalid duration
	os.Setenv("TEST_DURATION", "invalid")
	result = getEnvDuration("TEST_DURATION", 1*time.Second)
	if result != 1*time.Second {
		t.Errorf("Expected 1s, got %v", result)
	}

	// Clean up
	os.Clearenv()
}
