package bot

import (
	"context"
	"log/slog"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == b.BotID || m.Author.Bot {
		return
	}

	if b.checkBadWords(s, m) {
		return
	}

	prefix := b.cfg.BotPrefix
	if m.GuildID != "" {
		p, err := b.store.GetPrefix(context.Background(), m.GuildID)
		if err == nil && p != "" {
			prefix = p
		}
	}

	content := strings.ToLower(m.Content)

	b.handlePassiveXP(s, m)

	if strings.HasPrefix(m.Content, prefix) {
		if b.isRateLimited(s, m) {
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
			b.cmdPing(s, m)
		case "version", "about":
			b.cmdVersion(s, m)
		case "uptime":
			b.cmdUptime(s, m)
		case "serverinfo":
			b.cmdServerInfo(s, m)
		case "avatar":
			b.cmdAvatar(s, m)
		case "roll":
			b.cmdRoll(s, m, args)
		case "kick":
			b.cmdKick(s, m, args)
		case "ban":
			b.cmdBan(s, m, args)
		case "mute":
			b.cmdMute(s, m, args)
		case "8ball":
			b.cmd8Ball(s, m, args)
		case "coinflip":
			b.cmdCoinflip(s, m)
		case "poll":
			b.cmdPoll(s, m, args)
		case "trivia":
			b.cmdTrivia(s, m)
		case "rank":
			b.cmdRank(s, m)
		case "leaderboard":
			b.cmdLeaderboard(s, m)
		case "daily":
			b.cmdDaily(s, m)
		case "balance":
			b.cmdBalance(s, m)
		case "give":
			b.cmdGive(s, m, args)
		case "remind":
			b.cmdRemind(s, m, args)
		case "help":
			b.cmdHelp(s, m)
		case "setprefix":
			b.cmdSetPrefix(s, m, args)
		case "slots":
			b.cmdSlots(s, m, args)
		case "blackjack":
			b.cmdBlackjack(s, m, args)
		case "warn":
			b.cmdWarn(s, m, args)
		case "warnings":
			b.cmdWarnings(s, m)
		case "wordle":
			b.cmdWordle(s, m)
		case "guess":
			b.cmdGuess(s, m, args)
		}
		return
	}

	if strings.Contains(content, "hello bot") {
		if !b.isRateLimited(s, m) {
			embed := newEmbed("👋 Hello!", "Hey "+m.Author.Mention()+"! How can I help you today?")
			SendEmbed(s, m.ChannelID, embed)
		}
	}

	if strings.Contains(content, "react to this") {
		_ = s.MessageReactionAdd(m.ChannelID, m.ID, "👀")
		_ = s.MessageReactionAdd(m.ChannelID, m.ID, "👍")
	}
}

func (b *Bot) slashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionMessageComponent {
		b.handleBlackjackInteraction(s, i)
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

func (b *Bot) welcomeHandler(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
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

	buf, err := b.generateWelcomeImage(avatarURL, m.User.Username)
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
