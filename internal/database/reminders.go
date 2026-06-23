package database

import (
	"context"
	"database/sql"
	"time"
)

// Reminder represents a scheduled reminder.
type Reminder struct {
	ID        int
	UserID    string
	ChannelID string
	GuildID   string
	Message   string
	RemindAt  time.Time
}

type remindersRepo struct {
	stmtCreateRem *sql.Stmt
	stmtDeleteRem *sql.Stmt
	stmtGetRem    *sql.Stmt
}

func (rm *remindersRepo) prepare(db *sql.DB) error {
	var err error
	rm.stmtCreateRem, err = db.Prepare("INSERT INTO reminders (user_id, channel_id, guild_id, message, remind_at, created_at) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	rm.stmtDeleteRem, err = db.Prepare("DELETE FROM reminders WHERE id = ?")
	if err != nil {
		return err
	}
	rm.stmtGetRem, err = db.Prepare("SELECT id, user_id, channel_id, guild_id, message, remind_at FROM reminders")
	return err
}

func (rm *remindersRepo) CreateReminder(ctx context.Context, userID, channelID, guildID, message string, remindAt time.Time) (int, error) {
	res, err := rm.stmtCreateRem.ExecContext(ctx, userID, channelID, guildID, message, remindAt.Format(time.RFC3339), time.Now().Format(time.RFC3339))
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (rm *remindersRepo) DeleteReminder(ctx context.Context, id int) error {
	_, err := rm.stmtDeleteRem.ExecContext(ctx, id)
	return err
}

func (rm *remindersRepo) GetPendingReminders(ctx context.Context) ([]*Reminder, error) {
	rows, err := rm.stmtGetRem.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reminders []*Reminder
	for rows.Next() {
		rem := &Reminder{}
		var remindAtStr string
		if err := rows.Scan(&rem.ID, &rem.UserID, &rem.ChannelID, &rem.GuildID, &rem.Message, &remindAtStr); err != nil {
			return nil, err
		}
		t, err := time.Parse(time.RFC3339, remindAtStr)
		if err != nil {
			continue
		}
		rem.RemindAt = t
		reminders = append(reminders, rem)
	}
	return reminders, rows.Err()
}
