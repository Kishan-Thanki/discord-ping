package bot

import (
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

var rateLimits sync.Map

const rateLimitDuration = 2 * time.Second

func isRateLimited(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	if m.Author.ID == BotID {
		return false
	}

	val, exists := rateLimits.Load(m.Author.ID)
	if exists {
		lastCmd := val.(time.Time)
		if time.Since(lastCmd) < rateLimitDuration {
			embed := newEmbed("⏳ Slow Down!", "Please wait a moment before using another command.")
			SendEmbed(s, m.ChannelID, embed)
			return true
		}
	}

	rateLimits.Store(m.Author.ID, time.Now())
	return false
}
