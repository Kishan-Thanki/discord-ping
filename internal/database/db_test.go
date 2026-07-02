package database

import (
	"context"
	"database/sql"
	"testing"
)

func TestNewRepository(t *testing.T) {

	repo, err := NewRepository(":memory:")
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}
	defer repo.Close(context.Background())

	if repo.db == nil {
		t.Error("expected database connection to be initialized")
	}
}

func TestPrefixOperations(t *testing.T) {
	repo, err := NewRepository(":memory:")
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}
	defer repo.Close(context.Background())

	ctx := context.Background()
	guildID := "1234567890"

	prefix, err := repo.GetPrefix(ctx, guildID)
	if err == nil {
		t.Errorf("expected error (sql.ErrNoRows) for non-existent guild, got prefix: %q", prefix)
	}
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got: %v", err)
	}

	err = repo.SetPrefix(ctx, guildID, "?")
	if err != nil {
		t.Fatalf("failed to set prefix: %v", err)
	}

	prefix, err = repo.GetPrefix(ctx, guildID)
	if err != nil {
		t.Fatalf("failed to get prefix: %v", err)
	}
	if prefix != "?" {
		t.Errorf("expected prefix '?', got: %q", prefix)
	}

	err = repo.SetPrefix(ctx, guildID, "$")
	if err != nil {
		t.Fatalf("failed to update prefix: %v", err)
	}

	prefix, err = repo.GetPrefix(ctx, guildID)
	if err != nil {
		t.Fatalf("failed to get prefix after update: %v", err)
	}
	if prefix != "$" {
		t.Errorf("expected updated prefix '$', got: %q", prefix)
	}
}

func TestPrefixCaching(t *testing.T) {
	repo, err := NewRepository(":memory:")
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}
	defer repo.Close(context.Background())

	ctx := context.Background()
	guildID := "cache-guild"

	_, err = repo.GetPrefix(ctx, guildID)
	if err != sql.ErrNoRows {
		t.Fatalf("expected ErrNoRows, got: %v", err)
	}

	val, ok := repo.cache.Load(guildID)
	if !ok {
		t.Error("expected guild to be in cache")
	}
	if val.(string) != "" {
		t.Errorf("expected cached prefix to be empty string, got: %q", val)
	}

	err = repo.SetPrefix(ctx, guildID, "@")
	if err != nil {
		t.Fatalf("failed to set prefix: %v", err)
	}

	val, ok = repo.cache.Load(guildID)
	if !ok {
		t.Error("expected guild to be in cache after SetPrefix")
	}
	if val.(string) != "@" {
		t.Errorf("expected cached prefix to be '@', got: %q", val)
	}

	repo.cache.Delete(guildID)

	p, err := repo.GetPrefix(ctx, guildID)
	if err != nil {
		t.Fatalf("failed to get prefix: %v", err)
	}
	if p != "@" {
		t.Errorf("expected prefix '@', got: %q", p)
	}

	val, ok = repo.cache.Load(guildID)
	if !ok {
		t.Error("expected cache to be repopulated after GetPrefix miss")
	}
	if val.(string) != "@" {
		t.Errorf("expected cached prefix to be '@' after repopulation, got: %q", val)
	}
}
