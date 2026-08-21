package database

import (
	"context"
	"database/sql"
	"log/slog"

	_ "modernc.org/sqlite"
)

type Repository struct {
	db *sql.DB
	configRepo
}

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

	for _, pragma := range pragmas {
		if _, err := db.ExecContext(context.Background(), pragma); err != nil {
			_ = db.Close()
			return nil, err
		}
	}

	repo := &Repository{db: db}

	if err := repo.createTables(context.Background()); err != nil {
		_ = db.Close()
		return nil, err
	}

	if err := repo.configRepo.prepare(db); err != nil {
		_ = db.Close()
		return nil, err
	}

	slog.Info("database initialized successfully")

	return repo, nil
}

func (r *Repository) Close(_ context.Context) {
	if r == nil {
		return
	}

	stmts := []*sql.Stmt{
		r.stmtGetPrefix,
		r.stmtSetPrefix,
	}

	for _, stmt := range stmts {
		if stmt != nil {
			_ = stmt.Close()
		}
	}

	if r.db != nil {
		slog.Info("closing database connection")
		_ = r.db.Close()
		r.db = nil
	}
}

func (r *Repository) createTables(ctx context.Context) error {
	const query = `
		CREATE TABLE IF NOT EXISTS server_config (
			guild_id TEXT PRIMARY KEY,
			prefix   TEXT NOT NULL DEFAULT '!'
		);
	`

	_, err := r.db.ExecContext(ctx, query)
	return err
}
