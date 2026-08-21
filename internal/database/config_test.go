package database

import (
	"context"
	"database/sql"
	"testing"
)

func newTestRepository(t *testing.T) *Repository {
	t.Helper()

	repo, err := NewRepository(":memory:")
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}

	t.Cleanup(func() {
		repo.Close(context.Background())
	})

	return repo
}

func TestConfigRepoGetPrefixMissingGuild(t *testing.T) {
	repo := newTestRepository(t)

	prefix, err := repo.GetPrefix(context.Background(), "missing-guild")
	if err != sql.ErrNoRows {
		t.Fatalf("expected sql.ErrNoRows, got prefix=%q err=%v", prefix, err)
	}

	if _, ok := repo.cache.Load("missing-guild"); ok {
		t.Fatal("missing guild must not be cached")
	}
}

func TestConfigRepoSetAndGetPrefix(t *testing.T) {
	repo := newTestRepository(t)

	ctx := context.Background()
	guildID := "1234567890"
	expectedPrefix := "!"

	if err := repo.SetPrefix(ctx, guildID, expectedPrefix); err != nil {
		t.Fatalf("failed to set prefix: %v", err)
	}

	prefix, err := repo.GetPrefix(ctx, guildID)
	if err != nil {
		t.Fatalf("failed to get prefix: %v", err)
	}

	if prefix != expectedPrefix {
		t.Errorf("expected prefix %q, got %q", expectedPrefix, prefix)
	}
}

func TestConfigRepoUpdatePrefix(t *testing.T) {
	repo := newTestRepository(t)

	ctx := context.Background()
	guildID := "1234567890"

	if err := repo.SetPrefix(ctx, guildID, "!"); err != nil {
		t.Fatalf("failed to set initial prefix: %v", err)
	}

	if err := repo.SetPrefix(ctx, guildID, "?"); err != nil {
		t.Fatalf("failed to update prefix: %v", err)
	}

	prefix, err := repo.GetPrefix(ctx, guildID)
	if err != nil {
		t.Fatalf("failed to get updated prefix: %v", err)
	}

	if prefix != "?" {
		t.Errorf("expected updated prefix %q, got %q", "?", prefix)
	}
}

func TestConfigRepoCacheAfterSet(t *testing.T) {
	repo := newTestRepository(t)

	ctx := context.Background()
	guildID := "cache-guild"
	expectedPrefix := "@"

	if err := repo.SetPrefix(ctx, guildID, expectedPrefix); err != nil {
		t.Fatalf("failed to set prefix: %v", err)
	}

	val, ok := repo.cache.Load(guildID)
	if !ok {
		t.Fatal("expected guild to be present in cache")
	}

	prefix, ok := val.(string)
	if !ok {
		t.Fatalf("expected cached value to be string, got %T", val)
	}

	if prefix != expectedPrefix {
		t.Errorf("expected cached prefix %q, got %q", expectedPrefix, prefix)
	}
}

func TestConfigRepoCacheRepopulation(t *testing.T) {
	repo := newTestRepository(t)

	ctx := context.Background()
	guildID := "repopulate-guild"
	expectedPrefix := "$"

	if err := repo.SetPrefix(ctx, guildID, expectedPrefix); err != nil {
		t.Fatalf("failed to set prefix: %v", err)
	}

	repo.cache.Delete(guildID)

	prefix, err := repo.GetPrefix(ctx, guildID)
	if err != nil {
		t.Fatalf("failed to get prefix after deleting cache: %v", err)
	}

	if prefix != expectedPrefix {
		t.Errorf("expected prefix %q, got %q", expectedPrefix, prefix)
	}

	val, ok := repo.cache.Load(guildID)
	if !ok {
		t.Fatal("expected cache to be repopulated")
	}

	cachedPrefix, ok := val.(string)
	if !ok {
		t.Fatalf("expected cached value to be string, got %T", val)
	}

	if cachedPrefix != expectedPrefix {
		t.Errorf("expected cached prefix %q, got %q", expectedPrefix, cachedPrefix)
	}
}

func TestConfigRepoEmptyPrefix(t *testing.T) {
	repo := newTestRepository(t)

	ctx := context.Background()
	guildID := "empty-prefix-guild"

	if err := repo.SetPrefix(ctx, guildID, ""); err != nil {
		t.Fatalf("failed to set empty prefix: %v", err)
	}

	prefix, err := repo.GetPrefix(ctx, guildID)
	if err != nil {
		t.Fatalf("expected empty prefix to be valid, got error: %v", err)
	}

	if prefix != "" {
		t.Errorf("expected empty prefix, got %q", prefix)
	}

	val, ok := repo.cache.Load(guildID)
	if !ok {
		t.Fatal("expected empty prefix to be cached")
	}

	cachedPrefix, ok := val.(string)
	if !ok {
		t.Fatalf("expected cached value to be string, got %T", val)
	}

	if cachedPrefix != "" {
		t.Errorf("expected cached prefix to be empty, got %q", cachedPrefix)
	}
}

func TestConfigRepoMultipleGuilds(t *testing.T) {
	repo := newTestRepository(t)

	ctx := context.Background()

	tests := []struct {
		guildID string
		prefix  string
	}{
		{
			guildID: "guild-1",
			prefix:  "!",
		},
		{
			guildID: "guild-2",
			prefix:  "?",
		},
		{
			guildID: "guild-3",
			prefix:  ">>",
		},
	}

	for _, tt := range tests {
		if err := repo.SetPrefix(ctx, tt.guildID, tt.prefix); err != nil {
			t.Fatalf(
				"failed to set prefix for guild %q: %v",
				tt.guildID,
				err,
			)
		}
	}

	for _, tt := range tests {
		prefix, err := repo.GetPrefix(ctx, tt.guildID)
		if err != nil {
			t.Fatalf(
				"failed to get prefix for guild %q: %v",
				tt.guildID,
				err,
			)
		}

		if prefix != tt.prefix {
			t.Errorf(
				"guild %q: expected prefix %q, got %q",
				tt.guildID,
				tt.prefix,
				prefix,
			)
		}
	}
}

func TestConfigRepoSetPrefixAfterMissingLookup(t *testing.T) {
	repo := newTestRepository(t)

	ctx := context.Background()
	guildID := "previously-missing-guild"

	_, err := repo.GetPrefix(ctx, guildID)
	if err != sql.ErrNoRows {
		t.Fatalf("expected initial sql.ErrNoRows, got %v", err)
	}

	if _, ok := repo.cache.Load(guildID); ok {
		t.Fatal("missing guild must not be cached")
	}

	if err := repo.SetPrefix(ctx, guildID, "!"); err != nil {
		t.Fatalf("failed to set prefix: %v", err)
	}

	prefix, err := repo.GetPrefix(ctx, guildID)
	if err != nil {
		t.Fatalf("failed to get prefix after SetPrefix: %v", err)
	}

	if prefix != "!" {
		t.Errorf("expected prefix %q, got %q", "!", prefix)
	}
}

func TestConfigRepoContextCancellation(t *testing.T) {
	repo := newTestRepository(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := repo.GetPrefix(ctx, "cancelled-guild"); err == nil {
		t.Fatal("expected GetPrefix to fail with cancelled context")
	}

	if err := repo.SetPrefix(ctx, "cancelled-guild", "!"); err == nil {
		t.Fatal("expected SetPrefix to fail with cancelled context")
	}
}

func TestConfigRepoPrepareStatements(t *testing.T) {
	repo := newTestRepository(t)

	if repo.stmtGetPrefix == nil {
		t.Fatal("expected stmtGetPrefix to be initialized")
	}

	if repo.stmtSetPrefix == nil {
		t.Fatal("expected stmtSetPrefix to be initialized")
	}
}

func TestConfigRepoDatabasePersistsAfterCacheDeletion(t *testing.T) {
	repo := newTestRepository(t)

	ctx := context.Background()
	guildID := "persistent-guild"
	expectedPrefix := "&"

	if err := repo.SetPrefix(ctx, guildID, expectedPrefix); err != nil {
		t.Fatalf("failed to set prefix: %v", err)
	}

	repo.cache.Delete(guildID)

	prefix, err := repo.GetPrefix(ctx, guildID)
	if err != nil {
		t.Fatalf("failed to read prefix from database: %v", err)
	}

	if prefix != expectedPrefix {
		t.Errorf("expected prefix %q, got %q", expectedPrefix, prefix)
	}
}
