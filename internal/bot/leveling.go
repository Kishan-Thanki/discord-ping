package bot

import (
	"context"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) handlePassiveXP(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.GuildID == "" {
		return
	}

	user, err := b.store.GetUser(context.Background(), m.Author.ID, m.GuildID)
	if err != nil {
		return
	}

	xpGained := rand.Intn(11) + 15 // 15 to 25
	newXP := user.XP + xpGained
	newLevel := int(math.Floor(math.Sqrt(float64(newXP) / 100.0)))

	_ = b.store.UpdateUserXP(context.Background(), m.Author.ID, m.GuildID, newXP, newLevel)

	if newLevel > user.Level {
		embed := newEmbed(
			"🎉 Level Up!",
			"Congratulations <@"+m.Author.ID+">! You've reached **Level "+strconv.Itoa(newLevel)+"**!",
		)
		SendEmbed(s, m.ChannelID, embed)
	}
}

func (b *Bot) cmdRank(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.GuildID == "" {
		return
	}

	user, err := b.store.GetUser(context.Background(), m.Author.ID, m.GuildID)
	if err != nil {
		return
	}

	embed := newEmbed("🏆 Your Rank", "")
	embed.Fields = []*discordgo.MessageEmbedField{
		{Name: "Level", Value: strconv.Itoa(user.Level), Inline: true},
		{Name: "XP", Value: strconv.Itoa(user.XP), Inline: true},
	}
	embed.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: m.Author.AvatarURL("256")}

	SendEmbed(s, m.ChannelID, embed)
}

func (b *Bot) cmdLeaderboard(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.GuildID == "" {
		return
	}

	leaders, err := b.store.GetLeaderboard(context.Background(), m.GuildID, 10)
	if err != nil || len(leaders) == 0 {
		embed := newEmbed("🏆 Leaderboard", "No data yet!")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	var desc strings.Builder
	for i, u := range leaders {
		desc.WriteString("**#")
		desc.WriteString(strconv.Itoa(i + 1))
		desc.WriteString("** <@")
		desc.WriteString(u.UserID)
		desc.WriteString("> - Level ")
		desc.WriteString(strconv.Itoa(u.Level))
		desc.WriteString(" (")
		desc.WriteString(strconv.Itoa(u.XP))
		desc.WriteString(" XP)\n")
	}

	embed := newEmbed("🏆 Server Leaderboard", desc.String())
	SendEmbed(s, m.ChannelID, embed)
}
