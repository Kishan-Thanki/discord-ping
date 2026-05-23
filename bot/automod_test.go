package bot

import (
	"testing"
)

func TestBadWordsRegex(t *testing.T) {
	tests := []struct {
		message string
		blocked bool
	}{
		// Clean messages
		{"Hello world!", false},
		{"Can I get some free nitr0?", false}, // Typo, doesn't match exactly
		{"I am having a fantastic day", false},

		// Bad words
		{"You are a bitch", true},
		{"fUck off", true}, // Case insensitive
		{"this is shit", true},
		{"asshole", true},
		{"Kill yourself KYS now", true},

		// Spam links
		{"Click here: discord.gift/123456", true},
		{"I have free nitro for you", true},
		{"discord.gift", true},
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			result := badWordsRegex.MatchString(tt.message)
			if result != tt.blocked {
				t.Errorf("Message %q: expected blocked=%v, got %v", tt.message, tt.blocked, result)
			}
		})
	}
}
