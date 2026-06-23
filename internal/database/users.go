package database

import (
	"context"
	"database/sql"
)

// User represents a user record in the database.
type User struct {
	UserID    string
	GuildID   string
	XP        int
	Level     int
	Balance   int
	LastDaily string
}

type usersRepo struct {
	stmtGetUser     *sql.Stmt
	stmtInsertUser  *sql.Stmt
	stmtUpdateXP    *sql.Stmt
	stmtUpdateBal   *sql.Stmt
	stmtUpdateDaily *sql.Stmt
	stmtGetLB       *sql.Stmt
}

func (u *usersRepo) prepare(db *sql.DB) error {
	var err error
	u.stmtGetUser, err = db.Prepare("SELECT xp, level, balance, last_daily FROM users WHERE user_id = ? AND guild_id = ?")
	if err != nil {
		return err
	}
	u.stmtInsertUser, err = db.Prepare("INSERT INTO users (user_id, guild_id) VALUES (?, ?)")
	if err != nil {
		return err
	}
	u.stmtUpdateXP, err = db.Prepare("UPDATE users SET xp = ?, level = ? WHERE user_id = ? AND guild_id = ?")
	if err != nil {
		return err
	}
	u.stmtUpdateBal, err = db.Prepare("UPDATE users SET balance = ? WHERE user_id = ? AND guild_id = ?")
	if err != nil {
		return err
	}
	u.stmtUpdateDaily, err = db.Prepare("UPDATE users SET balance = ?, last_daily = ? WHERE user_id = ? AND guild_id = ?")
	if err != nil {
		return err
	}
	u.stmtGetLB, err = db.Prepare("SELECT user_id, xp, level FROM users WHERE guild_id = ? ORDER BY xp DESC LIMIT ?")
	return err
}

func (u *usersRepo) GetUser(ctx context.Context, userID, guildID string) (*User, error) {
	usr := &User{UserID: userID, GuildID: guildID}

	err := u.stmtGetUser.QueryRowContext(ctx, userID, guildID).Scan(&usr.XP, &usr.Level, &usr.Balance, &usr.LastDaily)
	if err == sql.ErrNoRows {
		_, err = u.stmtInsertUser.ExecContext(ctx, userID, guildID)
		if err != nil {
			return nil, err
		}
		return usr, nil
	}
	if err != nil {
		return nil, err
	}
	return usr, nil
}

func (u *usersRepo) UpdateUserXP(ctx context.Context, userID, guildID string, xp, level int) error {
	_, err := u.stmtUpdateXP.ExecContext(ctx, xp, level, userID, guildID)
	return err
}

func (u *usersRepo) UpdateUserBalance(ctx context.Context, userID, guildID string, balance int) error {
	_, err := u.stmtUpdateBal.ExecContext(ctx, balance, userID, guildID)
	return err
}

func (u *usersRepo) UpdateUserDaily(ctx context.Context, userID, guildID string, balance int, lastDaily string) error {
	_, err := u.stmtUpdateDaily.ExecContext(ctx, balance, lastDaily, userID, guildID)
	return err
}

func (u *usersRepo) GetLeaderboard(ctx context.Context, guildID string, limit int) ([]*User, error) {
	rows, err := u.stmtGetLB.QueryContext(ctx, guildID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		usr := &User{GuildID: guildID}
		if err := rows.Scan(&usr.UserID, &usr.XP, &usr.Level); err != nil {
			return nil, err
		}
		users = append(users, usr)
	}
	return users, rows.Err()
}
