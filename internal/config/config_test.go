package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// Save original env vars
	originalEnv := make(map[string]string)
	envVars := []string{
		"SERVER_PORT", "SERVER_READ_TIMEOUT", "SERVER_WRITE_TIMEOUT",
		"SERVER_IDLE_TIMEOUT", "ENVIRONMENT", "DB_HOST", "DB_PORT",
		"DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE",
		"LOG_LEVEL", "LOG_FORMAT",
	}

	for _, key := range envVars {
		originalEnv[key] = os.Getenv(key)
		os.Unsetenv(key)
	}

	// Restore env vars after test
	defer func() {
		for key, value := range originalEnv {
			if value != "" {
				os.Setenv(key, value)
			}
		}
	}()

	t.Run("default configuration", func(t *testing.T) {
		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}

		if cfg.Server.Port != "8080" {
			t.Errorf("Expected default port 8080, got %s", cfg.Server.Port)
		}

		if cfg.Database.Host != "localhost" {
			t.Errorf("Expected default host localhost, got %s", cfg.Database.Host)
		}

		if cfg.Logger.Level != "info" {
			t.Errorf("Expected default log level info, got %s", cfg.Logger.Level)
		}
	})

	t.Run("custom configuration from env", func(t *testing.T) {
		os.Setenv("SERVER_PORT", "9090")
		os.Setenv("DB_HOST", "db.example.com")
		os.Setenv("DB_PORT", "5433")
		os.Setenv("LOG_LEVEL", "debug")
		os.Setenv("ENVIRONMENT", "production")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}

		if cfg.Server.Port != "9090" {
			t.Errorf("Expected port 9090, got %s", cfg.Server.Port)
		}

		if cfg.Database.Host != "db.example.com" {
			t.Errorf("Expected host db.example.com, got %s", cfg.Database.Host)
		}

		if cfg.Database.Port != "5433" {
			t.Errorf("Expected port 5433, got %s", cfg.Database.Port)
		}

		if cfg.Logger.Level != "debug" {
			t.Errorf("Expected log level debug, got %s", cfg.Logger.Level)
		}

		if cfg.Server.Environment != "production" {
			t.Errorf("Expected environment production, got %s", cfg.Server.Environment)
		}

		// Cleanup
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("ENVIRONMENT")
	})
}

func TestDatabaseConfig_ConnectionString(t *testing.T) {
	cfg := &DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "testuser",
		Password: "testpass",
		DBName:   "testdb",
		SSLMode:  "disable",
	}

	expected := "host=localhost port=5432 user=testuser password=testpass dbname=testdb sslmode=disable"
	result := cfg.ConnectionString()

	if result != expected {
		t.Errorf("Expected connection string %s, got %s", expected, result)
	}
}

func TestGetDurationEnv(t *testing.T) {
	tests := []struct {
		name         string
		envKey       string
		envValue     string
		defaultValue time.Duration
		expected     time.Duration
	}{
		{
			name:         "valid duration",
			envKey:       "TEST_DURATION",
			envValue:     "5s",
			defaultValue: 10 * time.Second,
			expected:     5 * time.Second,
		},
		{
			name:         "invalid duration uses default",
			envKey:       "TEST_DURATION",
			envValue:     "invalid",
			defaultValue: 10 * time.Second,
			expected:     10 * time.Second,
		},
		{
			name:         "empty value uses default",
			envKey:       "TEST_DURATION",
			envValue:     "",
			defaultValue: 10 * time.Second,
			expected:     10 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.envKey, tt.envValue)
				defer os.Unsetenv(tt.envKey)
			}

			result := getDurationEnv(tt.envKey, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
