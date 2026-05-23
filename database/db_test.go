package database

import (
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	err := InitDB(":memory:")
	if err != nil {
		panic(err)
	}
	code := m.Run()
	CloseDB()
	os.Exit(code)
}

func TestUserCRUD(t *testing.T) {
	userID := "123"
	guildID := "456"

	// 1. Get User (should create if not exists)
	u, err := GetUser(userID, guildID)
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}
	if u.XP != 0 || u.Level != 0 || u.Balance != 0 {
		t.Errorf("Expected new user to have 0 stats, got XP: %d, Level: %d, Balance: %d", u.XP, u.Level, u.Balance)
	}

	// 2. Update XP
	err = UpdateUserXP(userID, guildID, 150, 1)
	if err != nil {
		t.Fatalf("UpdateUserXP failed: %v", err)
	}

	u, _ = GetUser(userID, guildID)
	if u.XP != 150 || u.Level != 1 {
		t.Errorf("XP update failed. Expected 150/1, got %d/%d", u.XP, u.Level)
	}

	// 3. Update Balance
	err = UpdateUserBalance(userID, guildID, 500)
	if err != nil {
		t.Fatalf("UpdateUserBalance failed: %v", err)
	}

	u, _ = GetUser(userID, guildID)
	if u.Balance != 500 {
		t.Errorf("Balance update failed. Expected 500, got %d", u.Balance)
	}

	// 4. Update Daily
	nowStr := time.Now().Format(time.RFC3339)
	err = UpdateUserDaily(userID, guildID, 600, nowStr)
	if err != nil {
		t.Fatalf("UpdateUserDaily failed: %v", err)
	}

	u, _ = GetUser(userID, guildID)
	if u.Balance != 600 || u.LastDaily != nowStr {
		t.Errorf("Daily update failed. Expected 600/%s, got %d/%s", nowStr, u.Balance, u.LastDaily)
	}
}

func TestLeaderboard(t *testing.T) {
	guildID := "leaderboard_test"

	UpdateUserXP("user1", guildID, 100, 1)
	UpdateUserXP("user2", guildID, 300, 3)
	UpdateUserXP("user3", guildID, 200, 2)

	// Update without GetUser first in this test, must ensure they exist.
	// Actually UpdateUserXP doesn't insert if not exists. Let's call GetUser to create them.
	GetUser("user1", guildID)
	GetUser("user2", guildID)
	GetUser("user3", guildID)

	UpdateUserXP("user1", guildID, 100, 1)
	UpdateUserXP("user2", guildID, 300, 3)
	UpdateUserXP("user3", guildID, 200, 2)

	lb, err := GetLeaderboard(guildID, 10)
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
	guildID := "prefix_test"

	// Default should be !
	p, err := GetPrefix(guildID)
	if err != nil {
		t.Fatalf("GetPrefix failed: %v", err)
	}
	if p != "!" {
		t.Errorf("Expected default prefix '!', got '%s'", p)
	}

	// Set custom
	err = SetPrefix(guildID, "?")
	if err != nil {
		t.Fatalf("SetPrefix failed: %v", err)
	}

	p, _ = GetPrefix(guildID)
	if p != "?" {
		t.Errorf("Expected prefix '?', got '%s'", p)
	}

	// Update existing
	err = SetPrefix(guildID, "$")
	if err != nil {
		t.Fatalf("SetPrefix failed: %v", err)
	}

	p, _ = GetPrefix(guildID)
	if p != "$" {
		t.Errorf("Expected prefix '$', got '%s'", p)
	}
}

func TestReminders(t *testing.T) {
	id, err := CreateReminder("123", "chan", "guild", "hello", time.Now())
	if err != nil {
		t.Fatalf("CreateReminder failed: %v", err)
	}

	pending, err := GetPendingReminders()
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

	err = DeleteReminder(id)
	if err != nil {
		t.Fatalf("DeleteReminder failed: %v", err)
	}
}
