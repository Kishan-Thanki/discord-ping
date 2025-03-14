package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func ReadConfig() error {
	fmt.Println("Reading .env file...")

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return err
	}

	Token = os.Getenv("TOKEN")
	BotPrefix = os.Getenv("BOT_PREFIX")

	if Token == "" || BotPrefix == "" {
		return fmt.Errorf("missing required environment variables")
	}

	fmt.Println("Config loaded successfully!")
	return nil
}

var (
	Token     string
	BotPrefix string
)
