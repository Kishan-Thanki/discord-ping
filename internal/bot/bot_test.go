package bot

import (
	"context"
	"fmt"
	"testing"
	"time"

	"discord-ping/internal/config"
	"github.com/bwmarrin/discordgo"
)

type mockStore struct {
	prefixes map[string]string
	closed   bool
}

func (m *mockStore) GetPrefix(ctx context.Context, guildID string) (string, error) {
	p, ok := m.prefixes[guildID]
	if !ok {
		return "", fmt.Errorf("no prefix found")
	}
	return p, nil
}

func (m *mockStore) SetPrefix(ctx context.Context, guildID, prefix string) error {
	m.prefixes[guildID] = prefix
	return nil
}

func (m *mockStore) Close(ctx context.Context) {
	m.closed = true
}

func TestNewBot(t *testing.T) {
	cfg := &config.Config{
		Token:     "dummy-token",
		BotPrefix: "?",
	}
	store := &mockStore{prefixes: make(map[string]string)}
	b := NewBot(cfg, store)

	if b.cfg.BotPrefix != "?" {
		t.Errorf("expected BotPrefix '?', got %q", b.cfg.BotPrefix)
	}
	if b.stopCleanup == nil {
		t.Error("expected stopCleanup channel to be initialized")
	}
}

func TestRateLimiterAndCleaner(t *testing.T) {
	cfg := &config.Config{
		Token:     "dummy-token",
		BotPrefix: "!",
	}
	store := &mockStore{prefixes: make(map[string]string)}
	b := NewBot(cfg, store)
	b.BotID = "bot-123"

	sess, err := discordgo.New("Bot test-token")
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	sess.StateEnabled = false

	msg := &discordgo.MessageCreate{
		Message: &discordgo.Message{
			Author: &discordgo.User{
				ID: "user-456",
			},
			ChannelID: "channel-789",
		},
	}

	limited := b.isRateLimited(sess, msg)
	if limited {
		t.Error("expected first call to not be rate limited")
	}

	limited = b.isRateLimited(sess, msg)
	if !limited {
		t.Error("expected second immediate call to be rate limited")
	}

	b.rateLimits.Store("user-expired", time.Now().Add(-10*time.Second))
	b.rateLimits.Store("user-active", time.Now())

	now := time.Now()
	b.rateLimits.Range(func(key, val interface{}) bool {
		lastCmd := val.(time.Time)
		if now.Sub(lastCmd) > rateLimitDuration {
			b.rateLimits.Delete(key)
		}
		return true
	})

	_, activeExists := b.rateLimits.Load("user-active")
	if !activeExists {
		t.Error("expected active user to still be in rate limits")
	}

	_, expiredExists := b.rateLimits.Load("user-expired")
	if expiredExists {
		t.Error("expected expired user to be cleaned up")
	}
}

func TestHasPermissionStopsPanic(t *testing.T) {
	sess, err := discordgo.New("Bot test-token")
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	sess.StateEnabled = false
	sess.State = nil

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("hasPermission panicked: %v", r)
		}
	}()

	perms := hasPermission(sess, "user-1", "channel-1", discordgo.PermissionAdministrator)
	if perms {
		t.Error("expected hasPermission to return false on default session")
	}
}
