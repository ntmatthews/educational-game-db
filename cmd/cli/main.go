package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"educational-game-db/internal/database"
	"educational-game-db/internal/models"
	"educational-game-db/internal/server"
	"github.com/spf13/cobra"
)

var (
	dbPath string
	port   string
	db     *database.Database
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "educational-game-db",
		Short: "Educational Game Database CLI",
		Long:  "A command-line interface for managing student accounts in the educational game database.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			var err error
			db, err = database.NewDatabase(dbPath)
			if err != nil {
				log.Fatalf("Failed to connect to database: %v", err)
			}
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if db != nil {
				db.Close()
			}
		},
	}

	rootCmd.PersistentFlags().StringVar(&dbPath, "db", "accounts.db", "Database file path")

	// Create account command
	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new student account",
		Run:   createAccount,
	}

	// List accounts command
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all student accounts",
		Run:   listAccounts,
	}

	// Get account command
	var getCmd = &cobra.Command{
		Use:   "get [id]",
		Short: "Get account by ID",
		Args:  cobra.ExactArgs(1),
		Run:   getAccount,
	}

	// Update account command
	var updateCmd = &cobra.Command{
		Use:   "update [id]",
		Short: "Update account by ID",
		Args:  cobra.ExactArgs(1),
		Run:   updateAccount,
	}

	// Delete account command
	var deleteCmd = &cobra.Command{
		Use:   "delete [id]",
		Short: "Delete account by ID",
		Args:  cobra.ExactArgs(1),
		Run:   deleteAccount,
	}

	// Stats command
	var statsCmd = &cobra.Command{
		Use:   "stats",
		Short: "Show account statistics",
		Run:   showStats,
	}

	// Web server command
	var webCmd = &cobra.Command{
		Use:   "web",
		Short: "Start the web server",
		Run:   startWebServer,
	}
	webCmd.Flags().StringVar(&port, "port", "8080", "Web server port")

	// Interactive mode command
	var interactiveCmd = &cobra.Command{
		Use:   "interactive",
		Short: "Start interactive mode",
		Run:   startInteractive,
	}

	rootCmd.AddCommand(createCmd, listCmd, getCmd, updateCmd, deleteCmd, statsCmd, webCmd, interactiveCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func createAccount(cmd *cobra.Command, args []string) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Username: ")
	scanner.Scan()
	username := scanner.Text()

	fmt.Print("Email: ")
	scanner.Scan()
	email := scanner.Text()

	fmt.Print("Password: ")
	scanner.Scan()
	password := scanner.Text()

	fmt.Print("First Name: ")
	scanner.Scan()
	firstName := scanner.Text()

	fmt.Print("Last Name: ")
	scanner.Scan()
	lastName := scanner.Text()

	fmt.Print("Grade (0 for none): ")
	scanner.Scan()
	grade, _ := strconv.Atoi(scanner.Text())

	fmt.Print("School: ")
	scanner.Scan()
	school := scanner.Text()

	req := models.CreateAccountRequest{
		Username:  username,
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
		Grade:     grade,
		School:    school,
	}

	account, err := db.CreateAccount(req)
	if err != nil {
		fmt.Printf("Error creating account: %v\n", err)
		return
	}

	fmt.Printf("Account created successfully!\n")
	fmt.Printf("ID: %d\n", account.ID)
	fmt.Printf("Username: %s\n", account.Username)
	fmt.Printf("Email: %s\n", account.Email)
}

func listAccounts(cmd *cobra.Command, args []string) {
	accounts, err := db.GetAllAccounts()
	if err != nil {
		fmt.Printf("Error listing accounts: %v\n", err)
		return
	}

	if len(accounts) == 0 {
		fmt.Println("No accounts found.")
		return
	}

	fmt.Printf("%-5s %-15s %-25s %-15s %-15s %-10s %-5s %-10s\n", 
		"ID", "Username", "Email", "First Name", "Last Name", "Grade", "Level", "XP")
	fmt.Println(strings.Repeat("-", 100))

	for _, account := range accounts {
		status := "Active"
		if !account.IsActive {
			status = "Inactive"
		}
		
		fmt.Printf("%-5d %-15s %-25s %-15s %-15s %-10d %-5d %-10d [%s]\n",
			account.ID, account.Username, account.Email, account.FirstName,
			account.LastName, account.Grade, account.GameLevel, account.Experience, status)
	}
}

func getAccount(cmd *cobra.Command, args []string) {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("Invalid account ID: %v\n", err)
		return
	}

	account, err := db.GetAccountByID(id)
	if err != nil {
		fmt.Printf("Error getting account: %v\n", err)
		return
	}

	fmt.Printf("Account Details:\n")
	fmt.Printf("ID: %d\n", account.ID)
	fmt.Printf("Username: %s\n", account.Username)
	fmt.Printf("Email: %s\n", account.Email)
	fmt.Printf("Name: %s %s\n", account.FirstName, account.LastName)
	fmt.Printf("Grade: %d\n", account.Grade)
	fmt.Printf("School: %s\n", account.School)
	fmt.Printf("Game Level: %d\n", account.GameLevel)
	fmt.Printf("Experience: %d\n", account.Experience)
	fmt.Printf("Active: %t\n", account.IsActive)
	fmt.Printf("Created: %s\n", account.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated: %s\n", account.UpdatedAt.Format("2006-01-02 15:04:05"))
}

func updateAccount(cmd *cobra.Command, args []string) {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("Invalid account ID: %v\n", err)
		return
	}

	// Get current account
	current, err := db.GetAccountByID(id)
	if err != nil {
		fmt.Printf("Error getting account: %v\n", err)
		return
	}

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("First Name [%s]: ", current.FirstName)
	scanner.Scan()
	firstName := scanner.Text()
	if firstName == "" {
		firstName = current.FirstName
	}

	fmt.Printf("Last Name [%s]: ", current.LastName)
	scanner.Scan()
	lastName := scanner.Text()
	if lastName == "" {
		lastName = current.LastName
	}

	fmt.Printf("Grade [%d]: ", current.Grade)
	scanner.Scan()
	gradeStr := scanner.Text()
	grade := current.Grade
	if gradeStr != "" {
		grade, _ = strconv.Atoi(gradeStr)
	}

	fmt.Printf("School [%s]: ", current.School)
	scanner.Scan()
	school := scanner.Text()
	if school == "" {
		school = current.School
	}

	fmt.Printf("Game Level [%d]: ", current.GameLevel)
	scanner.Scan()
	levelStr := scanner.Text()
	gameLevel := current.GameLevel
	if levelStr != "" {
		gameLevel, _ = strconv.Atoi(levelStr)
	}

	fmt.Printf("Experience [%d]: ", current.Experience)
	scanner.Scan()
	expStr := scanner.Text()
	experience := current.Experience
	if expStr != "" {
		experience, _ = strconv.Atoi(expStr)
	}

	req := models.UpdateAccountRequest{
		FirstName:  firstName,
		LastName:   lastName,
		Grade:      grade,
		School:     school,
		GameLevel:  gameLevel,
		Experience: experience,
		IsActive:   current.IsActive,
	}

	account, err := db.UpdateAccount(id, req)
	if err != nil {
		fmt.Printf("Error updating account: %v\n", err)
		return
	}

	fmt.Printf("Account updated successfully!\n")
	fmt.Printf("Updated: %s\n", account.UpdatedAt.Format("2006-01-02 15:04:05"))
}

func deleteAccount(cmd *cobra.Command, args []string) {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("Invalid account ID: %v\n", err)
		return
	}

	// Get account details first
	account, err := db.GetAccountByID(id)
	if err != nil {
		fmt.Printf("Error getting account: %v\n", err)
		return
	}

	fmt.Printf("Are you sure you want to delete account for %s (%s)? (y/N): ", 
		account.Username, account.Email)
	
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	confirmation := strings.ToLower(scanner.Text())

	if confirmation != "y" && confirmation != "yes" {
		fmt.Println("Account deletion cancelled.")
		return
	}

	err = db.DeleteAccount(id)
	if err != nil {
		fmt.Printf("Error deleting account: %v\n", err)
		return
	}

	fmt.Printf("Account for %s deleted successfully.\n", account.Username)
}

func showStats(cmd *cobra.Command, args []string) {
	stats, err := db.GetAccountStats()
	if err != nil {
		fmt.Printf("Error getting stats: %v\n", err)
		return
	}

	fmt.Printf("Account Statistics:\n")
	fmt.Printf("Total Accounts: %d\n", stats.TotalAccounts)
	fmt.Printf("Active Accounts: %d\n", stats.ActiveAccounts)
	fmt.Printf("Average Game Level: %.2f\n", stats.AverageGameLevel)
	fmt.Printf("Total Experience: %d\n", stats.TotalExperience)
}

func startWebServer(cmd *cobra.Command, args []string) {
	srv := server.NewServer(db, port)
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start web server: %v", err)
	}
}

func startInteractive(cmd *cobra.Command, args []string) {
	scanner := bufio.NewScanner(os.Stdin)
	
	fmt.Println("Educational Game Database - Interactive Mode")
	fmt.Println("Type 'help' for available commands, 'exit' to quit")
	
	for {
		fmt.Print("> ")
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())
		
		if input == "exit" || input == "quit" {
			break
		}
		
		switch input {
		case "help":
			fmt.Println("Available commands:")
			fmt.Println("  list    - List all accounts")
			fmt.Println("  create  - Create new account")
			fmt.Println("  stats   - Show statistics")
			fmt.Println("  web     - Start web server")
			fmt.Println("  help    - Show this help")
			fmt.Println("  exit    - Exit interactive mode")
		case "list":
			listAccounts(nil, nil)
		case "create":
			createAccount(nil, nil)
		case "stats":
			showStats(nil, nil)
		case "web":
			fmt.Printf("Starting web server on port %s...\n", port)
			startWebServer(nil, nil)
		default:
			fmt.Printf("Unknown command: %s. Type 'help' for available commands.\n", input)
		}
	}
	
	fmt.Println("Goodbye!")
}
