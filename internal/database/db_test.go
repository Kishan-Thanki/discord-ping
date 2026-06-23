package database

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"
)

var repo *Repository

func TestMain(m *testing.M) {
	var err error
	repo, err = NewRepository(":memory:")
	if err != nil {
		panic(err)
	}
	code := m.Run()
	repo.Close(context.Background())
	os.Exit(code)
}

func TestUserCRUD(t *testing.T) {
	ctx := context.Background()
	userID := "123"
	guildID := "456"

	u, err := repo.GetUser(ctx, userID, guildID)
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}
	if u.XP != 0 || u.Level != 0 || u.Balance != 0 {
		t.Errorf("Expected new user to have 0 stats, got XP: %d, Level: %d, Balance: %d", u.XP, u.Level, u.Balance)
	}

	err = repo.UpdateUserXP(ctx, userID, guildID, 150, 1)
	if err != nil {
		t.Fatalf("UpdateUserXP failed: %v", err)
	}

	u, _ = repo.GetUser(ctx, userID, guildID)
	if u.XP != 150 || u.Level != 1 {
		t.Errorf("XP update failed. Expected 150/1, got %d/%d", u.XP, u.Level)
	}

	err = repo.UpdateUserBalance(ctx, userID, guildID, 500)
	if err != nil {
		t.Fatalf("UpdateUserBalance failed: %v", err)
	}

	u, _ = repo.GetUser(ctx, userID, guildID)
	if u.Balance != 500 {
		t.Errorf("Balance update failed. Expected 500, got %d", u.Balance)
	}

	nowStr := time.Now().Format(time.RFC3339)
	err = repo.UpdateUserDaily(ctx, userID, guildID, 600, nowStr)
	if err != nil {
		t.Fatalf("UpdateUserDaily failed: %v", err)
	}

	u, _ = repo.GetUser(ctx, userID, guildID)
	if u.Balance != 600 || u.LastDaily != nowStr {
		t.Errorf("Daily update failed. Expected 600/%s, got %d/%s", nowStr, u.Balance, u.LastDaily)
	}
}

func TestLeaderboard(t *testing.T) {
	ctx := context.Background()
	guildID := "leaderboard_test"

	repo.GetUser(ctx, "user1", guildID)
	repo.GetUser(ctx, "user2", guildID)
	repo.GetUser(ctx, "user3", guildID)

	repo.UpdateUserXP(ctx, "user1", guildID, 100, 1)
	repo.UpdateUserXP(ctx, "user2", guildID, 300, 3)
	repo.UpdateUserXP(ctx, "user3", guildID, 200, 2)

	lb, err := repo.GetLeaderboard(ctx, guildID, 10)
	if err != nil {
		t.Fatalf("GetLeaderboard failed: %v", err)
	}

	if len(lb) != 3 {
		t.Fatalf("Expected 3 users in leaderboard, got %d", len(lb))
	}

	if lb[0].UserID != "user2" || lb[1].UserID != "user3" || lb[2].UserID != "user1" {
		t.Errorf("Leaderboard order incorrect: %v", lb)
	}
}

func TestPrefix(t *testing.T) {
	ctx := context.Background()
	guildID := "prefix_test"

	p, err := repo.GetPrefix(ctx, guildID)
	if err != sql.ErrNoRows {
		t.Fatalf("Expected ErrNoRows for new guild, got: %v", err)
	}
	if p != "" {
		t.Errorf("Expected empty prefix, got '%s'", p)
	}

	err = repo.SetPrefix(ctx, guildID, "?")
	if err != nil {
		t.Fatalf("SetPrefix failed: %v", err)
	}

	p, _ = repo.GetPrefix(ctx, guildID)
	if p != "?" {
		t.Errorf("Expected prefix '?', got '%s'", p)
	}

	err = repo.SetPrefix(ctx, guildID, "$")
	if err != nil {
		t.Fatalf("SetPrefix failed: %v", err)
	}

	p, _ = repo.GetPrefix(ctx, guildID)
	if p != "$" {
		t.Errorf("Expected prefix '$', got '%s'", p)
	}
}

func TestReminders(t *testing.T) {
	ctx := context.Background()
	id, err := repo.CreateReminder(ctx, "123", "chan", "guild", "hello", time.Now())
	if err != nil {
		t.Fatalf("CreateReminder failed: %v", err)
	}

	pending, err := repo.GetPendingReminders(ctx)
	if err != nil {
		t.Fatalf("GetPendingReminders failed: %v", err)
	}

	found := false
	for _, r := range pending {
		if r.ID == id {
			found = true
			if r.Message != "hello" {
				t.Errorf("Expected message 'hello', got '%s'", r.Message)
			}
			break
		}
	}

	if !found {
		t.Errorf("Created reminder not found in pending list")
	}

	err = repo.DeleteReminder(ctx, id)
	if err != nil {
		t.Fatalf("DeleteReminder failed: %v", err)
	}
}
