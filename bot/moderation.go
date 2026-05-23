package bot

import (
	"log/slog"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func cmdKick(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if m.GuildID == "" {
		embed := newEmbed("❌ Error", "This command can only be used in a server.")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	if !hasPermission(s, m.Author.ID, m.ChannelID, discordgo.PermissionKickMembers) {
		embed := newEmbed("🚫 Permission Denied", "You need the **Kick Members** permission to use this command.")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	if len(m.Mentions) == 0 {
		embed := newEmbed("❌ Usage", "Usage: `!kick @user [reason]`")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	target := m.Mentions[0]

	if target.ID == BotID {
		embed := newEmbed("🤖 Nope", "I can't kick myself!")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	reason := "No reason provided"
	if len(args) > 2 {
		reason = strings.Join(args[2:], " ")
	}

	err := s.GuildMemberDeleteWithReason(m.GuildID, target.ID, reason)
	if err != nil {
		slog.Error("Failed to kick member", "target", target.Username, "error", err)
		embed := newEmbed("❌ Error", "Failed to kick "+target.Mention()+". Make sure my role is higher than theirs.")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	slog.Info("Member kicked", "target", target.Username, "by", m.Author.Username, "reason", reason)
	embed := newEmbed("👢 Member Kicked", "**"+target.Username+"** has been kicked by "+m.Author.Mention()+".")
	embed.Fields = []*discordgo.MessageEmbedField{
		{Name: "Reason", Value: reason, Inline: false},
	}
	SendEmbed(s, m.ChannelID, embed)
}

func cmdBan(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if m.GuildID == "" {
		embed := newEmbed("❌ Error", "This command can only be used in a server.")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	if !hasPermission(s, m.Author.ID, m.ChannelID, discordgo.PermissionBanMembers) {
		embed := newEmbed("🚫 Permission Denied", "You need the **Ban Members** permission to use this command.")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	if len(m.Mentions) == 0 {
		embed := newEmbed("❌ Usage", "Usage: `!ban @user [reason]`")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	target := m.Mentions[0]

	if target.ID == BotID {
		embed := newEmbed("🤖 Nope", "I can't ban myself!")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	reason := "No reason provided"
	if len(args) > 2 {
		reason = strings.Join(args[2:], " ")
	}

	err := s.GuildBanCreateWithReason(m.GuildID, target.ID, reason, 0)
	if err != nil {
		slog.Error("Failed to ban member", "target", target.Username, "error", err)
		embed := newEmbed("❌ Error", "Failed to ban "+target.Mention()+". Make sure my role is higher than theirs.")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	slog.Info("Member banned", "target", target.Username, "by", m.Author.Username, "reason", reason)
	embed := newEmbed("🔨 Member Banned", "**"+target.Username+"** has been banned by "+m.Author.Mention()+".")
	embed.Fields = []*discordgo.MessageEmbedField{
		{Name: "Reason", Value: reason, Inline: false},
	}
	SendEmbed(s, m.ChannelID, embed)
}

func cmdMute(s *discordgo.Session, m *discordgo.MessageCreate, _ []string) {
	if m.GuildID == "" {
		embed := newEmbed("❌ Error", "This command can only be used in a server.")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	if !hasPermission(s, m.Author.ID, m.ChannelID, discordgo.PermissionModerateMembers) {
		embed := newEmbed("🚫 Permission Denied", "You need the **Moderate Members** permission to use this command.")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	if len(m.Mentions) == 0 {
		embed := newEmbed("❌ Usage", "Usage: `!mute @user`")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	target := m.Mentions[0]

	if target.ID == BotID {
		embed := newEmbed("🤖 Nope", "I can't mute myself!")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	muteDuration := 10 * time.Minute
	until := time.Now().Add(muteDuration)

	err := s.GuildMemberTimeout(m.GuildID, target.ID, &until)
	if err != nil {
		slog.Error("Failed to mute member", "target", target.Username, "error", err)
		embed := newEmbed("❌ Error", "Failed to mute "+target.Mention()+". Make sure my role is higher than theirs.")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	slog.Info("Member muted", "target", target.Username, "by", m.Author.Username, "duration", muteDuration.String())
	embed := newEmbed("🔇 Member Muted", "**"+target.Username+"** has been muted for **10 minutes** by "+m.Author.Mention()+".")
	SendEmbed(s, m.ChannelID, embed)
}

func hasPermission(s *discordgo.Session, userID, channelID string, perm int64) bool {
	perms, err := s.State.UserChannelPermissions(userID, channelID)
	if err != nil {
		slog.Warn("Failed to check permissions via state, falling back", "error", err)

		// Fallback: check if user is guild owner
		channel, err := s.Channel(channelID)
		if err != nil {
			return false
		}
		guild, err := s.Guild(channel.GuildID)
		if err != nil {
			return false
		}
		if guild.OwnerID == userID {
			return true
		}

		// Check roles manually
		member, err := s.GuildMember(channel.GuildID, userID)
		if err != nil {
			return false
		}
		for _, roleID := range member.Roles {
			for _, role := range guild.Roles {
				if role.ID == roleID && role.Permissions&perm != 0 {
					return true
				}
			}
		}
		return false
	}

	return perms&perm != 0
}
