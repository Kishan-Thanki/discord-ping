package database

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	_ "modernc.org/sqlite"
)

// Repository encapsulates all database access. Every piece of state —
// the connection AND the prepared statements — lives on this struct.
// There are ZERO package-level variables.
type Repository struct {
	db              *sql.DB
	stmtGetUser     *sql.Stmt
	stmtInsertUser  *sql.Stmt
	stmtUpdateXP    *sql.Stmt
	stmtUpdateBal   *sql.Stmt
	stmtUpdateDaily *sql.Stmt
	stmtGetLB       *sql.Stmt
	stmtGetPrefix   *sql.Stmt
	stmtSetPrefix   *sql.Stmt
	stmtCreateRem   *sql.Stmt
	stmtDeleteRem   *sql.Stmt
	stmtGetRem      *sql.Stmt
	stmtAddWarn     *sql.Stmt
	stmtGetWarnCnt  *sql.Stmt
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
	if err := r.prepareStatements(); err != nil {
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

func (r *Repository) prepareStatements() error {
	var err error

	r.stmtGetUser, err = r.db.Prepare("SELECT xp, level, balance, last_daily FROM users WHERE user_id = ? AND guild_id = ?")
	if err != nil {
		return err
	}
	r.stmtInsertUser, err = r.db.Prepare("INSERT INTO users (user_id, guild_id) VALUES (?, ?)")
	if err != nil {
		return err
	}
	r.stmtUpdateXP, err = r.db.Prepare("UPDATE users SET xp = ?, level = ? WHERE user_id = ? AND guild_id = ?")
	if err != nil {
		return err
	}
	r.stmtUpdateBal, err = r.db.Prepare("UPDATE users SET balance = ? WHERE user_id = ? AND guild_id = ?")
	if err != nil {
		return err
	}
	r.stmtUpdateDaily, err = r.db.Prepare("UPDATE users SET balance = ?, last_daily = ? WHERE user_id = ? AND guild_id = ?")
	if err != nil {
		return err
	}
	r.stmtGetLB, err = r.db.Prepare("SELECT user_id, xp, level FROM users WHERE guild_id = ? ORDER BY xp DESC LIMIT ?")
	if err != nil {
		return err
	}
	r.stmtGetPrefix, err = r.db.Prepare("SELECT prefix FROM server_config WHERE guild_id = ?")
	if err != nil {
		return err
	}
	r.stmtSetPrefix, err = r.db.Prepare("INSERT INTO server_config (guild_id, prefix) VALUES (?, ?) ON CONFLICT(guild_id) DO UPDATE SET prefix = ?")
	if err != nil {
		return err
	}
	r.stmtCreateRem, err = r.db.Prepare("INSERT INTO reminders (user_id, channel_id, guild_id, message, remind_at, created_at) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	r.stmtDeleteRem, err = r.db.Prepare("DELETE FROM reminders WHERE id = ?")
	if err != nil {
		return err
	}
	r.stmtGetRem, err = r.db.Prepare("SELECT id, user_id, channel_id, guild_id, message, remind_at FROM reminders")
	if err != nil {
		return err
	}
	r.stmtAddWarn, err = r.db.Prepare("INSERT INTO warnings (user_id, guild_id, reason, timestamp) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	r.stmtGetWarnCnt, err = r.db.Prepare("SELECT COUNT(*) FROM warnings WHERE user_id = ? AND guild_id = ?")
	if err != nil {
		return err
	}

	return nil
}

// ---------------------------------------------------------------------------
// Domain Types
// ---------------------------------------------------------------------------

// User represents a user record in the database.
type User struct {
	UserID    string
	GuildID   string
	XP        int
	Level     int
	Balance   int
	LastDaily string
}

// Reminder represents a scheduled reminder.
type Reminder struct {
	ID        int
	UserID    string
	ChannelID string
	GuildID   string
	Message   string
	RemindAt  time.Time
}

// ---------------------------------------------------------------------------
// Query Methods — every method uses the context-aware variants
// ---------------------------------------------------------------------------

func (r *Repository) GetUser(ctx context.Context, userID, guildID string) (*User, error) {
	u := &User{UserID: userID, GuildID: guildID}

	err := r.stmtGetUser.QueryRowContext(ctx, userID, guildID).Scan(&u.XP, &u.Level, &u.Balance, &u.LastDaily)
	if err == sql.ErrNoRows {
		_, err = r.stmtInsertUser.ExecContext(ctx, userID, guildID)
		if err != nil {
			return nil, err
		}
		return u, nil
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *Repository) UpdateUserXP(ctx context.Context, userID, guildID string, xp, level int) error {
	_, err := r.stmtUpdateXP.ExecContext(ctx, xp, level, userID, guildID)
	return err
}

func (r *Repository) UpdateUserBalance(ctx context.Context, userID, guildID string, balance int) error {
	_, err := r.stmtUpdateBal.ExecContext(ctx, balance, userID, guildID)
	return err
}

func (r *Repository) UpdateUserDaily(ctx context.Context, userID, guildID string, balance int, lastDaily string) error {
	_, err := r.stmtUpdateDaily.ExecContext(ctx, balance, lastDaily, userID, guildID)
	return err
}

func (r *Repository) GetLeaderboard(ctx context.Context, guildID string, limit int) ([]*User, error) {
	rows, err := r.stmtGetLB.QueryContext(ctx, guildID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		u := &User{GuildID: guildID}
		if err := rows.Scan(&u.UserID, &u.XP, &u.Level); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *Repository) GetPrefix(ctx context.Context, guildID string) (string, error) {
	var prefix string
	err := r.stmtGetPrefix.QueryRowContext(ctx, guildID).Scan(&prefix)
	if err == sql.ErrNoRows {
		return "!", nil
	}
	if err != nil {
		return "", err
	}
	return prefix, nil
}

func (r *Repository) SetPrefix(ctx context.Context, guildID, prefix string) error {
	_, err := r.stmtSetPrefix.ExecContext(ctx, guildID, prefix, prefix)
	return err
}

func (r *Repository) CreateReminder(ctx context.Context, userID, channelID, guildID, message string, remindAt time.Time) (int, error) {
	res, err := r.stmtCreateRem.ExecContext(ctx, userID, channelID, guildID, message, remindAt.Format(time.RFC3339), time.Now().Format(time.RFC3339))
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *Repository) DeleteReminder(ctx context.Context, id int) error {
	_, err := r.stmtDeleteRem.ExecContext(ctx, id)
	return err
}

func (r *Repository) GetPendingReminders(ctx context.Context) ([]*Reminder, error) {
	rows, err := r.stmtGetRem.QueryContext(ctx)
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

func (r *Repository) AddWarning(ctx context.Context, userID, guildID, reason string) error {
	_, err := r.stmtAddWarn.ExecContext(ctx, userID, guildID, reason, time.Now().Format(time.RFC3339))
	return err
}

func (r *Repository) GetWarningCount(ctx context.Context, userID, guildID string) (int, error) {
	var count int
	err := r.stmtGetWarnCnt.QueryRowContext(ctx, userID, guildID).Scan(&count)
	return count, err
}
