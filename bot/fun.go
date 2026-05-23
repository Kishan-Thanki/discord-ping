package bot

import (
	"encoding/json"
	"html"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Kishan-Thanki/discord-ping/database"
	"github.com/bwmarrin/discordgo"
)

var eightBallResponses = []string{
	"It is certain.", "It is decidedly so.", "Without a doubt.",
	"Yes definitely.", "You may rely on it.", "As I see it, yes.",
	"Most likely.", "Outlook good.", "Yes.", "Signs point to yes.",
	"Reply hazy, try again.", "Ask again later.", "Better not tell you now.",
	"Cannot predict now.", "Concentrate and ask again.", "Don't count on it.",
	"My reply is no.", "My sources say no.", "Outlook not so good.",
	"Very doubtful.",
}

func cmd8Ball(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 2 {
		embed := newEmbed("❌ Error", "You need to ask a question! `!8ball [question]`")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	question := strings.Join(args[1:], " ")
	response := eightBallResponses[rand.Intn(len(eightBallResponses))]

	embed := newEmbed("🎱 Magic 8-Ball", "")
	embed.Fields = []*discordgo.MessageEmbedField{
		{Name: "Question", Value: question, Inline: false},
		{Name: "Answer", Value: response, Inline: false},
	}

	SendEmbed(s, m.ChannelID, embed)
}

func cmdCoinflip(s *discordgo.Session, m *discordgo.MessageCreate) {
	result := "Tails"
	if rand.Intn(2) == 0 {
		result = "Heads"
	}

	embed := newEmbed("🪙 Coinflip", "The coin landed on **"+result+"**!")
	SendEmbed(s, m.ChannelID, embed)
}

var pollEmojis = []string{"1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣", "6️⃣", "7️⃣", "8️⃣", "9️⃣", "🔟"}

func cmdPoll(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	// Simple argument parser expecting quotes
	raw := strings.Join(args[1:], " ")
	parts := strings.Split(raw, "\"")

	var parsed []string
	for i, part := range parts {
		if i%2 != 0 && strings.TrimSpace(part) != "" {
			parsed = append(parsed, part)
		}
	}

	if len(parsed) < 3 || len(parsed) > 11 {
		embed := newEmbed("❌ Usage", "Usage: `!poll \"Question\" \"Option 1\" \"Option 2\" ...` (Max 10 options)")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	question := parsed[0]
	options := parsed[1:]

	desc := ""
	for i, opt := range options {
		desc += pollEmojis[i] + " " + opt + "\n"
	}

	embed := newEmbed("📊 Poll: "+question, desc)
	embed.Author = &discordgo.MessageEmbedAuthor{
		Name:    m.Author.Username,
		IconURL: m.Author.AvatarURL("256"),
	}

	msg, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err == nil {
		for i := range options {
			_ = s.MessageReactionAdd(m.ChannelID, msg.ID, pollEmojis[i])
		}
	}
}

// Trivia structs
type triviaResponse struct {
	ResponseCode int              `json:"response_code"`
	Results      []triviaQuestion `json:"results"`
}

type triviaQuestion struct {
	Category         string   `json:"category"`
	Type             string   `json:"type"`
	Difficulty       string   `json:"difficulty"`
	Question         string   `json:"question"`
	CorrectAnswer    string   `json:"correct_answer"`
	IncorrectAnswers []string `json:"incorrect_answers"`
}

var activeTrivia = sync.Map{}

func cmdTrivia(s *discordgo.Session, m *discordgo.MessageCreate) {
	if _, active := activeTrivia.Load(m.ChannelID); active {
		embed := newEmbed("❌ Error", "A trivia question is already active in this channel!")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	activeTrivia.Store(m.ChannelID, true)
	defer activeTrivia.Delete(m.ChannelID)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("https://opentdb.com/api.php?amount=1&type=multiple")
	if err != nil {
		embed := newEmbed("❌ Error", "Failed to fetch trivia question.")
		SendEmbed(s, m.ChannelID, embed)
		return
	}
	defer resp.Body.Close()

	var tr triviaResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil || len(tr.Results) == 0 {
		embed := newEmbed("❌ Error", "Failed to parse trivia question.")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	q := tr.Results[0]
	question := html.UnescapeString(q.Question)
	correct := html.UnescapeString(q.CorrectAnswer)

	options := []string{correct}
	for _, inc := range q.IncorrectAnswers {
		options = append(options, html.UnescapeString(inc))
	}

	// Shuffle options
	rand.Shuffle(len(options), func(i, j int) {
		options[i], options[j] = options[j], options[i]
	})

	correctIdx := 0
	desc := ""
	for i, opt := range options {
		if opt == correct {
			correctIdx = i
		}
		desc += "**" + strconv.Itoa(i+1) + ".** " + opt + "\n"
	}

	embed := newEmbed("🧠 Trivia Time!", question)
	embed.Description = desc
	embed.Footer = &discordgo.MessageEmbedFooter{Text: "Category: " + html.UnescapeString(q.Category) + " | Difficulty: " + strings.Title(q.Difficulty)}

	SendEmbed(s, m.ChannelID, embed)

	// Create message collector
	stopChan := make(chan struct{})
	winnerChan := make(chan string)

	removeHandler := s.AddHandler(func(ss *discordgo.Session, mc *discordgo.MessageCreate) {
		if mc.ChannelID != m.ChannelID || mc.Author.Bot {
			return
		}
		if mc.Content == strconv.Itoa(correctIdx+1) {
			select {
			case winnerChan <- mc.Author.ID:
			default:
			}
		}
	})
	defer removeHandler()

	select {
	case winnerID := <-winnerChan:
		close(stopChan)

		// Award XP and Coins
		dbUser, _ := database.GetUser(winnerID, m.GuildID)
		if dbUser != nil {
			database.UpdateUserXP(winnerID, m.GuildID, dbUser.XP+50, dbUser.Level)
			database.UpdateUserBalance(winnerID, m.GuildID, dbUser.Balance+25)
		}

		embed := newEmbed("🎉 Trivia Winner!", "<@"+winnerID+"> got it right!\n\nThe answer was **"+correct+"**.\n*Awarded 50 XP and 25 coins.*")
		SendEmbed(s, m.ChannelID, embed)

	case <-time.After(15 * time.Second):
		close(stopChan)
		embed := newEmbed("⌛ Trivia Over", "Time's up! No one got it right.\n\nThe correct answer was **"+correct+"**.")
		SendEmbed(s, m.ChannelID, embed)
	}
}
