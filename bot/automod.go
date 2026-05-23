package bot

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Kishan-Thanki/discord-ping/database"
	"github.com/bwmarrin/discordgo"
)

var badWordsRegex = regexp.MustCompile(`(?i)\b(fuck|shit|bitch|asshole|cunt|dick|faggot|nigger|kys|discord\.gift|free nitro)\b`)

func checkBadWords(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	if m.GuildID == "" || m.Author.Bot {
		return false
	}

	if badWordsRegex.MatchString(m.Content) {
		// Delete the message
		_ = s.ChannelMessageDelete(m.ChannelID, m.ID)

		// Warn the user
		applyWarning(s, m.GuildID, m.ChannelID, m.Author.ID, "Triggered auto-mod filter (bad words/spam)")
		return true
	}
	return false
}

func cmdWarn(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
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

	applyWarning(s, m.GuildID, m.ChannelID, target.ID, reason)
}

func cmdWarnings(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.GuildID == "" {
		return
	}

	target := m.Author
	if len(m.Mentions) > 0 {
		target = m.Mentions[0]
	}

	count, err := database.GetWarningCount(target.ID, m.GuildID)
	if err != nil {
		count = 0
	}

	embed := newEmbed("⚠️ Warnings", "<@"+target.ID+"> has **"+strconv.Itoa(count)+"** warnings.")
	SendEmbed(s, m.ChannelID, embed)
}

func applyWarning(s *discordgo.Session, guildID, channelID, targetID, reason string) {
	err := database.AddWarning(targetID, guildID, reason)
	if err != nil {
		return
	}

	count, _ := database.GetWarningCount(targetID, guildID)

	embed := newEmbed("⚠️ Warning Issued", "<@"+targetID+"> has been warned.\n**Reason:** "+reason+"\n**Total Warnings:** "+strconv.Itoa(count)+"/3")
	SendEmbed(s, channelID, embed)

	// Escalation: 3 warnings = 10 minute timeout
	if count >= 3 {
		timeoutUntil := time.Now().Add(10 * time.Minute)
		err = s.GuildMemberTimeout(guildID, targetID, &timeoutUntil)
		if err == nil {
			muteEmbed := newEmbed("🔇 Auto-Mod Action", "<@"+targetID+"> has reached 3 warnings and was automatically muted for 10 minutes.")
			SendEmbed(s, channelID, muteEmbed)
		}
	}
}
