package main

import (
	"forum/internal/db"
	"forum/internal/handlers"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func TestCheckUsername(t *testing.T) {
	Testinit()

	defer db.DB.Close()
	tests := []struct {
		input    string
		expected bool
	}{
		{"valid123", true},
		{"short", true},
		{"four", false},
		{"invalid!", false},
	}

	for _, test := range tests {
		result := handlers.CheckUsername(test.input)
		if result != test.expected {
			t.Errorf("checkUsername(%v) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func TestCheckPassword(t *testing.T) {
	Testinit()
	defer db.DB.Close()
	tests := []struct {
		input    string
		expected bool
	}{
		{"valid123", true},
		{"short", true},
		{"four", false},
		{"validNow!", true},
	}

	for _, test := range tests {
		result := handlers.CheckPassword(test.input)
		if result != test.expected {
			t.Errorf("checkUsername(%v) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func TestNameOrEmailExists(t *testing.T) {
	Testinit()
	defer db.DB.Close()

	// Create users table
	_, err := db.DB.Exec(`CREATE TABLE IF NOT EXISTS users (id TEXT PRIMARY KEY, email TEXT UNIQUE NOT NULL, username TEXT UNIQUE NOT NULL, password TEXT NOT NULL);`)
	if err != nil {
		t.Fatalf("failed to create users table: %v", err)
	}

	// Insert a test user
	_, err = db.DB.Exec(`INSERT INTO users (id, email, username, password) VALUES (?, ?, ?, ?);`, "test-id", "testuser@example.com", "testuser", "password")
	if err != nil {
		t.Fatalf("failed to insert test user: %v", err)
	}

	result := handlers.NameOremailExists("testuser")
	if !result {
		t.Errorf("nameOremailExists returned false; want true")
	}
}

func TestCreateSession(t *testing.T) {
	sessionToken, err := handlers.CreateSession()
	if err != nil {
		t.Fatalf("createSession returned error: %v", err)
	}

	if _, err := uuid.FromString(sessionToken); err != nil {
		t.Errorf("createSession returned invalid UUID: %v", sessionToken)
	}
}

func TestSaveSession(t *testing.T) {
	// Use in-memory SQLite database
	Testinit()
	defer db.DB.Close()

	// Create sessions table
	_, err := db.DB.Exec(`CREATE TABLE IF NOT EXISTS sessions (id INTEGER PRIMARY KEY AUTOINCREMENT, user_id TEXT NOT NULL, username TEXT, session_token TEXT UNIQUE NOT NULL, expires_at DATETIME NOT NULL);`)
	if err != nil {
		t.Fatalf("failed to create sessions table: %v", err)
	}

	err = handlers.SaveSession("userID", "username", "sessionToken", time.Now())
	if err != nil {
		t.Errorf("saveSession returned error: %v", err)
	}
}
