package bot

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"discord-ping/internal/config"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/image/font"
)

// Bot is the central struct that holds all runtime state.
// There are ZERO package-level mutable variables — everything lives here.
type Bot struct {
	cfg       *config.Config
	store     Store
	goBot     *discordgo.Session
	BotID     string
	startTime time.Time

	blackjackGames sync.Map
	wordleGames    sync.Map
	triviaActive   sync.Map
	rateLimits     sync.Map

	fontFace48 font.Face
	fontFace32 font.Face
}

// NewBot constructs a Bot with its injected dependencies.
func NewBot(cfg *config.Config, store Store) *Bot {
	return &Bot{
		cfg:   cfg,
		store: store,
	}
}

var slashCommands = []*discordgo.ApplicationCommand{
	{
		Name:        "ping",
		Description: "Replies with pong and network latency",
	},
}

// Start initializes the Discord session and connects to the gateway.
// It returns an error if any critical step fails — the caller decides what to do.
func (b *Bot) Start() error {
	b.startTime = time.Now()

	b.initFonts()

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
		discordgo.IntentsGuildMembers |
		discordgo.IntentsGuilds

	b.goBot.StateEnabled = false

	b.goBot.AddHandler(b.messageHandler)
	b.goBot.AddHandler(b.slashCommandHandler)
	b.goBot.AddHandler(b.welcomeHandler)

	if err := b.goBot.Open(); err != nil {
		return fmt.Errorf("opening discord connection: %w", err)
	}

	b.LoadReminders(b.goBot)

	_ = b.goBot.UpdateListeningStatus(b.cfg.BotPrefix + "ping")

	for _, cmd := range slashCommands {
		if _, err := b.goBot.ApplicationCommandCreate(b.goBot.State.User.ID, "", cmd); err != nil {
			slog.Error("Failed to register slash command", "command", cmd.Name, "error", err)
		}
	}

	slog.Info("Bot is running", "user", u.Username, "id", u.ID)
	return nil
}

// Stop gracefully shuts down the bot and closes the database.
func (b *Bot) Stop() {
	if b.goBot != nil {
		slog.Info("Shutting down bot gracefully")
		b.goBot.Close()
	}
	b.store.Close(context.Background())
}
