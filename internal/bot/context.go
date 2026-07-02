package bot

import (
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

type PingContext interface {
	Respond(embed *discordgo.MessageEmbed) error
	SentAt() time.Time
}

type MessageContext struct {
	Session *discordgo.Session
	Message *discordgo.MessageCreate
	sentMsg *discordgo.Message
}

func (c *MessageContext) Respond(embed *discordgo.MessageEmbed) error {
	if c.sentMsg == nil {
		var err error
		c.sentMsg, err = c.Session.ChannelMessageSendEmbed(c.Message.ChannelID, embed)
		return err
	}
	_, err := c.Session.ChannelMessageEditEmbed(c.Message.ChannelID, c.sentMsg.ID, embed)
	return err
}

func (c *MessageContext) SentAt() time.Time {
	return c.Message.Timestamp
}

type InteractionContext struct {
	Session     *discordgo.Session
	Interaction *discordgo.InteractionCreate
	responded   bool
}

func (c *InteractionContext) Respond(embed *discordgo.MessageEmbed) error {
	if !c.responded {
		err := c.Session.InteractionRespond(c.Interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		})
		if err == nil {
			c.responded = true
		}
		return err
	}
	_, err := c.Session.InteractionResponseEdit(c.Interaction.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	})
	return err
}

func (c *InteractionContext) SentAt() time.Time {
	if id, err := strconv.ParseInt(c.Interaction.Interaction.ID, 10, 64); err == nil {
		return time.UnixMilli((id >> 22) + 1420070400000)
	}
	return time.Now()
}
