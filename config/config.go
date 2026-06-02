package config

import (
	"errors"
	"log/slog"
	"os"
)

var (
	Token     string
	BotPrefix string
	Version   = "v1.0.0"
)

func ReadConfig() error {
	Token = os.Getenv("TOKEN")
	BotPrefix = os.Getenv("BOT_PREFIX")

	if Token == "" || BotPrefix == "" {
		return errors.New("missing required environment variables (TOKEN, BOT_PREFIX)")
	}

	slog.Info("Config loaded successfully")
	return nil
}
