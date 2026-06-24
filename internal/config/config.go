package config

import (
	"fmt"
	"log/slog"
	"os"
)

const Version = "v1.1.0"

// Config represents the single source of truth for the bot's environment.
// Explicitly define every configuration value the bot needs here.
type Config struct {
	Token     string
	BotPrefix string
}

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
