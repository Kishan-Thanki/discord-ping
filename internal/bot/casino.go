package bot

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

var slotEmojis = []string{"🍒", "🔔", "💎", "🍋", "🍉", "7️⃣"}

func (b *Bot) cmdSlots(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if m.GuildID == "" {
		return
	}

	if len(args) < 2 {
		embed := newEmbed("❌ Usage", "Usage: `!slots <bet>`")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	bet, err := strconv.Atoi(args[1])
	if err != nil || bet <= 0 {
		embed := newEmbed("❌ Error", "Please provide a valid positive bet amount.")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	user, err := b.store.GetUser(context.Background(), m.Author.ID, m.GuildID)
	if err != nil {
		return
	}

	if user.Balance < bet {
		embed := newEmbed("💸 Insufficient Funds", "You don't have enough coins. Your balance is **"+strconv.Itoa(user.Balance)+"**.")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	// Deduct bet immediately
	err = b.store.UpdateUserBalance(context.Background(), m.Author.ID, m.GuildID, user.Balance-bet)
	if err != nil {
		return
	}

	// Initial spinning message
	embed := newEmbed("🎰 Slots", "Spinning...\n\n| 🔄 | 🔄 | 🔄 |")
	msg, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		// Refund if failed to send
		_ = b.store.UpdateUserBalance(context.Background(), m.Author.ID, m.GuildID, user.Balance)
		return
	}

	// Animation frames
	for i := 0; i < 2; i++ {
		time.Sleep(800 * time.Millisecond)
		frame := "| " + getRandomSlot() + " | " + getRandomSlot() + " | " + getRandomSlot() + " |"
		embed.Description = "Spinning...\n\n" + frame
		EditEmbed(s, m.ChannelID, msg.ID, embed)
	}

	time.Sleep(800 * time.Millisecond)

	// Final result
	res1 := getRandomSlot()
	res2 := getRandomSlot()
	res3 := getRandomSlot()

	finalStr := "| " + res1 + " | " + res2 + " | " + res3 + " |"

	payout := 0
	msgText := ""

	if res1 == res2 && res2 == res3 {
		payout = bet * 10
		msgText = "🎉 **JACKPOT!** You won **" + strconv.Itoa(payout) + " coins**! (10x payout)"
	} else if res1 == res2 || res2 == res3 || res1 == res3 {
		payout = bet * 2
		msgText = "✨ You won **" + strconv.Itoa(payout) + " coins**! (2x payout)"
	} else {
		msgText = "😢 You lost. Better luck next time!"
	}

	if payout > 0 {
		// Re-fetch user to avoid race conditions with other commands
		u, _ := b.store.GetUser(context.Background(), m.Author.ID, m.GuildID)
		_ = b.store.UpdateUserBalance(context.Background(), m.Author.ID, m.GuildID, u.Balance+payout)
	}

	embed.Description = finalStr + "\n\n" + msgText
	EditEmbed(s, m.ChannelID, msg.ID, embed)
}

func getRandomSlot() string {
	return slotEmojis[rand.Intn(len(slotEmojis))]
}
