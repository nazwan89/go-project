package utils

import (
	"os"
	"time"
)

// AppConfig holds all application configuration loaded from environment variables.
// LoadConfig() is the single point of truth for env-driven configuration.
type AppConfig struct {
	AppName  string
	Port     string
	Timezone string
	Location *time.Location
}

// LoadConfig reads configuration from environment variables with safe defaults.
// Defaults: AppName="go-project", Port="8080", Timezone="UTC".
// If TIMEZONE is set to an invalid location string, UTC is used silently.
func LoadConfig() AppConfig {
	appName := os.Getenv("APP_NAME")
	if appName == "" {
		appName = "go-project"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	tz := os.Getenv("TIMEZONE")
	if tz == "" {
		tz = "UTC"
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		loc = time.UTC
	}

	return AppConfig{
		AppName:  appName,
		Port:     port,
		Timezone: tz,
		Location: loc,
	}
}
