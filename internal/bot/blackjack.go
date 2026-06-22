package bot

import (
	"context"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

var suits = []string{"♠️", "♥️", "♦️", "♣️"}
var ranks = []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

type Card struct {
	Suit  string
	Rank  string
	Value int
}

func (c Card) String() string {
	return c.Rank + c.Suit
}

type Deck struct {
	Cards []Card
}

func newDeck() *Deck {
	d := &Deck{}
	for _, suit := range suits {
		for _, rank := range ranks {
			val := 0
			switch rank {
			case "A":
				val = 11
			case "J", "Q", "K":
				val = 10
			default:
				val, _ = strconv.Atoi(rank)
			}
			d.Cards = append(d.Cards, Card{Suit: suit, Rank: rank, Value: val})
		}
	}
	rand.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})
	return d
}

func (d *Deck) Draw() Card {
	c := d.Cards[0]
	d.Cards = d.Cards[1:]
	return c
}

func calculateHand(hand []Card) int {
	total := 0
	aces := 0
	for _, c := range hand {
		total += c.Value
		if c.Rank == "A" {
			aces++
		}
	}
	for total > 21 && aces > 0 {
		total -= 10
		aces--
	}
	return total
}

type BlackjackGame struct {
	UserID     string
	GuildID    string
	Bet        int
	Deck       *Deck
	PlayerHand []Card
	DealerHand []Card
	MessageID  string
	ChannelID  string
	LastActive time.Time
	mu         sync.Mutex
}

func (b *Bot) cmdBlackjack(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if m.GuildID == "" {
		return
	}

	if _, active := b.blackjackGames.Load(m.Author.ID); active {
		embed := newEmbed("❌ Error", "You already have an active Blackjack game!")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	if len(args) < 2 {
		embed := newEmbed("❌ Usage", "Usage: `!blackjack <bet>`")
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
	if err != nil || user.Balance < bet {
		embed := newEmbed("💸 Insufficient Funds", "You don't have enough coins. Your balance is **"+strconv.Itoa(user.Balance)+"**.")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	_ = b.store.UpdateUserBalance(context.Background(), m.Author.ID, m.GuildID, user.Balance-bet)

	deck := newDeck()
	game := &BlackjackGame{
		UserID:     m.Author.ID,
		GuildID:    m.GuildID,
		Bet:        bet,
		Deck:       deck,
		PlayerHand: []Card{deck.Draw(), deck.Draw()},
		DealerHand: []Card{deck.Draw(), deck.Draw()},
		ChannelID:  m.ChannelID,
		LastActive: time.Now(),
	}

	embed := renderBlackjackBoard(game, false)

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Hit",
					Style:    discordgo.SuccessButton,
					CustomID: "bj_hit",
				},
				discordgo.Button{
					Label:    "Stand",
					Style:    discordgo.DangerButton,
					CustomID: "bj_stand",
				},
			},
		},
	}

	msg, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Embed:      embed,
		Components: components,
	})
	if err == nil {
		game.MessageID = msg.ID
		b.blackjackGames.Store(m.Author.ID, game)

		playerTotal := calculateHand(game.PlayerHand)
		if playerTotal == 21 {
			game.mu.Lock()
			b.resolveBlackjackGame(s, game)
			game.mu.Unlock()
		}
	} else {
		_ = b.store.UpdateUserBalance(context.Background(), m.Author.ID, m.GuildID, user.Balance)
	}
}

func renderBlackjackBoard(game *BlackjackGame, gameOver bool) *discordgo.MessageEmbed {
	playerTotal := calculateHand(game.PlayerHand)
	dealerTotal := calculateHand(game.DealerHand)

	playerStr := ""
	for _, c := range game.PlayerHand {
		playerStr += c.String() + " "
	}

	dealerStr := ""
	if gameOver {
		for _, c := range game.DealerHand {
			dealerStr += c.String() + " "
		}
	} else {
		dealerStr = game.DealerHand[0].String() + " ❓"
	}

	embed := newEmbed("🃏 Blackjack", "Bet: **"+strconv.Itoa(game.Bet)+" coins**")

	embed.Fields = []*discordgo.MessageEmbedField{
		{Name: "Your Hand (" + strconv.Itoa(playerTotal) + ")", Value: playerStr, Inline: true},
	}

	if gameOver {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{Name: "Dealer Hand (" + strconv.Itoa(dealerTotal) + ")", Value: dealerStr, Inline: true})
	} else {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{Name: "Dealer Hand", Value: dealerStr, Inline: true})
	}

	return embed
}

func (b *Bot) handleBlackjackInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	customID := i.MessageComponentData().CustomID
	if customID != "bj_hit" && customID != "bj_stand" {
		return
	}

	val, exists := b.blackjackGames.Load(i.Member.User.ID)
	if !exists {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This game is no longer active or it's not yours!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	game := val.(*BlackjackGame)

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})

	game.mu.Lock()
	defer game.mu.Unlock()

	switch customID {
	case "bj_hit":
		game.PlayerHand = append(game.PlayerHand, game.Deck.Draw())
		playerTotal := calculateHand(game.PlayerHand)

		if playerTotal >= 21 {
			b.resolveBlackjackGame(s, game)
		} else {
			game.LastActive = time.Now()
			embed := renderBlackjackBoard(game, false)
			EditComplex(s, &discordgo.MessageEdit{
				Channel:    game.ChannelID,
				ID:         game.MessageID,
				Embed:      embed,
				Components: &i.Message.Components,
			})
		}
	case "bj_stand":
		b.resolveBlackjackGame(s, game)
	}
}

func (b *Bot) resolveBlackjackGame(s *discordgo.Session, game *BlackjackGame) {
	b.blackjackGames.Delete(game.UserID)

	playerTotal := calculateHand(game.PlayerHand)

	if playerTotal <= 21 {
		dealerTotal := calculateHand(game.DealerHand)
		for dealerTotal < 17 {
			game.DealerHand = append(game.DealerHand, game.Deck.Draw())
			dealerTotal = calculateHand(game.DealerHand)
		}
	}

	dealerTotal := calculateHand(game.DealerHand)
	embed := renderBlackjackBoard(game, true)

	payout := 0
	msg := ""

	if playerTotal > 21 {
		msg = "💥 You busted! Dealer wins."
	} else if dealerTotal > 21 {
		payout = game.Bet * 2
		msg = "🎉 Dealer busted! You won **" + strconv.Itoa(payout) + " coins**!"
	} else if playerTotal > dealerTotal {
		if playerTotal == 21 && len(game.PlayerHand) == 2 {
			payout = int(float64(game.Bet) * 2.5)
			msg = "🃏 **BLACKJACK!** You won **" + strconv.Itoa(payout) + " coins**!"
		} else {
			payout = game.Bet * 2
			msg = "🎉 You win! You won **" + strconv.Itoa(payout) + " coins**!"
		}
	} else if playerTotal < dealerTotal {
		msg = "😢 Dealer wins."
	} else {
		payout = game.Bet
		msg = "🤝 Push! Your bet was refunded."
	}

	embed.Description = msg

	if payout > 0 {
		user, _ := b.store.GetUser(context.Background(), game.UserID, game.GuildID)
		_ = b.store.UpdateUserBalance(context.Background(), game.UserID, game.GuildID, user.Balance+payout)
	}

	EditComplex(s, &discordgo.MessageEdit{
		Channel:    game.ChannelID,
		ID:         game.MessageID,
		Embed:      embed,
		Components: &[]discordgo.MessageComponent{},
	})
}
