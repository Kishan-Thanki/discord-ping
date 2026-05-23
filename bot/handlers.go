package bot

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/Kishan-Thanki/discord-ping/database"
	"github.com/bwmarrin/discordgo"
)

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotID || m.Author.Bot {
		return
	}

	if checkBadWords(s, m) {
		return
	}

	prefix := "!"
	if m.GuildID != "" {
		p, err := database.GetPrefix(m.GuildID)
		if err == nil {
			prefix = p
		}
	}

	content := strings.ToLower(m.Content)

	handlePassiveXP(s, m)

	if strings.HasPrefix(m.Content, prefix) {
		if isRateLimited(s, m) {
			return
		}

		cmd := strings.TrimPrefix(m.Content, prefix)
		args := strings.Fields(cmd)
		if len(args) == 0 {
			return
		}

		slog.Info("Command received",
			"command", args[0],
			"user", m.Author.Username,
			"channel", m.ChannelID,
			"guild", m.GuildID,
		)

		switch strings.ToLower(args[0]) {
		case "ping":
			cmdPing(s, m)
		case "version", "about":
			cmdVersion(s, m)
		case "uptime":
			cmdUptime(s, m)
		case "serverinfo":
			cmdServerInfo(s, m)
		case "avatar":
			cmdAvatar(s, m)
		case "roll":
			cmdRoll(s, m, args)
		case "kick":
			cmdKick(s, m, args)
		case "ban":
			cmdBan(s, m, args)
		case "mute":
			cmdMute(s, m, args)
		case "8ball":
			cmd8Ball(s, m, args)
		case "coinflip":
			cmdCoinflip(s, m)
		case "poll":
			cmdPoll(s, m, args)
		case "trivia":
			cmdTrivia(s, m)
		case "rank":
			cmdRank(s, m)
		case "leaderboard":
			cmdLeaderboard(s, m)
		case "daily":
			cmdDaily(s, m)
		case "balance":
			cmdBalance(s, m)
		case "give":
			cmdGive(s, m, args)
		case "remind":
			cmdRemind(s, m, args)
		case "help":
			cmdHelp(s, m)
		case "setprefix":
			cmdSetPrefix(s, m, args)
		case "slots":
			cmdSlots(s, m, args)
		case "blackjack":
			cmdBlackjack(s, m, args)
		case "warn":
			cmdWarn(s, m, args)
		case "warnings":
			cmdWarnings(s, m)
		case "wordle":
			cmdWordle(s, m)
		case "guess":
			cmdGuess(s, m, args)
		}
		return
	}

	if strings.Contains(content, "hello bot") {
		if !isRateLimited(s, m) {
			embed := newEmbed("👋 Hello!", "Hey "+m.Author.Mention()+"! How can I help you today?")
			SendEmbed(s, m.ChannelID, embed)
		}
	}

	if strings.Contains(content, "react to this") {
		_ = s.MessageReactionAdd(m.ChannelID, m.ID, "👀")
		_ = s.MessageReactionAdd(m.ChannelID, m.ID, "👍")
	}
}

func slashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionMessageComponent {
		handleBlackjackInteraction(s, i)
		return
	}

	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	slog.Info("Slash command received", "command", i.ApplicationCommandData().Name)

	switch i.ApplicationCommandData().Name {
	case "ping":
		latency := s.HeartbeatLatency().Milliseconds()
		embed := newEmbed("🏓 Pong!", "Latency: **"+strconv.FormatInt(latency, 10)+"ms**")

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		})
		if err != nil {
			slog.Error("Failed to respond to slash command", "error", err)
		}
	}
}

func welcomeHandler(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	guild, err := s.Guild(m.GuildID)
	if err != nil {
		slog.Error("Failed to get guild for welcome message", "error", err)
		return
	}

	if guild.SystemChannelID == "" {
		slog.Warn("No system channel configured for guild", "guild", guild.Name)
		return
	}

	avatarURL := m.User.AvatarURL("")

	buf, err := generateWelcomeImage(avatarURL, m.User.Username)
	if err == nil {
		SendComplex(s, guild.SystemChannelID, &discordgo.MessageSend{
			Content: "Welcome to the server, <@" + m.User.ID + ">!",
			Files: []*discordgo.File{
				{
					Name:   "welcome.png",
					Reader: buf,
				},
			},
		})
	} else {
		slog.Error("Failed to generate welcome image", "err", err)
		embed := newEmbed("👋 Welcome!", "Welcome to the server, <@"+m.User.ID+">!")
		SendEmbed(s, guild.SystemChannelID, embed)
	}
}
