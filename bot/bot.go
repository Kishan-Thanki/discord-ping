package bot

import (
	"fmt"

	"github.com/Kishan-Thanki/discord-ping/config"
	"github.com/bwmarrin/discordgo"
)

var BotID string

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotID {
		return
	}

	if m.Content == "ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "pong")
	}
}

func Start() {
	if config.Token == "" {
		fmt.Println("Error: DISCORD_BOT_TOKEN is missing. Set it in your environment variables.")
		return
	}

	goBot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("Error creating bot session:", err)
		return
	}

	u, err := goBot.User("@me")
	if err != nil {
		fmt.Println("Error getting bot user:", err)
		return
	}

	BotID = u.ID
	goBot.AddHandler(messageHandler)

	err = goBot.Open()
	if err != nil {
		fmt.Println("Error opening connection:", err)
		return
	}

	fmt.Println("Bot is running!")
}
