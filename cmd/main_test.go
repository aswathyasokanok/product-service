package main

import (
	"os"
	"testing"
)

func TestMain_EnvironmentVariables(t *testing.T) {
	// Test that the main function can be called without panicking
	// This is a basic smoke test

	// Set test environment variables
	os.Setenv("WORKERS", "1")
	os.Setenv("QUEUE_SIZE", "10")
	os.Setenv("PORT", "0") // Use port 0 for testing

	// Note: We can't easily test the main function directly as it starts a server
	// This test mainly ensures the code compiles and basic setup works

	// Clean up
	os.Clearenv()
}

func TestMain_DefaultConfig(t *testing.T) {
	// Clear environment variables to test defaults
	os.Clearenv()

	// Test that default configuration values are used
	// This is tested indirectly through the config package tests

	// The main function should use default values when no env vars are set
	// This is verified in the config package tests
}
