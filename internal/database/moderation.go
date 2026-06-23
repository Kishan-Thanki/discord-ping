package database

import (
	"context"
	"database/sql"
	"time"
)

type moderationRepo struct {
	stmtAddWarn    *sql.Stmt
	stmtGetWarnCnt *sql.Stmt
}

func (m *moderationRepo) prepare(db *sql.DB) error {
	var err error
	m.stmtAddWarn, err = db.Prepare("INSERT INTO warnings (user_id, guild_id, reason, timestamp) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	m.stmtGetWarnCnt, err = db.Prepare("SELECT COUNT(*) FROM warnings WHERE user_id = ? AND guild_id = ?")
	return err
}

func (m *moderationRepo) AddWarning(ctx context.Context, userID, guildID, reason string) error {
	_, err := m.stmtAddWarn.ExecContext(ctx, userID, guildID, reason, time.Now().Format(time.RFC3339))
	return err
}

func (m *moderationRepo) GetWarningCount(ctx context.Context, userID, guildID string) (int, error) {
	var count int
	err := m.stmtGetWarnCnt.QueryRowContext(ctx, userID, guildID).Scan(&count)
	return count, err
}
