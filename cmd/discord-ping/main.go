package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Kishan-Thanki/discord-ping/internal/bot"
	"github.com/Kishan-Thanki/discord-ping/internal/config"
	"github.com/Kishan-Thanki/discord-ping/internal/database"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		return
	}

	slog.Info("Starting go-discord-ping", "version", config.Version)

	repo, err := database.NewRepository("bot.db")
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		return
	}

	b := bot.NewBot(cfg, repo)
	if err := b.Start(); err != nil {
		slog.Error("Failed to start bot", "error", err)
		return
	}

	slog.Info("Press CTRL-C to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	b.Stop()
}
