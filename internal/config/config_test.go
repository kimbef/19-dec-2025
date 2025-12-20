package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// Set environment variables
	os.Setenv("PORT", "9090")
	os.Setenv("DB_HOST", "testhost")
	os.Setenv("LOG_LEVEL", "debug")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.Server.Port != "9090" {
		t.Errorf("expected port 9090, got %s", cfg.Server.Port)
	}

	if cfg.Database.Host != "testhost" {
		t.Errorf("expected host testhost, got %s", cfg.Database.Host)
	}

	if cfg.Log.Level != "debug" {
		t.Errorf("expected log level debug, got %s", cfg.Log.Level)
	}

	// Clean up
	os.Unsetenv("PORT")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("LOG_LEVEL")
}

func TestLoadDefaults(t *testing.T) {
	// Clear relevant env vars
	os.Clearenv()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.Server.Port != "8080" {
		t.Errorf("expected default port 8080, got %s", cfg.Server.Port)
	}

	if cfg.Database.Host != "localhost" {
		t.Errorf("expected default host localhost, got %s", cfg.Database.Host)
	}

	if cfg.Log.Level != "info" {
		t.Errorf("expected default log level info, got %s", cfg.Log.Level)
	}
}

func TestGetDSN(t *testing.T) {
	dbCfg := DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "testuser",
		Password: "testpass",
		DBName:   "testdb",
		SSLMode:  "disable",
	}

	expected := "host=localhost port=5432 user=testuser password=testpass dbname=testdb sslmode=disable"
	if dsn := dbCfg.GetDSN(); dsn != expected {
		t.Errorf("expected DSN %s, got %s", expected, dsn)
	}
}

func TestGetIntEnv(t *testing.T) {
	os.Setenv("TEST_INT", "42")
	defer os.Unsetenv("TEST_INT")

	result := getIntEnv("TEST_INT", 10)
	if result != 42 {
		t.Errorf("expected 42, got %d", result)
	}

	result = getIntEnv("NONEXISTENT", 10)
	if result != 10 {
		t.Errorf("expected default 10, got %d", result)
	}
}

func TestGetDurationEnv(t *testing.T) {
	os.Setenv("TEST_DURATION", "5m")
	defer os.Unsetenv("TEST_DURATION")

	result := getDurationEnv("TEST_DURATION", 1*time.Minute)
	if result != 5*time.Minute {
		t.Errorf("expected 5m, got %v", result)
	}

	result = getDurationEnv("NONEXISTENT", 1*time.Minute)
	if result != 1*time.Minute {
		t.Errorf("expected default 1m, got %v", result)
	}
}
