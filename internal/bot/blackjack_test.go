package bot

import (
	"testing"
)

func TestCalculateHand(t *testing.T) {
	tests := []struct {
		name     string
		hand     []Card
		expected int
	}{
		{
			name: "Normal Cards",
			hand: []Card{
				{Rank: "10", Value: 10},
				{Rank: "7", Value: 7},
			},
			expected: 17,
		},
		{
			name: "Face Cards",
			hand: []Card{
				{Rank: "K", Value: 10},
				{Rank: "J", Value: 10},
			},
			expected: 20,
		},
		{
			name: "Blackjack (Ace + 10)",
			hand: []Card{
				{Rank: "A", Value: 11},
				{Rank: "K", Value: 10},
			},
			expected: 21,
		},
		{
			name: "Ace Downgrade to 1",
			hand: []Card{
				{Rank: "A", Value: 11},
				{Rank: "10", Value: 10},
				{Rank: "5", Value: 5},
			},
			expected: 16, // 11 + 10 + 5 = 26 -> Ace drops to 1 -> 1 + 10 + 5 = 16
		},
		{
			name: "Double Ace",
			hand: []Card{
				{Rank: "A", Value: 11},
				{Rank: "A", Value: 11},
			},
			expected: 12, // 11 + 11 = 22 -> One ace drops to 1 -> 1 + 11 = 12
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			total := calculateHand(tt.hand)
			if total != tt.expected {
				t.Errorf("calculateHand() = %v, want %v", total, tt.expected)
			}
		})
	}
}

func TestNewDeck(t *testing.T) {
	deck := newDeck()
	if len(deck.Cards) != 52 {
		t.Errorf("newDeck() generated %v cards, want 52", len(deck.Cards))
	}

	// Verify all suits exist
	suitsMap := make(map[string]int)
	for _, c := range deck.Cards {
		suitsMap[c.Suit]++
	}

	if len(suitsMap) != 4 {
		t.Errorf("newDeck() generated %v suits, want 4", len(suitsMap))
	}
}
