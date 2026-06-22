package bot

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/Kishan-Thanki/discord-ping/internal/config"
	"github.com/bwmarrin/discordgo"
)

const embedColor = 0x5865F2 // Discord Blurple

func newEmbed(title, description string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       embedColor,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "ping-bot",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

func (b *Bot) cmdPing(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := newEmbed("🏓 Pong!", "Latency calculation...")
	msg, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		return
	}
	latency := time.Since(m.Timestamp).Milliseconds()
	embed.Description = "Latency: **" + strconv.FormatInt(latency, 10) + "ms**"
	EditEmbed(s, m.ChannelID, msg.ID, embed)
}

func (b *Bot) cmdVersion(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := newEmbed("ℹ️ About", "go-discord-ping is a hyper-optimized Discord utility bot built in Go.")
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

	embed := newEmbed("⏱️ Uptime", "I've been running for **"+strconv.Itoa(hours)+"h "+strconv.Itoa(minutes)+"m "+strconv.Itoa(seconds)+"s**")
	SendEmbed(s, m.ChannelID, embed)
}

func (b *Bot) cmdServerInfo(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.GuildID == "" {
		embed := newEmbed("❌ Error", "This command can only be used in a server.")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	guild, err := s.Guild(m.GuildID)
	if err != nil {
		embed := newEmbed("❌ Error", "Failed to fetch server info.")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	createdAt := snowflakeToTime(guild.ID)

	embed := newEmbed("📊 Server Info", "")
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "Name",
			Value:  guild.Name,
			Inline: true,
		},
		{
			Name:   "Owner",
			Value:  "<@" + guild.OwnerID + ">",
			Inline: true,
		},
		{
			Name:   "Members",
			Value:  strconv.Itoa(guild.MemberCount),
			Inline: true,
		},
		{
			Name:   "Channels",
			Value:  strconv.Itoa(len(guild.Channels)),
			Inline: true,
		},
		{
			Name:   "Roles",
			Value:  strconv.Itoa(len(guild.Roles)),
			Inline: true,
		},
		{
			Name:   "Created",
			Value:  "<t:" + strconv.FormatInt(createdAt.Unix(), 10) + ":R>",
			Inline: true,
		},
	}

	if guild.Icon != "" {
		iconURL := "https://cdn.discordapp.com/icons/" + guild.ID + "/" + guild.Icon + ".png?size=256"
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: iconURL,
		}
	}

	SendEmbed(s, m.ChannelID, embed)
}

func (b *Bot) cmdAvatar(s *discordgo.Session, m *discordgo.MessageCreate) {
	var targetUser *discordgo.User

	if len(m.Mentions) > 0 {
		targetUser = m.Mentions[0]
	} else {
		targetUser = m.Author
	}

	avatarURL := targetUser.AvatarURL("512")

	embed := newEmbed(
		"🖼️ "+targetUser.Username+"'s Avatar",
		"",
	)
	embed.Image = &discordgo.MessageEmbedImage{
		URL: avatarURL,
	}

	SendEmbed(s, m.ChannelID, embed)
}

func (b *Bot) cmdRoll(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	maxVal := 6

	if len(args) > 1 {
		parsed, err := strconv.Atoi(args[1])
		if err == nil && parsed > 1 {
			maxVal = parsed
		}
	}

	result := rand.Intn(maxVal) + 1

	embed := newEmbed("🎲 Dice Roll", "You rolled a **"+strconv.Itoa(result)+"** (1-"+strconv.Itoa(maxVal)+")")
	SendEmbed(s, m.ChannelID, embed)
}

func snowflakeToTime(id string) time.Time {
	snowflake, _ := strconv.ParseInt(id, 10, 64)
	timestamp := (snowflake >> 22) + 1420070400000
	return time.UnixMilli(timestamp)
}
