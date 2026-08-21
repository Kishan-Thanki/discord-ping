package database

import (
	"context"
	"testing"
)

func TestNewRepository(t *testing.T) {
	repo := newTestRepository(t)

	if repo == nil {
		t.Fatal("expected repository, got nil")
	}

	if repo.db == nil {
		t.Fatal("expected database connection to be initialized")
	}

	if repo.stmtGetPrefix == nil {
		t.Fatal("expected get-prefix statement to be initialized")
	}

	if repo.stmtSetPrefix == nil {
		t.Fatal("expected set-prefix statement to be initialized")
	}
}

func TestRepositoryCreatesServerConfigTable(t *testing.T) {
	repo := newTestRepository(t)

	var tableName string
	err := repo.db.QueryRowContext(
		context.Background(),
		"SELECT name FROM sqlite_master WHERE type = 'table' AND name = 'server_config'",
	).Scan(&tableName)
	if err != nil {
		t.Fatalf("failed to verify server_config table: %v", err)
	}

	if tableName != "server_config" {
		t.Fatalf("expected table name %q, got %q", "server_config", tableName)
	}
}

func TestRepositoryDefaultPrefix(t *testing.T) {
	repo := newTestRepository(t)

	ctx := context.Background()
	guildID := "default-prefix-guild"

	_, err := repo.db.ExecContext(
		ctx,
		"INSERT INTO server_config (guild_id) VALUES (?)",
		guildID,
	)
	if err != nil {
		t.Fatalf("failed to insert guild without prefix: %v", err)
	}

	var prefix string
	if err := repo.db.QueryRowContext(
		ctx,
		"SELECT prefix FROM server_config WHERE guild_id = ?",
		guildID,
	).Scan(&prefix); err != nil {
		t.Fatalf("failed to read default prefix: %v", err)
	}

	if prefix != "!" {
		t.Errorf("expected default prefix %q, got %q", "!", prefix)
	}
}

func TestRepositoryPrefixCannotBeNull(t *testing.T) {
	repo := newTestRepository(t)

	_, err := repo.db.ExecContext(
		context.Background(),
		"INSERT INTO server_config (guild_id, prefix) VALUES (?, NULL)",
		"null-prefix-guild",
	)
	if err == nil {
		t.Fatal("expected inserting NULL prefix to fail")
	}
}

func TestRepositoryClose(t *testing.T) {
	repo, err := NewRepository(":memory:")
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}

	repo.Close(context.Background())

	if repo.db != nil {
		t.Error("expected repository database handle to be nil after Close")
	}
}

func TestRepositoryCloseNil(t *testing.T) {
	var repo *Repository

	repo.Close(context.Background())
}
