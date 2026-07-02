package bot

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

func SendEmbed(s *discordgo.Session, channelID string, embed *discordgo.MessageEmbed) {
	_, err := s.ChannelMessageSendEmbed(channelID, embed)
	if err != nil {
		slog.Error("Failed to send embed message", "channel_id", channelID, "error", err)
	}
}

func EditEmbed(s *discordgo.Session, channelID, messageID string, embed *discordgo.MessageEmbed) {
	_, err := s.ChannelMessageEditEmbed(channelID, messageID, embed)
	if err != nil {
		slog.Error("Failed to edit embed message", "channel_id", channelID, "message_id", messageID, "error", err)
	}
}

func hasPermission(s *discordgo.Session, userID, channelID string, permission int64) bool {
	var perms int64
	var err error
	if s.StateEnabled && s.State != nil {
		perms, err = s.State.UserChannelPermissions(userID, channelID)
	} else {
		err = discordgo.ErrStateNotFound
	}
	if err != nil {
		perms, err = s.UserChannelPermissions(userID, channelID)
		if err != nil {
			return false
		}
	}
	return perms&permission == permission
}
