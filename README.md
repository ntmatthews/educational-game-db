# Educational Game Database

A comprehensive database system for managing student accounts in educational games, built with Go. Features both terminal-based CLI and web interface with PWA capabilities.

## Features

- **Account Management**: Create, read, update, delete student accounts
- **Terminal Interface**: Command-line interface for admin operations
- **Web Interface**: Modern PWA with HTML, CSS, and JavaScript
- **Progressive Web App**: Installable with offline support
- **Security**: Password hashing and basic authentication
- **Database**: SQLite for data persistence
- **Git Integration**: Version control ready
- **Cross-platform**: Works on macOS, Linux, and Windows
- **VS Code Integration**: Optimized for Visual Studio Code development

## Installation

1. Ensure you have Go 1.21+ installed
2. Clone this repository
3. Run `go mod tidy` to install dependencies
4. Build and run the application

```bash
# Build the application
go build -o educational-game-db cmd/cli/main.go

# Or run directly
go run cmd/cli/main.go
```

## Usage

### Terminal Interface
```bash
# Create a new student account
./educational-game-db create

# List all accounts
./educational-game-db list

# Get account by ID
./educational-game-db get 1

# Update account
./educational-game-db update 1

# Delete account
./educational-game-db delete 1

# Show statistics
./educational-game-db stats

# Start interactive mode
./educational-game-db interactive

# Start web server
./educational-game-db web --port 8081
```

### Web Interface

1. Start the web server:
   ```bash
   ./educational-game-db web --port 8081
   ```

2. Open your browser to:
   - Student Portal: http://localhost:8081
   - Admin Dashboard: http://localhost:8081/admin

3. **PWA Installation**: Click the "Install App" button to install as a Progressive Web App

### VS Code Development

This project is optimized for VS Code development with:
- Go extension for syntax highlighting and debugging
- Live Server extension for web development
- GitLens for Git integration
- Thunder Client for API testing

Use the VS Code task runner (Ctrl+Shift+P → "Tasks: Run Task") to:
- Build the application
- Start the web server
- Run tests

## API Endpoints

- `GET /api/accounts` - List all accounts
- `POST /api/accounts` - Create new account
- `GET /api/accounts/:id` - Get account by ID
- `PUT /api/accounts/:id` - Update account
- `DELETE /api/accounts/:id` - Delete account
- `GET /api/stats` - Get account statistics
- `POST /api/login` - Student login

## Progressive Web App Features

- **Offline Support**: Service worker caches resources for offline use
- **Installable**: Can be installed on mobile and desktop devices
- **Responsive Design**: Works on all screen sizes
- **App-like Experience**: Fullscreen mode with app icons

## Git Integration

The project is fully configured for Git with:
- `.gitignore` file excluding database files and build artifacts
- GitHub/GitLab ready with proper project structure
- VS Code Git integration with GitLens

### Remote Repository Setup

```bash
# Add GitHub remote
git remote add origin https://github.com/yourusername/educational-game-db.git

# Add GitLab remote  
git remote add gitlab https://gitlab.com/yourusername/educational-game-db.git

# Push to GitHub
git push -u origin main

# Push to GitLab
git push -u gitlab main
```

## Development

### Project Structure
```
.
├── cmd/cli/main.go              # CLI application entry point
├── internal/
│   ├── database/database.go     # Database operations
│   ├── handlers/handlers.go     # HTTP request handlers
│   ├── models/account.go        # Data models
│   └── server/server.go         # Web server setup
├── web/
│   ├── static/                  # Static web assets (CSS, JS, icons)
│   └── templates/               # HTML templates
├── .vscode/tasks.json           # VS Code tasks
├── go.mod                       # Go module definition
└── README.md                    # This file
```

### Database Schema

The SQLite database includes:
- **accounts** table with student information
- Indexed columns for performance
- Password hashing with bcrypt
- Automatic timestamps

### Security Features

- Password hashing using bcrypt
- Input validation and sanitization
- SQL injection prevention with prepared statements
- HTTPS ready (configure with TLS certificates)

## Testing

Test the API endpoints using:
- VS Code Thunder Client extension
- curl commands
- Web interface forms

Example API test:
```bash
# Create account
curl -X POST http://localhost:8081/api/accounts \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123","first_name":"John","last_name":"Doe","grade":5,"school":"Test School"}'

# Get all accounts
curl http://localhost:8081/api/accounts
```

## License

MIT License
