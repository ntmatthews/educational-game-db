package models

import (
	"time"
)

// Account represents a student account in the educational game
type Account struct {
	ID           int       `json:"id" db:"id"`
	Username     string    `json:"username" db:"username" binding:"required"`
	Email        string    `json:"email" db:"email" binding:"required,email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	FirstName    string    `json:"first_name" db:"first_name"`
	LastName     string    `json:"last_name" db:"last_name"`
	Grade        int       `json:"grade" db:"grade"`
	School       string    `json:"school" db:"school"`
	GameLevel    int       `json:"game_level" db:"game_level"`
	Experience   int       `json:"experience" db:"experience"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	IsActive     bool      `json:"is_active" db:"is_active"`
}

// CreateAccountRequest represents the request payload for creating an account
type CreateAccountRequest struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Grade     int    `json:"grade"`
	School    string `json:"school"`
}

// UpdateAccountRequest represents the request payload for updating an account
type UpdateAccountRequest struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Grade      int    `json:"grade"`
	School     string `json:"school"`
	GameLevel  int    `json:"game_level"`
	Experience int    `json:"experience"`
	IsActive   bool   `json:"is_active"`
}

// AccountStats represents aggregated statistics about accounts
type AccountStats struct {
	TotalAccounts    int     `json:"total_accounts"`
	ActiveAccounts   int     `json:"active_accounts"`
	AverageGameLevel float64 `json:"average_game_level"`
	TotalExperience  int     `json:"total_experience"`
}
