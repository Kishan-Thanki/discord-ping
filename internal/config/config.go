package config

import (
	"fmt"
	"log/slog"
	"os"
)

// Version is the application version. It is a constant because it never changes at runtime.
const Version = "v1.0.0"

// Config represents the single source of truth for the bot's environment.
// We explicitly define every configuration value the bot needs here.
type Config struct {
	Token     string
	BotPrefix string
}

// Load reads from environment variables and constructs the typed Config struct.
// It fails fast if required variables are missing.
func Load() (*Config, error) {
	token := os.Getenv("TOKEN")
	if token == "" {
		return nil, fmt.Errorf("required environment variable TOKEN is not set")
	}

	prefix := os.Getenv("BOT_PREFIX")
	if prefix == "" {
		return nil, fmt.Errorf("required environment variable BOT_PREFIX is not set")
	}

	slog.Info("Config loaded successfully")
	return &Config{
		Token:     token,
		BotPrefix: prefix,
	}, nil
}
