package bot

import (
	"context"
	"log/slog"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == b.BotID || m.Author.Bot {
		return
	}

	prefix := b.cfg.BotPrefix
	if m.GuildID != "" {
		p, err := b.store.GetPrefix(context.Background(), m.GuildID)
		if err == nil && p != "" {
			prefix = p
		}
	}

	if strings.HasPrefix(m.Content, prefix) {
		if b.isRateLimited(s, m) {
			return
		}

		cmd := strings.TrimPrefix(m.Content, prefix)
		args := strings.Fields(cmd)
		if len(args) == 0 {
			return
		}

		slog.Info("Command received",
			"command", args[0],
			"user", m.Author.Username,
			"channel", m.ChannelID,
			"guild", m.GuildID,
		)

		switch strings.ToLower(args[0]) {
		case "ping":
			b.cmdPing(s, m)
		case "version", "about":
			b.cmdVersion(s, m)
		case "uptime":
			b.cmdUptime(s, m)
		case "help":
			b.cmdHelp(s, m)
		case "setprefix":
			b.cmdSetPrefix(s, m, args)
		}
	}
}

func (b *Bot) slashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	slog.Info("Slash command received", "command", i.ApplicationCommandData().Name)

	switch i.ApplicationCommandData().Name {
	case "ping":
		// Quick response for slash command since we don't have message timing
		latency := s.HeartbeatLatency().Milliseconds()
		embed := newEmbed("🏓 Ping Metrics", "WebSocket Heartbeat: **"+strconv.FormatInt(latency, 10)+"ms**")

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		})
		if err != nil {
			slog.Error("Failed to respond to slash command", "error", err)
		}
	}
}
