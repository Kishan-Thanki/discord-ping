package bot

import (
	"strconv"
	"time"

	"discord-ping/internal/config"
	"github.com/bwmarrin/discordgo"
)

const embedColor = 0x5865F2 // Discord Blurple

func newEmbed(title, description string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       embedColor,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "discord-ping | High-Performance Diagnostics",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

func (b *Bot) cmdPing(s *discordgo.Session, m *discordgo.MessageCreate) {
	start := time.Now()
	embed := newEmbed("🏓 Ping Metrics", "Measuring latency...")
	msg, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		return
	}
	
	// Calculate API Round Trip
	apiLatency := time.Since(start).Milliseconds()
	
	// Get WebSocket Heartbeat Latency
	heartbeat := s.HeartbeatLatency().Milliseconds()
	
	// Message transit time
	messageTransit := time.Since(m.Timestamp).Milliseconds()

	embed.Description = ""
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name: "🌐 WebSocket Heartbeat",
			Value: "**" + strconv.FormatInt(heartbeat, 10) + "ms**\n*(Connection to Discord Gateway)*",
			Inline: false,
		},
		{
			Name: "⚡ API Round-Trip",
			Value: "**" + strconv.FormatInt(apiLatency, 10) + "ms**\n*(Time to send and receive message)*",
			Inline: false,
		},
		{
			Name: "📨 Message Transit",
			Value: "**" + strconv.FormatInt(messageTransit, 10) + "ms**\n*(Time since user sent command)*",
			Inline: false,
		},
	}
	
	EditEmbed(s, m.ChannelID, msg.ID, embed)
}

func (b *Bot) cmdVersion(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := newEmbed("ℹ️ About", "discord-ping is a hyper-optimized Discord diagnostic bot built in Go.")
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name: "Version", Value: "`" + config.Version + "`", Inline: true,
	})
	SendEmbed(s, m.ChannelID, embed)
}

func (b *Bot) cmdUptime(s *discordgo.Session, m *discordgo.MessageCreate) {
	uptime := time.Since(b.startTime)

	hours := int(uptime.Hours())
	minutes := int(uptime.Minutes()) % 60
	seconds := int(uptime.Seconds()) % 60

	embed := newEmbed("⏱️ Uptime", "Online and monitoring for **"+strconv.Itoa(hours)+"h "+strconv.Itoa(minutes)+"m "+strconv.Itoa(seconds)+"s**")
	SendEmbed(s, m.ChannelID, embed)
}
