package database

import (
	"context"
	"database/sql"
	"sync"
)

type configRepo struct {
	stmtGetPrefix *sql.Stmt
	stmtSetPrefix *sql.Stmt
	cache         sync.Map
}

func (c *configRepo) prepare(db *sql.DB) error {
	var err error

	c.stmtGetPrefix, err = db.Prepare(
		"SELECT prefix FROM server_config WHERE guild_id = ?",
	)
	if err != nil {
		return err
	}

	c.stmtSetPrefix, err = db.Prepare(
		"INSERT INTO server_config (guild_id, prefix) VALUES (?, ?) " +
			"ON CONFLICT(guild_id) DO UPDATE SET prefix = excluded.prefix",
	)
	if err != nil {
		_ = c.stmtGetPrefix.Close()
		c.stmtGetPrefix = nil
		return err
	}

	return nil
}

func (c *configRepo) GetPrefix(ctx context.Context, guildID string) (string, error) {
	if val, ok := c.cache.Load(guildID); ok {
		prefix, ok := val.(string)
		if ok {
			return prefix, nil
		}

		c.cache.Delete(guildID)
	}

	var prefix string
	if err := c.stmtGetPrefix.QueryRowContext(ctx, guildID).Scan(&prefix); err != nil {
		return "", err
	}

	c.cache.Store(guildID, prefix)
	return prefix, nil
}

func (c *configRepo) SetPrefix(ctx context.Context, guildID, prefix string) error {
	if _, err := c.stmtSetPrefix.ExecContext(ctx, guildID, prefix); err != nil {
		return err
	}

	c.cache.Store(guildID, prefix)
	return nil
}
