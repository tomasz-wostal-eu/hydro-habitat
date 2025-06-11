package config

import (
	"os"
	"testing"
)

// TestLoad_DefaultValues checks whether the Load function correctly loads
// default values when no environment variables are set.
func TestLoad_DefaultValues(t *testing.T) {
	// Ensure that environment variables are clean for the test.
	// t.Setenv is the preferred method in Go 1.17+
	t.Setenv("DATABASE_URL", "")
	t.Setenv("API_PORT", "")
	t.Setenv("GIN_MODE", "")
	t.Setenv("JWT_SECRET", "")

	// Remove the .env file if it exists, so it doesn't affect the test
	_ = os.Remove(".env")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() function returned an error: %v", err)
	}

	// Expected default values
	expectedDBURL := "postgres://user:password@localhost:5432/hydro_habitat?sslmode=disable"
	expectedPort := "8080"
	expectedGinMode := "debug"
	expectedJWTSecret := "defaultsecret"

	if cfg.DatabaseURL != expectedDBURL {
		t.Errorf("expected DatabaseURL '%s', got '%s'", expectedDBURL, cfg.DatabaseURL)
	}
	if cfg.APIPort != expectedPort {
		t.Errorf("expected APIPort '%s', got '%s'", expectedPort, cfg.APIPort)
	}
	if cfg.GinMode != expectedGinMode {
		t.Errorf("expected GinMode '%s', got '%s'", expectedGinMode, cfg.GinMode)
	}
	if cfg.JWTSecret != expectedJWTSecret {
		t.Errorf("expected JWTSecret '%s', got '%s'", expectedJWTSecret, cfg.JWTSecret)
	}
}

// TestLoad_WithEnvVariables checks whether the Load function correctly overrides
// default values using values from environment variables.
func TestLoad_WithEnvVariables(t *testing.T) {
	// Set custom values for environment variables
	// t.Setenv automatically cleans up variables after the test finishes.
	expectedDBURL := "postgres://test:test@db:5432/testdb?sslmode=require"
	expectedPort := "9999"
	expectedGinMode := "release"
	expectedJWTSecret := "supersecretkey"

	t.Setenv("DATABASE_URL", expectedDBURL)
	t.Setenv("API_PORT", expectedPort)
	t.Setenv("GIN_MODE", expectedGinMode)
	t.Setenv("JWT_SECRET", expectedJWTSecret)

	// Remove the .env file if it exists, so it doesn't affect the test
	_ = os.Remove(".env")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() function returned an error: %v", err)
	}

	if cfg.DatabaseURL != expectedDBURL {
		t.Errorf("expected DatabaseURL '%s', got '%s'", expectedDBURL, cfg.DatabaseURL)
	}
	if cfg.APIPort != expectedPort {
		t.Errorf("expected APIPort '%s', got '%s'", expectedPort, cfg.APIPort)
	}
	if cfg.GinMode != expectedGinMode {
		t.Errorf("expected GinMode '%s', got '%s'", expectedGinMode, cfg.GinMode)
	}
	if cfg.JWTSecret != expectedJWTSecret {
		t.Errorf("expected JWTSecret '%s', got '%s'", expectedJWTSecret, cfg.JWTSecret)
	}
}

// TestLoad_WithDotEnvFile checks whether configuration is correctly
// loaded from the .env file.
func TestLoad_WithDotEnvFile(t *testing.T) {
	// Prepare the content of the .env file
	envContent := `
API_PORT=4000
GIN_MODE=test
# DATABASE_URL variable intentionally omitted to test fallback
`
	// Create a temporary .env file
	err := os.WriteFile(".env", []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create .env file: %v", err)
	}
	// Register a cleanup function that will remove the .env file after the test
	t.Cleanup(func() {
		_ = os.Remove(".env")
	})

	// Also check the precedence of environment variable over .env file
	// `godotenv` by default does not override existing environment variables.
	expectedJWTSecret := "secret_from_env_variable"
	t.Setenv("JWT_SECRET", expectedJWTSecret)
	t.Setenv("DATABASE_URL", "") // Clear the variable to use the default

	cfg, loadErr := Load()
	if loadErr != nil {
		t.Fatalf("Load() function returned an error: %v", loadErr)
	}

	// Check the values
	if cfg.APIPort != "4000" {
		t.Errorf("expected APIPort '4000' from .env file, got '%s'", cfg.APIPort)
	}
	if cfg.GinMode != "test" {
		t.Errorf("expected GinMode 'test' from .env file, got '%s'", cfg.GinMode)
	}
	if cfg.JWTSecret != expectedJWTSecret {
		t.Errorf("expected JWTSecret '%s' from environment variable, got '%s'", expectedJWTSecret, cfg.JWTSecret)
	}
	// Check if DatabaseURL used the default value
	expectedDBURL := "postgres://user:password@localhost:5432/hydro_habitat?sslmode=disable"
	if cfg.DatabaseURL != expectedDBURL {
		t.Errorf("expected default DatabaseURL '%s', got '%s'", expectedDBURL, cfg.DatabaseURL)
	}
}

// TestGetEnv tests the helper function getEnv.
func TestGetEnv(t *testing.T) {
	t.Run("variable exists", func(t *testing.T) {
		key := "TEST_VARIABLE_EXISTS"
		expectedValue := "hello world"
		t.Setenv(key, expectedValue)

		value := getEnv(key, "fallback")
		if value != expectedValue {
			t.Errorf("expected value '%s', got '%s'", expectedValue, value)
		}
	})

	t.Run("variable does not exist - use fallback", func(t *testing.T) {
		key := "TEST_VARIABLE_DOES_NOT_EXIST"
		fallbackValue := "default value"

		// Ensure that the variable is not set
		t.Setenv(key, "")
		_ = os.Unsetenv(key)

		value := getEnv(key, fallbackValue)
		if value != fallbackValue {
			t.Errorf("expected fallback value '%s', got '%s'", fallbackValue, value)
		}
	})
}
