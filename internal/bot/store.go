package bot

import (
	"context"
	"time"

	"discord-ping/internal/database"
)

// Store defines the data operations that the Bot depends on.
// The bot package does NOT need to know about SQL, prepared statements,
// or connection pools. It only cares about these behaviors.
//
// This is the Consumer-Defined Interface pattern.
type Store interface {
	GetUser(ctx context.Context, userID, guildID string) (*database.User, error)
	UpdateUserXP(ctx context.Context, userID, guildID string, xp, level int) error
	UpdateUserBalance(ctx context.Context, userID, guildID string, balance int) error
	UpdateUserDaily(ctx context.Context, userID, guildID string, balance int, lastDaily string) error
	GetLeaderboard(ctx context.Context, guildID string, limit int) ([]*database.User, error)
	GetPrefix(ctx context.Context, guildID string) (string, error)
	SetPrefix(ctx context.Context, guildID, prefix string) error
	CreateReminder(ctx context.Context, userID, channelID, guildID, message string, remindAt time.Time) (int, error)
	DeleteReminder(ctx context.Context, id int) error
	GetPendingReminders(ctx context.Context) ([]*database.Reminder, error)
	AddWarning(ctx context.Context, userID, guildID, reason string) error
	GetWarningCount(ctx context.Context, userID, guildID string) (int, error)
	Close(ctx context.Context)
}
