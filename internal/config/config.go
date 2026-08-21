package config

import (
	"fmt"
	"os"
	"strings"
)

const Version = "v1.2.0"

type Config struct {
	Token     string
	BotPrefix string
}

func Load() (*Config, error) {
	token := strings.TrimSpace(os.Getenv("TOKEN"))
	if token == "" {
		return nil, fmt.Errorf("required environment variable TOKEN is not set or empty")
	}

	prefix := strings.TrimSpace(os.Getenv("BOT_PREFIX"))
	if prefix == "" {
		return nil, fmt.Errorf("required environment variable BOT_PREFIX is not set or empty")
	}

	return &Config{
		Token:     token,
		BotPrefix: prefix,
	}, nil
}
