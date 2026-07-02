package bot

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

const rateLimitDuration = 2 * time.Second

func (b *Bot) isRateLimited(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	if m.Author.ID == b.BotID {
		return false
	}

	val, exists := b.rateLimits.Load(m.Author.ID)
	if exists {
		lastCmd := val.(time.Time)
		if time.Since(lastCmd) < rateLimitDuration {
			embed := newEmbed("⏳ Slow Down!", "Please wait a moment before using another command.")
			SendEmbed(s, m.ChannelID, embed)
			return true
		}
	}

	b.rateLimits.Store(m.Author.ID, time.Now())
	return false
}

func (b *Bot) startRateLimitCleaner() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-b.stopCleanup:
			return
		case <-ticker.C:
			now := time.Now()
			b.rateLimits.Range(func(key, val interface{}) bool {
				lastCmd, ok := val.(time.Time)
				if !ok || now.Sub(lastCmd) > rateLimitDuration {
					b.rateLimits.Delete(key)
				}
				return true
			})
		}
	}
}
