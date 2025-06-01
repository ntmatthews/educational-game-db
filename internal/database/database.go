package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"educational-game-db/internal/models"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	database := &Database{db: db}
	if err := database.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return database, nil
}

func (d *Database) createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS accounts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		first_name TEXT,
		last_name TEXT,
		grade INTEGER DEFAULT 0,
		school TEXT,
		game_level INTEGER DEFAULT 1,
		experience INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		is_active BOOLEAN DEFAULT TRUE
	);

	CREATE INDEX IF NOT EXISTS idx_username ON accounts(username);
	CREATE INDEX IF NOT EXISTS idx_email ON accounts(email);
	CREATE INDEX IF NOT EXISTS idx_is_active ON accounts(is_active);
	`

	_, err := d.db.Exec(query)
	return err
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) CreateAccount(req models.CreateAccountRequest) (*models.Account, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	query := `
	INSERT INTO accounts (username, email, password_hash, first_name, last_name, grade, school, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := d.db.Exec(query, req.Username, req.Email, string(hashedPassword), 
		req.FirstName, req.LastName, req.Grade, req.School, now, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get account ID: %w", err)
	}

	return d.GetAccountByID(int(id))
}

func (d *Database) GetAccountByID(id int) (*models.Account, error) {
	query := `
	SELECT id, username, email, password_hash, first_name, last_name, grade, school,
		   game_level, experience, created_at, updated_at, is_active
	FROM accounts WHERE id = ?
	`

	var account models.Account
	err := d.db.QueryRow(query, id).Scan(
		&account.ID, &account.Username, &account.Email, &account.PasswordHash,
		&account.FirstName, &account.LastName, &account.Grade, &account.School,
		&account.GameLevel, &account.Experience, &account.CreatedAt, &account.UpdatedAt,
		&account.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account not found")
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return &account, nil
}

func (d *Database) GetAccountByUsername(username string) (*models.Account, error) {
	query := `
	SELECT id, username, email, password_hash, first_name, last_name, grade, school,
		   game_level, experience, created_at, updated_at, is_active
	FROM accounts WHERE username = ?
	`

	var account models.Account
	err := d.db.QueryRow(query, username).Scan(
		&account.ID, &account.Username, &account.Email, &account.PasswordHash,
		&account.FirstName, &account.LastName, &account.Grade, &account.School,
		&account.GameLevel, &account.Experience, &account.CreatedAt, &account.UpdatedAt,
		&account.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account not found")
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return &account, nil
}

func (d *Database) GetAllAccounts() ([]models.Account, error) {
	query := `
	SELECT id, username, email, password_hash, first_name, last_name, grade, school,
		   game_level, experience, created_at, updated_at, is_active
	FROM accounts ORDER BY created_at DESC
	`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts: %w", err)
	}
	defer rows.Close()

	var accounts []models.Account
	for rows.Next() {
		var account models.Account
		err := rows.Scan(
			&account.ID, &account.Username, &account.Email, &account.PasswordHash,
			&account.FirstName, &account.LastName, &account.Grade, &account.School,
			&account.GameLevel, &account.Experience, &account.CreatedAt, &account.UpdatedAt,
			&account.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan account: %w", err)
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (d *Database) UpdateAccount(id int, req models.UpdateAccountRequest) (*models.Account, error) {
	query := `
	UPDATE accounts 
	SET first_name = ?, last_name = ?, grade = ?, school = ?, 
		game_level = ?, experience = ?, is_active = ?, updated_at = ?
	WHERE id = ?
	`

	now := time.Now()
	_, err := d.db.Exec(query, req.FirstName, req.LastName, req.Grade, req.School,
		req.GameLevel, req.Experience, req.IsActive, now, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update account: %w", err)
	}

	return d.GetAccountByID(id)
}

func (d *Database) DeleteAccount(id int) error {
	query := `DELETE FROM accounts WHERE id = ?`
	result, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("account not found")
	}

	return nil
}

func (d *Database) GetAccountStats() (*models.AccountStats, error) {
	query := `
	SELECT 
		COUNT(*) as total_accounts,
		SUM(CASE WHEN is_active = 1 THEN 1 ELSE 0 END) as active_accounts,
		AVG(game_level) as average_game_level,
		SUM(experience) as total_experience
	FROM accounts
	`

	var stats models.AccountStats
	err := d.db.QueryRow(query).Scan(
		&stats.TotalAccounts,
		&stats.ActiveAccounts,
		&stats.AverageGameLevel,
		&stats.TotalExperience,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get account stats: %w", err)
	}

	return &stats, nil
}

func (d *Database) VerifyPassword(username, password string) bool {
	account, err := d.GetAccountByUsername(username)
	if err != nil {
		log.Printf("Failed to get account for password verification: %v", err)
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(password))
	return err == nil
}
