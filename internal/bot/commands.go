package bot

import (
	"strconv"
	"time"

	"discord-ping/internal/config"

	"github.com/bwmarrin/discordgo"
)

const embedColor = 0x5865F2

func newEmbed(title, description string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       embedColor,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "discord-ping",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

func (b *Bot) cmdPing(s *discordgo.Session, ctx PingContext) {
	start := time.Now()
	embed := newEmbed("🏓 Ping Metrics", "Measuring latency...")
	if err := ctx.Respond(embed); err != nil {
		return
	}

	apiLatency := time.Since(start).Milliseconds()
	heartbeat := s.HeartbeatLatency().Milliseconds()
	transitTime := time.Since(ctx.SentAt()).Milliseconds()

	embed.Title = "🏓 Pong! **" + strconv.FormatInt(heartbeat, 10) + "ms**"
	embed.Description = ""
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "🌐 WebSocket",
			Value:  "*" + strconv.FormatInt(heartbeat, 10) + "ms*",
			Inline: true,
		},
		{
			Name:   "⚡ API Round-Trip",
			Value:  "*" + strconv.FormatInt(apiLatency, 10) + "ms*",
			Inline: true,
		},
		{
			Name:   "📨 Transit",
			Value:  "*" + strconv.FormatInt(transitTime, 10) + "ms*",
			Inline: true,
		},
	}

	_ = ctx.Respond(embed)
}

func (b *Bot) cmdVersion(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := newEmbed("ℹ️ About", "discord-ping is a diagnostic bot built in Go.")
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name: "Version", Value: "`" + config.Version + "`", Inline: true,
	})
	SendEmbed(s, m.ChannelID, embed)
}
