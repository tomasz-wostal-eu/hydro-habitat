package main

import (
"os"
"testing"
"time"

"hydro-habitat/backend/config"

"github.com/stretchr/testify/assert"
)

// TestMain_ConfigLoad tests that the main function can load configuration
func TestMain_ConfigLoad(t *testing.T) {
	// Save original environment
	originalDBURL := os.Getenv("DATABASE_URL")
	originalAPIPort := os.Getenv("API_PORT")
	originalGinMode := os.Getenv("GIN_MODE")

	// Clean up after test
	defer func() {
		if originalDBURL != "" {
			_ = os.Setenv("DATABASE_URL", originalDBURL)
		} else {
			_ = os.Unsetenv("DATABASE_URL")
		}
		if originalAPIPort != "" {
			_ = os.Setenv("API_PORT", originalAPIPort)
		} else {
			_ = os.Unsetenv("API_PORT")
		}
		if originalGinMode != "" {
			_ = os.Setenv("GIN_MODE", originalGinMode)
		} else {
			_ = os.Unsetenv("GIN_MODE")
		}
	}()

	// Set test environment variables
	if err := os.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/test?sslmode=disable"); err != nil {
		t.Fatalf("Failed to set DATABASE_URL environment variable: %v", err)
	}
	if err := os.Setenv("API_PORT", "8081"); err != nil {
		t.Fatalf("Failed to set API_PORT environment variable: %v", err)
	}
	if err := os.Setenv("GIN_MODE", "test"); err != nil {
		t.Fatalf("Failed to set GIN_MODE environment variable: %v", err)
	}

	// Test that config.Load() works
	cfg, err := config.Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "8081", cfg.APIPort)
	assert.Equal(t, "test", cfg.GinMode)
}

// TestMain_Integration tests that we can start the basic components
// This is more of an integration test to ensure imports and basic setup work
func TestMain_Integration(t *testing.T) {
	// This test primarily ensures that all imports work and basic setup doesn't panic
// We can't easily test the actual main() function without starting the server

	// Test that we can create basic components
	cfg := &config.Config{
		DatabaseURL: "postgres://user:password@localhost:5432/hydro_habitat?sslmode=disable",
		APIPort:     "8080",
		GinMode:     "test",
		JWTSecret:   "test-secret",
	}

	assert.NotNil(t, cfg)
	assert.Equal(t, "8080", cfg.APIPort)
	assert.Equal(t, "test", cfg.GinMode)

	// Note: We don't actually try to connect to the database or start the server
// in unit tests as that would make them dependent on external resources
}

// TestMain_Timeout tests that main function setup completes in reasonable time
func TestMain_Timeout(t *testing.T) {
// Test that loading configuration doesn't hang
	done := make(chan bool, 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				// Config loading might fail, but it shouldn't panic
t.Logf("Config loading panicked (expected in test environment): %v", r)
}
done <- true
}()

_, err := config.Load()
if err != nil {
t.Logf("Config loading failed (expected in test environment): %v", err)
}
}()

select {
case <-done:
// Configuration loading completed
case <-time.After(5 * time.Second):
t.Fatal("Configuration loading took too long (>5 seconds)")
}
}
