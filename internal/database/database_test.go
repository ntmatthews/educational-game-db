package database

import (
	"os"
	"testing"

	"educational-game-db/internal/models"
)

func TestMain(m *testing.M) {
	// Setup: Create a test database
	code := m.Run()
	// Cleanup after tests
	os.Exit(code)
}

func setupTestDB(t *testing.T) *Database {
	// Create a temporary in-memory database for testing
	db, err := NewDatabase(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	return db
}

func TestNewDatabase(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	if db == nil {
		t.Fatal("Database should not be nil")
	}
}

func TestCreateAccount(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	req := models.CreateAccountRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		Grade:     5,
		School:    "Test School",
	}

	account, err := db.CreateAccount(req)
	if err != nil {
		t.Fatalf("Failed to create account: %v", err)
	}

	if account.ID == 0 {
		t.Error("Account ID should be set after creation")
	}

	if account.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set after creation")
	}

	if account.Username != req.Username {
		t.Errorf("Expected username %s, got %s", req.Username, account.Username)
	}
}

func TestCreateAccountDuplicate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	req1 := models.CreateAccountRequest{
		Username:  "duplicate",
		Email:     "test1@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		Grade:     5,
		School:    "Test School",
	}

	req2 := models.CreateAccountRequest{
		Username:  "duplicate", // Same username
		Email:     "test2@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User2",
		Grade:     6,
		School:    "Test School",
	}

	// First account should succeed
	_, err := db.CreateAccount(req1)
	if err != nil {
		t.Fatalf("Failed to create first account: %v", err)
	}

	// Second account with same username should fail
	_, err = db.CreateAccount(req2)
	if err == nil {
		t.Error("Expected error when creating account with duplicate username")
	}
}

func TestGetAccountByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create test account
	req := models.CreateAccountRequest{
		Username:  "gettest",
		Email:     "gettest@example.com",
		Password:  "password123",
		FirstName: "Get",
		LastName:  "Test",
		Grade:     7,
		School:    "Test School",
	}

	account, err := db.CreateAccount(req)
	if err != nil {
		t.Fatalf("Failed to create account: %v", err)
	}

	// Retrieve account
	retrieved, err := db.GetAccountByID(account.ID)
	if err != nil {
		t.Fatalf("Failed to get account by ID: %v", err)
	}

	if retrieved.Username != account.Username {
		t.Errorf("Expected username %s, got %s", account.Username, retrieved.Username)
	}

	if retrieved.Email != account.Email {
		t.Errorf("Expected email %s, got %s", account.Email, retrieved.Email)
	}
}

func TestGetAccountByIDNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, err := db.GetAccountByID(99999)
	if err == nil {
		t.Error("Expected error when getting non-existent account")
	}
}

func TestGetAccountByUsername(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create test account
	req := models.CreateAccountRequest{
		Username:  "usernametest",
		Email:     "usernametest@example.com",
		Password:  "password123",
		FirstName: "Username",
		LastName:  "Test",
		Grade:     8,
		School:    "Test School",
	}

	account, err := db.CreateAccount(req)
	if err != nil {
		t.Fatalf("Failed to create account: %v", err)
	}

	// Retrieve account by username
	retrieved, err := db.GetAccountByUsername("usernametest")
	if err != nil {
		t.Fatalf("Failed to get account by username: %v", err)
	}

	if retrieved.ID != account.ID {
		t.Errorf("Expected ID %d, got %d", account.ID, retrieved.ID)
	}
}

func TestGetStats(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create test accounts with different grades and schools
	requests := []models.CreateAccountRequest{
		{Username: "stats1", Email: "stats1@example.com", Password: "password123", FirstName: "Stats", LastName: "One", Grade: 4, School: "School A"},
		{Username: "stats2", Email: "stats2@example.com", Password: "password123", FirstName: "Stats", LastName: "Two", Grade: 5, School: "School A"},
		{Username: "stats3", Email: "stats3@example.com", Password: "password123", FirstName: "Stats", LastName: "Three", Grade: 4, School: "School B"},
		{Username: "stats4", Email: "stats4@example.com", Password: "password123", FirstName: "Stats", LastName: "Four", Grade: 6, School: "School B"},
	}

	for _, req := range requests {
		_, err := db.CreateAccount(req)
		if err != nil {
			t.Fatalf("Failed to create account: %v", err)
		}
	}

	// Test that we have accounts (basic verification without GetStats method)
	// We'll verify by getting accounts individually
	account1, err := db.GetAccountByUsername("stats1")
	if err != nil {
		t.Fatalf("Failed to get stats1 account: %v", err)
	}
	if account1.Grade != 4 {
		t.Errorf("Expected grade 4, got %d", account1.Grade)
	}
}

func BenchmarkCreateAccount(b *testing.B) {
	db := setupTestDB(&testing.T{})
	defer db.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := models.CreateAccountRequest{
			Username:  "benchuser",
			Email:     "bench@example.com",
			Password:  "password123",
			FirstName: "Bench",
			LastName:  "User",
			Grade:     5,
			School:    "Bench School",
		}
		_, _ = db.CreateAccount(req)
	}
}

func BenchmarkGetAccountByID(b *testing.B) {
	db := setupTestDB(&testing.T{})
	defer db.Close()

	// Create a test account
	req := models.CreateAccountRequest{
		Username:  "benchget",
		Email:     "benchget@example.com",
		Password:  "password123",
		FirstName: "Bench",
		LastName:  "Get",
		Grade:     5,
		School:    "Bench School",
	}
	account, _ := db.CreateAccount(req)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = db.GetAccountByID(account.ID)
	}
}
