package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig_Defaults(t *testing.T) {
	os.Unsetenv("APP_NAME")
	os.Unsetenv("PORT")
	os.Unsetenv("TIMEZONE")

	cfg := LoadConfig()

	if cfg.AppName != "go-project" {
		t.Errorf("expected AppName 'go-project', got %q", cfg.AppName)
	}
	if cfg.Port != "8080" {
		t.Errorf("expected Port '8080', got %q", cfg.Port)
	}
	if cfg.Timezone != "UTC" {
		t.Errorf("expected Timezone 'UTC', got %q", cfg.Timezone)
	}
	if cfg.Location != time.UTC {
		t.Errorf("expected Location time.UTC, got %v", cfg.Location)
	}
}

func TestLoadConfig_FromEnv(t *testing.T) {
	os.Setenv("APP_NAME", "MyApp")
	os.Setenv("PORT", "9090")
	os.Setenv("TIMEZONE", "Asia/Kuala_Lumpur")
	defer func() {
		os.Unsetenv("APP_NAME")
		os.Unsetenv("PORT")
		os.Unsetenv("TIMEZONE")
	}()

	cfg := LoadConfig()

	if cfg.AppName != "MyApp" {
		t.Errorf("expected AppName 'MyApp', got %q", cfg.AppName)
	}
	if cfg.Port != "9090" {
		t.Errorf("expected Port '9090', got %q", cfg.Port)
	}
	if cfg.Timezone != "Asia/Kuala_Lumpur" {
		t.Errorf("expected Timezone 'Asia/Kuala_Lumpur', got %q", cfg.Timezone)
	}
	if cfg.Location == nil {
		t.Fatal("expected non-nil Location")
	}
}

func TestLoadConfig_InvalidTimezone(t *testing.T) {
	os.Setenv("TIMEZONE", "Invalid/Zone")
	defer os.Unsetenv("TIMEZONE")

	cfg := LoadConfig()

	if cfg.Location != time.UTC {
		t.Errorf("expected fallback to time.UTC for invalid timezone, got %v", cfg.Location)
	}
}
