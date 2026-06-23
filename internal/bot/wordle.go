package bot

import (
	"context"
	"math/rand"
	"strconv"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var wordleWords = []string{
	"APPLE", "TRAIN", "GHOST", "BRAIN", "PLANT", "WATER", "EARTH", "MAGIC", "SWORD", "SHIELD",
	"NIGHT", "DREAM", "FLAME", "STORM", "HEART", "BLOOD", "STONE", "RIVER", "OCEAN", "SPACE",
	"ALIEN", "ROBOT", "CRAZY", "FUNNY", "HAPPY", "SMILE", "LAUGH", "DANCE", "MUSIC", "SOUND",
	"VOICE", "LIGHT", "COLOR", "PAINT", "BRUSH", "PAPER", "BOOKS", "WRITE", "LEARN", "STUDY",
	"THINK", "SMART", "BRAVE", "PROUD", "FIGHT", "PEACE", "WORLD", "HOUSE", "CHAIR", "TABLE",
}

type WordleGame struct {
	UserID     string
	TargetWord string
	Guesses    []string
	Finished   bool
	mu         sync.Mutex
}

func (b *Bot) cmdWordle(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.GuildID == "" {
		return
	}

	if _, exists := b.wordleGames.Load(m.Author.ID); exists {
		embed := newEmbed("❌ Error", "You already have an active Wordle game! Make a guess using `!guess [word]`.")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	target := wordleWords[rand.Intn(len(wordleWords))]
	game := &WordleGame{
		UserID:     m.Author.ID,
		TargetWord: target,
		Guesses:    []string{},
		Finished:   false,
	}
	b.wordleGames.Store(m.Author.ID, game)

	embed := newEmbed("🟩 Discord Wordle", "A new game has started!\nUse `!guess <5-letter-word>` to play.\n\n⬜ ⬜ ⬜ ⬜ ⬜\n⬜ ⬜ ⬜ ⬜ ⬜\n⬜ ⬜ ⬜ ⬜ ⬜\n⬜ ⬜ ⬜ ⬜ ⬜\n⬜ ⬜ ⬜ ⬜ ⬜\n⬜ ⬜ ⬜ ⬜ ⬜")
	SendEmbed(s, m.ChannelID, embed)
}

func (b *Bot) cmdGuess(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if m.GuildID == "" {
		return
	}

	val, exists := b.wordleGames.Load(m.Author.ID)
	if !exists {
		embed := newEmbed("❌ Error", "You don't have an active Wordle game! Start one with `!wordle`.")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	if len(args) < 2 {
		embed := newEmbed("❌ Usage", "Usage: `!guess <5-letter-word>`")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	guess := strings.ToUpper(args[1])
	if len(guess) != 5 {
		embed := newEmbed("❌ Error", "Your guess must be exactly 5 letters long!")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	game := val.(*WordleGame)

	game.mu.Lock()
	defer game.mu.Unlock()

	game.Guesses = append(game.Guesses, guess)

	var boardDesc strings.Builder
	won := false

	for _, g := range game.Guesses {
		boardDesc.WriteString(renderWordleGuess(g, game.TargetWord))
		boardDesc.WriteString("  ")
		boardDesc.WriteString(g)
		boardDesc.WriteString("\n")
	}

	for i := len(game.Guesses); i < 6; i++ {
		boardDesc.WriteString("⬜ ⬜ ⬜ ⬜ ⬜\n")
	}

	if guess == game.TargetWord {
		won = true
		game.Finished = true
	} else if len(game.Guesses) >= 6 {
		game.Finished = true
	}

	embed := newEmbed("🟩 Discord Wordle", boardDesc.String())

	if game.Finished {
		b.wordleGames.Delete(m.Author.ID)

		if won {
			payouts := []int{1000, 500, 250, 100, 50, 25}
			reward := payouts[len(game.Guesses)-1]

			user, _ := b.store.GetUser(context.Background(), m.Author.ID, m.GuildID)
			_ = b.store.UpdateUserBalance(context.Background(), m.Author.ID, m.GuildID, user.Balance+reward)

			embed.Description += "\n\n🎉 **You won!** You guessed the word in " + strconv.Itoa(len(game.Guesses)) + " tries and earned **" + strconv.Itoa(reward) + " coins**!"
		} else {
			embed.Description += "\n\n😢 **Game Over!** The word was **" + game.TargetWord + "**."
		}
	} else {
		embed.Description += "\n\n*You have " + strconv.Itoa(6-len(game.Guesses)) + " guesses left.*"
	}

	SendEmbed(s, m.ChannelID, embed)
}

func renderWordleGuess(guess, target string) string {
	result := []string{"⬛", "⬛", "⬛", "⬛", "⬛"}
	targetRunes := []rune(target)
	guessRunes := []rune(guess)

	used := make([]bool, 5)

	for i := 0; i < 5; i++ {
		if guessRunes[i] == targetRunes[i] {
			result[i] = "🟩"
			used[i] = true
		}
	}

	for i := 0; i < 5; i++ {
		if result[i] == "🟩" {
			continue
		}

		for j := 0; j < 5; j++ {
			if !used[j] && guessRunes[i] == targetRunes[j] {
				result[i] = "🟨"
				used[j] = true
				break
			}
		}
	}

	return strings.Join(result, " ")
}
