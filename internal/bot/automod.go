package bot

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var badWordsRegex = regexp.MustCompile(`(?i)\b(fuck|shit|bitch|asshole|cunt|dick|faggot|nigger|kys|discord\.gift|free nitro)\b`)

func (b *Bot) checkBadWords(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	if m.GuildID == "" || m.Author.Bot {
		return false
	}

	if badWordsRegex.MatchString(m.Content) {
		_ = s.ChannelMessageDelete(m.ChannelID, m.ID)

		b.applyWarning(s, m.GuildID, m.ChannelID, m.Author.ID, "Triggered auto-mod filter (bad words/spam)")
		return true
	}
	return false
}

func (b *Bot) cmdWarn(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if m.GuildID == "" {
		return
	}

	if !hasPermission(s, m.Author.ID, m.ChannelID, discordgo.PermissionKickMembers) {
		embed := newEmbed("🚫 Permission Denied", "You don't have permission to warn members.")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	if len(m.Mentions) == 0 {
		embed := newEmbed("❌ Usage", "Usage: `!warn @user [reason]`")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	target := m.Mentions[0]
	if target.ID == m.Author.ID {
		embed := newEmbed("❌ Error", "You can't warn yourself.")
		SendEmbed(s, m.ChannelID, embed)
		return
	}
	if target.Bot {
		embed := newEmbed("❌ Error", "You can't warn a bot.")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	reason := "No reason provided"
	if len(args) > 2 {
		reason = strings.Join(args[2:], " ")
	}

	b.applyWarning(s, m.GuildID, m.ChannelID, target.ID, reason)
}

func (b *Bot) cmdWarnings(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.GuildID == "" {
		return
	}

	target := m.Author
	if len(m.Mentions) > 0 {
		target = m.Mentions[0]
	}

	count, err := b.store.GetWarningCount(context.Background(), target.ID, m.GuildID)
	if err != nil {
		count = 0
	}

	embed := newEmbed("⚠️ Warnings", "<@"+target.ID+"> has **"+strconv.Itoa(count)+"** warnings.")
	SendEmbed(s, m.ChannelID, embed)
}

func (b *Bot) applyWarning(s *discordgo.Session, guildID, channelID, targetID, reason string) {
	err := b.store.AddWarning(context.Background(), targetID, guildID, reason)
	if err != nil {
		return
	}

	count, _ := b.store.GetWarningCount(context.Background(), targetID, guildID)

	embed := newEmbed("⚠️ Warning Issued", "<@"+targetID+"> has been warned.\n**Reason:** "+reason+"\n**Total Warnings:** "+strconv.Itoa(count)+"/3")
	SendEmbed(s, channelID, embed)

	if count >= 3 {
		timeoutUntil := time.Now().Add(10 * time.Minute)
		err = s.GuildMemberTimeout(guildID, targetID, &timeoutUntil)
		if err == nil {
			muteEmbed := newEmbed("🔇 Auto-Mod Action", "<@"+targetID+"> has reached 3 warnings and was automatically muted for 10 minutes.")
			SendEmbed(s, channelID, muteEmbed)
		}
	}
}
