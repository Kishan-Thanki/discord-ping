package bot

import (
	"strconv"
	"time"

	"github.com/Kishan-Thanki/discord-ping/database"
	"github.com/bwmarrin/discordgo"
)

func cmdDaily(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.GuildID == "" {
		return
	}

	user, err := database.GetUser(m.Author.ID, m.GuildID)
	if err != nil {
		return
	}

	now := time.Now()

	if user.LastDaily != "" {
		last, err := time.Parse(time.RFC3339, user.LastDaily)
		if err == nil {
			if now.Sub(last) < 24*time.Hour {
				remaining := (24 * time.Hour) - now.Sub(last)
				hours := int(remaining.Hours())
				minutes := int(remaining.Minutes()) % 60

				embed := newEmbed("⏳ Daily Claimed", "You already claimed your daily coins! Come back in **"+strconv.Itoa(hours)+"h "+strconv.Itoa(minutes)+"m**.")
				SendEmbed(s, m.ChannelID, embed)
				return
			}
		}
	}

	reward := 100
	err = database.UpdateUserDaily(m.Author.ID, m.GuildID, user.Balance+reward, now.Format(time.RFC3339))
	if err != nil {
		return
	}

	embed := newEmbed("💰 Daily Reward", "You claimed your daily **"+strconv.Itoa(reward)+" coins**! Your new balance is **"+strconv.Itoa(user.Balance+reward)+"**.")
	SendEmbed(s, m.ChannelID, embed)
}

func cmdBalance(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.GuildID == "" {
		return
	}

	target := m.Author
	if len(m.Mentions) > 0 {
		target = m.Mentions[0]
	}

	user, err := database.GetUser(target.ID, m.GuildID)
	if err != nil {
		return
	}

	embed := newEmbed("🏦 Balance", "<@"+target.ID+"> has **"+strconv.Itoa(user.Balance)+" coins**.")
	SendEmbed(s, m.ChannelID, embed)
}

func cmdGive(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if m.GuildID == "" {
		return
	}

	if len(m.Mentions) == 0 || len(args) < 3 {
		embed := newEmbed("❌ Usage", "Usage: `!give @user <amount>`")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	target := m.Mentions[0]
	if target.ID == m.Author.ID {
		embed := newEmbed("❌ Error", "You can't give coins to yourself.")
		SendEmbed(s, m.ChannelID, embed)
		return
	}
	if target.Bot {
		embed := newEmbed("❌ Error", "You can't give coins to a bot.")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	amount, err := strconv.Atoi(args[len(args)-1])
	if err != nil || amount <= 0 {
		embed := newEmbed("❌ Error", "Please provide a valid positive amount.")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	sender, _ := database.GetUser(m.Author.ID, m.GuildID)
	if sender.Balance < amount {
		embed := newEmbed("💸 Insufficient Funds", "You don't have enough coins. Your balance is **"+strconv.Itoa(sender.Balance)+"**.")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	receiver, _ := database.GetUser(target.ID, m.GuildID)

	database.UpdateUserBalance(m.Author.ID, m.GuildID, sender.Balance-amount)
	database.UpdateUserBalance(target.ID, m.GuildID, receiver.Balance+amount)

	embed := newEmbed("🤝 Transfer Complete", "You gave **"+strconv.Itoa(amount)+" coins** to <@"+target.ID+">.")
	SendEmbed(s, m.ChannelID, embed)
}
