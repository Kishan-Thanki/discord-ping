package bot

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

// SendEmbed safely sends an embed message and logs any resulting network error.
func SendEmbed(s *discordgo.Session, channelID string, embed *discordgo.MessageEmbed) {
	_, err := s.ChannelMessageSendEmbed(channelID, embed)
	if err != nil {
		slog.Error("Failed to send embed message", "channel_id", channelID, "error", err)
	}
}

// EditEmbed safely edits an existing message with a new embed and logs any error.
func EditEmbed(s *discordgo.Session, channelID, messageID string, embed *discordgo.MessageEmbed) {
	_, err := s.ChannelMessageEditEmbed(channelID, messageID, embed)
	if err != nil {
		slog.Error("Failed to edit embed message", "channel_id", channelID, "message_id", messageID, "error", err)
	}
}

// SendComplex safely sends a complex message (buttons, files, etc) and logs any error.
func SendComplex(s *discordgo.Session, channelID string, data *discordgo.MessageSend) {
	_, err := s.ChannelMessageSendComplex(channelID, data)
	if err != nil {
		slog.Error("Failed to send complex message", "channel_id", channelID, "error", err)
	}
}

// EditComplex safely edits a complex message and logs any error.
func EditComplex(s *discordgo.Session, data *discordgo.MessageEdit) {
	_, err := s.ChannelMessageEditComplex(data)
	if err != nil {
		slog.Error("Failed to edit complex message", "channel_id", data.Channel, "message_id", data.ID, "error", err)
	}
}

// hasPermission checks if a user has a specific permission in a channel.
func hasPermission(s *discordgo.Session, userID, channelID string, permission int64) bool {
	perms, err := s.State.UserChannelPermissions(userID, channelID)
	if err != nil {
		perms, err = s.UserChannelPermissions(userID, channelID)
		if err != nil {
			return false
		}
	}
	return perms&permission == permission
}
