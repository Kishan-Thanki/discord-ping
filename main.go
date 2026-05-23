package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Kishan-Thanki/discord-ping/bot"
	"github.com/Kishan-Thanki/discord-ping/config"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	err := godotenv.Load()
	if err != nil {
		slog.Warn("No .env file found, ensure environment variables are set")
	}

	err = config.ReadConfig()
	if err != nil {
		slog.Error("Failed to read config", "error", err)
		return
	}

	slog.Info("Starting go-discord-ping", "version", config.Version)
	bot.Start()

	slog.Info("Press CTRL-C to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	bot.Stop()
}
