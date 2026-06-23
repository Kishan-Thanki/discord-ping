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
	usersRepo
	configRepo
	remindersRepo
	moderationRepo
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
	if err := r.usersRepo.prepare(db); err != nil {
		return nil, err
	}
	if err := r.configRepo.prepare(db); err != nil {
		return nil, err
	}
	if err := r.remindersRepo.prepare(db); err != nil {
		return nil, err
	}
	if err := r.moderationRepo.prepare(db); err != nil {
		return nil, err
	}

	slog.Info("Database initialized successfully with WAL & Prepared Statements")
	return r, nil
}

// Close releases all prepared statements and closes the database connection.
func (r *Repository) Close(_ context.Context) {
	stmts := []*sql.Stmt{
		r.stmtGetUser, r.stmtInsertUser, r.stmtUpdateXP, r.stmtUpdateBal,
		r.stmtUpdateDaily, r.stmtGetLB, r.stmtGetPrefix, r.stmtSetPrefix,
		r.stmtCreateRem, r.stmtDeleteRem, r.stmtGetRem, r.stmtAddWarn, r.stmtGetWarnCnt,
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
		`CREATE TABLE IF NOT EXISTS users (
			user_id    TEXT NOT NULL,
			guild_id   TEXT NOT NULL,
			xp         INTEGER DEFAULT 0,
			level      INTEGER DEFAULT 0,
			balance    INTEGER DEFAULT 0,
			last_daily TEXT DEFAULT '',
			PRIMARY KEY (user_id, guild_id)
		);`,
		`CREATE TABLE IF NOT EXISTS server_config (
			guild_id        TEXT PRIMARY KEY,
			prefix          TEXT DEFAULT '!',
			welcome_channel TEXT DEFAULT '',
			log_channel     TEXT DEFAULT ''
		);`,
		`CREATE TABLE IF NOT EXISTS reminders (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id    TEXT NOT NULL,
			channel_id TEXT NOT NULL,
			guild_id   TEXT NOT NULL,
			message    TEXT NOT NULL,
			remind_at  TEXT NOT NULL,
			created_at TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS warnings (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id    TEXT NOT NULL,
			guild_id   TEXT NOT NULL,
			reason     TEXT NOT NULL,
			timestamp  TEXT NOT NULL
		);`,
	}

	for _, q := range queries {
		if _, err := r.db.ExecContext(ctx, q); err != nil {
			return err
		}
	}
	return nil
}
