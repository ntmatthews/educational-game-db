# Educational Game Database

A comprehensive database system for managing student accounts in educational games, built with Go. Features both terminal-based CLI and web interface with PWA capabilities.

## Features

- **Account Management**: Create, read, update, delete student accounts
- **Terminal Interface**: Command-line interface for admin operations
- **Web Interface**: Modern PWA with HTML, CSS, and JavaScript
- **Security**: Password hashing and basic authentication
- **Database**: SQLite for data persistence
- **Git Integration**: Version control ready
- **Cross-platform**: Works on macOS, Linux, and Windows

## Installation

1. Ensure you have Go 1.21+ installed
2. Clone this repository
3. Run `go mod tidy` to install dependencies
4. Build and run the application

## Usage

### Terminal Interface
```bash
# Start the CLI
go run cmd/cli/main.go

# Available commands:
./educational-game-db create-account --username john --email john@example.com
./educational-game-db list-accounts
./educational-game-db web-server --port 8080
```

### Web Interface
Start the web server and visit http://localhost:8080

## API Endpoints

- `GET /api/accounts` - List all accounts
- `POST /api/accounts` - Create new account
- `GET /api/accounts/:id` - Get account by ID
- `PUT /api/accounts/:id` - Update account
- `DELETE /api/accounts/:id` - Delete account

## License

MIT License
