package bot

import (
	"log/slog"
	"strings"
	"time"

	"github.com/Kishan-Thanki/discord-ping/database"
	"github.com/bwmarrin/discordgo"
)

func LoadReminders(s *discordgo.Session) {
	reminders, err := database.GetPendingReminders()
	if err != nil {
		slog.Error("Failed to load reminders from DB", "error", err)
		return
	}

	for _, r := range reminders {
		delay := time.Until(r.RemindAt)
		if delay <= 0 {
			// Reminder is past due, fire immediately
			fireReminder(s, r.ID, r.UserID, r.Message)
		} else {
			// Schedule it
			scheduleReminder(s, r.ID, r.UserID, r.Message, delay)
		}
	}
	slog.Info("Loaded pending reminders", "count", len(reminders))
}

func scheduleReminder(s *discordgo.Session, id int, userID, message string, delay time.Duration) {
	time.AfterFunc(delay, func() {
		fireReminder(s, id, userID, message)
	})
}

func fireReminder(s *discordgo.Session, id int, userID, message string) {
	_ = database.DeleteReminder(id)

	ch, err := s.UserChannelCreate(userID)
	if err != nil {
		slog.Warn("Failed to create DM channel for reminder", "user", userID, "error", err)
		return
	}

	embed := newEmbed("⏰ Reminder", message)
	SendEmbed(s, ch.ID, embed)
}

func cmdRemind(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 3 {
		embed := newEmbed("❌ Usage", "Usage: `!remind <duration> <message>` (e.g. `!remind 30m Check the oven`)")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	durationStr := args[1]
	message := strings.Join(args[2:], " ")

	d, err := time.ParseDuration(durationStr)
	if err != nil || d <= 0 {
		embed := newEmbed("❌ Error", "Invalid duration format. Use formats like `30s`, `5m`, `2h`, `1d`.")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	remindAt := time.Now().Add(d)
	guildID := m.GuildID
	if guildID == "" {
		guildID = "DM"
	}

	id, err := database.CreateReminder(m.Author.ID, m.ChannelID, guildID, message, remindAt)
	if err != nil {
		slog.Error("Failed to save reminder to DB", "error", err)
		embed := newEmbed("❌ Error", "Failed to set reminder.")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	scheduleReminder(s, id, m.Author.ID, message, d)

	embed := newEmbed("⏰ Reminder Set", "I will remind you in **"+durationStr+"**.")
	SendEmbed(s, m.ChannelID, embed)
}
