package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	t.Setenv("TOKEN", "test-token")
	t.Setenv("BOT_PREFIX", "!")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned unexpected error: %v", err)
	}

	if cfg == nil {
		t.Fatal("Load() returned nil config")
	}

	if cfg.Token != "test-token" {
		t.Errorf("expected token %q, got %q", "test-token", cfg.Token)
	}

	if cfg.BotPrefix != "!" {
		t.Errorf("expected prefix %q, got %q", "!", cfg.BotPrefix)
	}
}

func TestLoadMissingToken(t *testing.T) {
	t.Setenv("TOKEN", "")
	t.Setenv("BOT_PREFIX", "!")

	cfg, err := Load()
	if err == nil {
		t.Fatal("Load() expected an error when TOKEN is missing")
	}

	if cfg != nil {
		t.Errorf("expected nil config, got %#v", cfg)
	}

	expected := "required environment variable TOKEN is not set or empty"
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestLoadMissingPrefix(t *testing.T) {
	t.Setenv("TOKEN", "test-token")
	t.Setenv("BOT_PREFIX", "")

	cfg, err := Load()
	if err == nil {
		t.Fatal("Load() expected an error when BOT_PREFIX is missing")
	}

	if cfg != nil {
		t.Errorf("expected nil config, got %#v", cfg)
	}

	expected := "required environment variable BOT_PREFIX is not set or empty"
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestLoadWhitespaceToken(t *testing.T) {
	t.Setenv("TOKEN", "   ")
	t.Setenv("BOT_PREFIX", "!")

	cfg, err := Load()
	if err == nil {
		t.Fatal("Load() expected an error for whitespace-only TOKEN")
	}

	if cfg != nil {
		t.Errorf("expected nil config, got %#v", cfg)
	}
}

func TestLoadWhitespacePrefix(t *testing.T) {
	t.Setenv("TOKEN", "test-token")
	t.Setenv("BOT_PREFIX", "   ")

	cfg, err := Load()
	if err == nil {
		t.Fatal("Load() expected an error for whitespace-only BOT_PREFIX")
	}

	if cfg != nil {
		t.Errorf("expected nil config, got %#v", cfg)
	}
}

func TestLoadTrimsValues(t *testing.T) {
	t.Setenv("TOKEN", "  test-token  ")
	t.Setenv("BOT_PREFIX", "  !  ")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned unexpected error: %v", err)
	}

	if cfg == nil {
		t.Fatal("Load() returned nil config")
	}

	if cfg.Token != "test-token" {
		t.Errorf("expected trimmed token %q, got %q", "test-token", cfg.Token)
	}

	if cfg.BotPrefix != "!" {
		t.Errorf("expected trimmed prefix %q, got %q", "!", cfg.BotPrefix)
	}
}

func TestLoadPreservesValidValues(t *testing.T) {
	tests := []struct {
		name   string
		token  string
		prefix string
	}{
		{
			name:   "single character prefix",
			token:  "token-123",
			prefix: "!",
		},
		{
			name:   "multi character prefix",
			token:  "token-456",
			prefix: ">>",
		},
		{
			name:   "unicode prefix",
			token:  "token-789",
			prefix: "§",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("TOKEN", tt.token)
			t.Setenv("BOT_PREFIX", tt.prefix)

			cfg, err := Load()
			if err != nil {
				t.Fatalf("Load() returned unexpected error: %v", err)
			}

			if cfg == nil {
				t.Fatal("Load() returned nil config")
			}

			if cfg.Token != tt.token {
				t.Errorf("expected token %q, got %q", tt.token, cfg.Token)
			}

			if cfg.BotPrefix != tt.prefix {
				t.Errorf("expected prefix %q, got %q", tt.prefix, cfg.BotPrefix)
			}
		})
	}
}

func TestVersion(t *testing.T) {
	const expectedVersion = "v1.2.0"

	if Version != expectedVersion {
		t.Errorf("expected Version %q, got %q", expectedVersion, Version)
	}
}
