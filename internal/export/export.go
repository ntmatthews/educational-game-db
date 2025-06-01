package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"educational-game-db/internal/database"
	"educational-game-db/internal/models"
)

// ExportService handles data export/import operations
type ExportService struct {
	db *database.Database
}

// NewExportService creates a new export service
func NewExportService(db *database.Database) *ExportService {
	return &ExportService{db: db}
}

// ExportToCSV exports all accounts to a CSV file
func (e *ExportService) ExportToCSV(filename string) error {
	accounts, err := e.db.GetAllAccounts()
	if err != nil {
		return fmt.Errorf("failed to get accounts: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"ID", "Username", "Email", "FirstName", "LastName", 
		"Grade", "School", "GameLevel", "Experience", 
		"CreatedAt", "UpdatedAt", "IsActive",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data
	for _, account := range accounts {
		record := []string{
			strconv.Itoa(account.ID),
			account.Username,
			account.Email,
			account.FirstName,
			account.LastName,
			strconv.Itoa(account.Grade),
			account.School,
			strconv.Itoa(account.GameLevel),
			strconv.Itoa(account.Experience),
			account.CreatedAt.Format(time.RFC3339),
			account.UpdatedAt.Format(time.RFC3339),
			strconv.FormatBool(account.IsActive),
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	return nil
}

// ExportToJSON exports all accounts to a JSON file
func (e *ExportService) ExportToJSON(filename string) error {
	accounts, err := e.db.GetAllAccounts()
	if err != nil {
		return fmt.Errorf("failed to get accounts: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create JSON file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	
	exportData := struct {
		ExportedAt time.Time         `json:"exported_at"`
		Count      int               `json:"count"`
		Accounts   []models.Account  `json:"accounts"`
	}{
		ExportedAt: time.Now(),
		Count:      len(accounts),
		Accounts:   accounts,
	}

	if err := encoder.Encode(exportData); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

// ImportFromCSV imports accounts from a CSV file
func (e *ExportService) ImportFromCSV(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) < 2 {
		return fmt.Errorf("CSV file must contain at least a header and one record")
	}

	// Skip header row
	for i := 1; i < len(records); i++ {
		record := records[i]
		if len(record) < 12 {
			continue // Skip incomplete records
		}

		// Parse record
		grade, _ := strconv.Atoi(record[5])
		gameLevel, _ := strconv.Atoi(record[7])
		experience, _ := strconv.Atoi(record[8])
		isActive, _ := strconv.ParseBool(record[11])

		req := models.CreateAccountRequest{
			Username:  record[1],
			Email:     record[2],
			Password:  "imported123", // Default password for imported accounts
			FirstName: record[3],
			LastName:  record[4],
			Grade:     grade,
			School:    record[6],
		}

		// Create account
		account, err := e.db.CreateAccount(req)
		if err != nil {
			// Log error but continue with other records
			fmt.Printf("Failed to import account %s: %v\n", req.Username, err)
			continue
		}

		// Update additional fields that aren't in CreateAccountRequest
		updateReq := models.UpdateAccountRequest{
			FirstName:  account.FirstName,
			LastName:   account.LastName,
			Grade:      account.Grade,
			School:     account.School,
			GameLevel:  gameLevel,
			Experience: experience,
			IsActive:   isActive,
		}

		_, err = e.db.UpdateAccount(account.ID, updateReq)
		if err != nil {
			fmt.Printf("Failed to update imported account %s: %v\n", req.Username, err)
		}
	}

	return nil
}

// ImportFromJSON imports accounts from a JSON file
func (e *ExportService) ImportFromJSON(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open JSON file: %w", err)
	}
	defer file.Close()

	var importData struct {
		Accounts []struct {
			Username   string `json:"username"`
			Email      string `json:"email"`
			FirstName  string `json:"first_name"`
			LastName   string `json:"last_name"`
			Grade      int    `json:"grade"`
			School     string `json:"school"`
			GameLevel  int    `json:"game_level"`
			Experience int    `json:"experience"`
			IsActive   bool   `json:"is_active"`
		} `json:"accounts"`
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&importData); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	for _, accountData := range importData.Accounts {
		req := models.CreateAccountRequest{
			Username:  accountData.Username,
			Email:     accountData.Email,
			Password:  "imported123", // Default password for imported accounts
			FirstName: accountData.FirstName,
			LastName:  accountData.LastName,
			Grade:     accountData.Grade,
			School:    accountData.School,
		}

		// Create account
		account, err := e.db.CreateAccount(req)
		if err != nil {
			fmt.Printf("Failed to import account %s: %v\n", req.Username, err)
			continue
		}

		// Update additional fields
		updateReq := models.UpdateAccountRequest{
			FirstName:  account.FirstName,
			LastName:   account.LastName,
			Grade:      account.Grade,
			School:     account.School,
			GameLevel:  accountData.GameLevel,
			Experience: accountData.Experience,
			IsActive:   accountData.IsActive,
		}

		_, err = e.db.UpdateAccount(account.ID, updateReq)
		if err != nil {
			fmt.Printf("Failed to update imported account %s: %v\n", req.Username, err)
		}
	}

	return nil
}

// ExportStats exports statistical data to JSON
func (e *ExportService) ExportStats(filename string) error {
	stats, err := e.db.GetAccountStats()
	if err != nil {
		return fmt.Errorf("failed to get stats: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create stats file: %w", err)
	}
	defer file.Close()

	exportData := struct {
		ExportedAt time.Time              `json:"exported_at"`
		Stats      models.AccountStats    `json:"stats"`
	}{
		ExportedAt: time.Now(),
		Stats:      *stats,
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	
	if err := encoder.Encode(exportData); err != nil {
		return fmt.Errorf("failed to encode stats JSON: %w", err)
	}

	return nil
}
