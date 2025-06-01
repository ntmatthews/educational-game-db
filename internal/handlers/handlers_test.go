package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"educational-game-db/internal/database"
	"educational-game-db/internal/models"

	"github.com/gin-gonic/gin"
)

func setupTestHandler() (*Handler, *database.Database) {
	// Create test database
	db, _ := database.NewDatabase(":memory:")
	
	// Create handler
	handler := NewHandler(db)
	
	return handler, db
}

func TestCreateAccountHandler(t *testing.T) {
	handler, db := setupTestHandler()
	defer db.Close()

	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	// Create request payload
	req := models.CreateAccountRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		Grade:     5,
		School:    "Test School",
	}

	jsonData, _ := json.Marshal(req)

	// Create HTTP request
	httpReq, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")

	// Create test response recorder
	w := httptest.NewRecorder()
	
	// Create gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = httpReq

	// Call handler
	handler.CreateAccount(c)

	// Check response
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	// Parse response
	var response models.Account
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.Username != req.Username {
		t.Errorf("Expected username %s, got %s", req.Username, response.Username)
	}
}

func TestCreateAccountHandlerInvalidData(t *testing.T) {
	handler, db := setupTestHandler()
	defer db.Close()

	gin.SetMode(gin.TestMode)

	// Create invalid request (missing required fields)
	req := models.CreateAccountRequest{
		Username: "testuser",
		// Missing email and password
	}

	jsonData, _ := json.Marshal(req)

	httpReq, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httpReq

	handler.CreateAccount(c)

	if w.Code == http.StatusCreated {
		t.Errorf("Expected error status, but got %d", w.Code)
	}
}

func TestGetAccountHandler(t *testing.T) {
	handler, db := setupTestHandler()
	defer db.Close()

	gin.SetMode(gin.TestMode)

	// Create test account first
	req := models.CreateAccountRequest{
		Username:  "gettest",
		Email:     "get@example.com",
		Password:  "password123",
		FirstName: "Get",
		LastName:  "Test",
		Grade:     7,
		School:    "Test School",
	}

	account, _ := db.CreateAccount(req)

	// Create HTTP request
	httpReq, _ := http.NewRequest("GET", "/accounts/1", nil)
	
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httpReq
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	handler.GetAccount(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response models.Account
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.Username != account.Username {
		t.Errorf("Expected username %s, got %s", account.Username, response.Username)
	}
}

func TestGetAccountHandlerNotFound(t *testing.T) {
	handler, db := setupTestHandler()
	defer db.Close()

	gin.SetMode(gin.TestMode)

	httpReq, _ := http.NewRequest("GET", "/accounts/99999", nil)
	
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httpReq
	c.Params = gin.Params{{Key: "id", Value: "99999"}}

	handler.GetAccount(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestGetAccountsHandler(t *testing.T) {
	handler, db := setupTestHandler()
	defer db.Close()

	gin.SetMode(gin.TestMode)

	// Create multiple test accounts
	accounts := []models.CreateAccountRequest{
		{Username: "list1", Email: "list1@example.com", Password: "password123", FirstName: "List", LastName: "One", Grade: 4, School: "Test School"},
		{Username: "list2", Email: "list2@example.com", Password: "password123", FirstName: "List", LastName: "Two", Grade: 5, School: "Test School"},
	}

	for _, req := range accounts {
		_, _ = db.CreateAccount(req)
	}

	httpReq, _ := http.NewRequest("GET", "/accounts", nil)
	
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httpReq

	handler.GetAccounts(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response []models.Account
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(response) != len(accounts) {
		t.Errorf("Expected %d accounts, got %d", len(accounts), len(response))
	}
}

func TestGetStatsHandler(t *testing.T) {
	handler, db := setupTestHandler()
	defer db.Close()

	gin.SetMode(gin.TestMode)

	// Create test accounts
	accounts := []models.CreateAccountRequest{
		{Username: "stats1", Email: "stats1@example.com", Password: "password123", FirstName: "Stats", LastName: "One", Grade: 4, School: "Test School"},
		{Username: "stats2", Email: "stats2@example.com", Password: "password123", FirstName: "Stats", LastName: "Two", Grade: 5, School: "Test School"},
	}

	for _, req := range accounts {
		_, _ = db.CreateAccount(req)
	}

	httpReq, _ := http.NewRequest("GET", "/stats", nil)
	
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httpReq

	handler.GetStats(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Basic validation that we got stats data
	if response["total_accounts"] == nil {
		t.Error("Expected total_accounts in response")
	}
}
