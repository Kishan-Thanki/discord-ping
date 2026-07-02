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
	c.stmtGetPrefix, err = db.Prepare("SELECT prefix FROM server_config WHERE guild_id = ?")
	if err != nil {
		return err
	}
	c.stmtSetPrefix, err = db.Prepare("INSERT INTO server_config (guild_id, prefix) VALUES (?, ?) ON CONFLICT(guild_id) DO UPDATE SET prefix = ?")
	return err
}

func (c *configRepo) GetPrefix(ctx context.Context, guildID string) (string, error) {
	if val, ok := c.cache.Load(guildID); ok {
		p := val.(string)
		if p == "" {
			return "", sql.ErrNoRows
		}
		return p, nil
	}

	var prefix string
	err := c.stmtGetPrefix.QueryRowContext(ctx, guildID).Scan(&prefix)
	if err == sql.ErrNoRows {
		c.cache.Store(guildID, "")
		return "", err
	}
	if err != nil {
		return "", err
	}

	c.cache.Store(guildID, prefix)
	return prefix, nil
}

func (c *configRepo) SetPrefix(ctx context.Context, guildID, prefix string) error {
	_, err := c.stmtSetPrefix.ExecContext(ctx, guildID, prefix, prefix)
	if err == nil {
		c.cache.Store(guildID, prefix)
	}
	return err
}
