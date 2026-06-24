package database

import (
	"context"
	"database/sql"
	"log/slog"

	_ "modernc.org/sqlite"
)

// Repository encapsulates all database access.
// Every piece of state — the connection AND the prepared statements — lives on this struct.
// There are ZERO package-level variables.
type Repository struct {
	db *sql.DB
	configRepo
}

// NewRepository opens a SQLite database, applies performance pragmas,
// creates tables, and prepares all statements.
func NewRepository(dataSourceName string) (*Repository, error) {
	db, err := sql.Open("sqlite", dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	pragmas := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
		"PRAGMA temp_store = MEMORY;",
		"PRAGMA cache_size = -16000;",
	}
	for _, p := range pragmas {
		if _, err := db.ExecContext(context.Background(), p); err != nil {
			return nil, err
		}
	}

	r := &Repository{db: db}

	if err := r.createTables(context.Background()); err != nil {
		return nil, err
	}
	if err := r.configRepo.prepare(db); err != nil {
		return nil, err
	}

	slog.Info("Database initialized successfully with WAL & Prepared Statements")
	return r, nil
}

// Close releases all prepared statements and closes the database connection.
func (r *Repository) Close(_ context.Context) {
	stmts := []*sql.Stmt{
		r.stmtGetPrefix, r.stmtSetPrefix,
	}
	for _, s := range stmts {
		if s != nil {
			s.Close()
		}
	}
	if r.db != nil {
		slog.Info("Closing database connection")
		r.db.Close()
	}
}

func (r *Repository) createTables(ctx context.Context) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS server_config (
			guild_id        TEXT PRIMARY KEY,
			prefix          TEXT DEFAULT '!'
		);`,
	}

	for _, q := range queries {
		if _, err := r.db.ExecContext(ctx, q); err != nil {
			return err
		}
	}
	return nil
}
