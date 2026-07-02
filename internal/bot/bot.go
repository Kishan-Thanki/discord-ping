package bot

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"discord-ping/internal/config"

	"github.com/bwmarrin/discordgo"
)

type Store interface {
	GetPrefix(ctx context.Context, guildID string) (string, error)
	SetPrefix(ctx context.Context, guildID, prefix string) error
	Close(ctx context.Context)
}

type Bot struct {
	cfg         *config.Config
	store       Store
	goBot       *discordgo.Session
	BotID       string
	rateLimits  sync.Map
	stopCleanup chan struct{}
}

func NewBot(cfg *config.Config, store Store) *Bot {
	return &Bot{
		cfg:         cfg,
		store:       store,
		stopCleanup: make(chan struct{}),
	}
}

var slashCommands = []*discordgo.ApplicationCommand{
	{
		Name:        "ping",
		Description: "Replies with pong and network latency",
	},
}

func (b *Bot) Start() error {
	var err error
	b.goBot, err = discordgo.New("Bot " + b.cfg.Token)
	if err != nil {
		return fmt.Errorf("creating discord session: %w", err)
	}

	u, err := b.goBot.User("@me")
	if err != nil {
		return fmt.Errorf("fetching bot user: %w", err)
	}

	b.BotID = u.ID

	b.goBot.Identify.Intents = discordgo.IntentsGuildMessages |
		discordgo.IntentsDirectMessages |
		discordgo.IntentsMessageContent |
		discordgo.IntentsGuilds

	b.goBot.StateEnabled = false

	b.goBot.AddHandler(b.messageHandler)
	b.goBot.AddHandler(b.slashCommandHandler)

	if err := b.goBot.Open(); err != nil {
		return fmt.Errorf("opening discord connection: %w", err)
	}

	_ = b.goBot.UpdateListeningStatus(b.cfg.BotPrefix + "ping")

	for _, cmd := range slashCommands {
		if _, err := b.goBot.ApplicationCommandCreate(b.BotID, "", cmd); err != nil {
			slog.Error("Failed to register slash command", "command", cmd.Name, "error", err)
		}
	}

	go b.startRateLimitCleaner()

	slog.Info("Bot is running", "user", u.Username, "id", u.ID)
	return nil
}

func (b *Bot) Stop() {
	close(b.stopCleanup)
	if b.goBot != nil {
		slog.Info("Shutting down bot gracefully")
		b.goBot.Close()
	}
	b.store.Close(context.Background())
}
