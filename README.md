# Go PostgreSQL Starter

A clean, production-ready Go web application starter template with PostgreSQL database integration, database migrations, and Docker support.

## 🚀 Features

- **Clean Architecture**: Well-organized project structure following Go best practices
- **PostgreSQL Integration**: Robust database connectivity with connection pooling
- **Database Migrations**: Automated schema management using Goose
- **Docker Support**: Easy development setup with Docker Compose
- **JSON API**: RESTful API with proper JSON response handling
- **Logging**: Structured logging throughout the application
- **Hot Reload**: Development server with Air for automatic reloading

## 📋 Prerequisites

Before running this project, make sure you have the following installed:

- [Go](https://golang.org/dl/) (version 1.24.0 or later)
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Air](https://github.com/air-verse/air) (for hot reload during development)

## 🛠️ Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/deepraj02/go-postgres-starter.git
   cd go-postgres-starter
   ```

2. **Install Go dependencies**
   ```bash
   go mod download
   ```

3. **Install Air for hot reload (optional but recommended for development)**
   ```bash
   go install github.com/air-verse/air@latest
   ```

## 🚀 Quick Start

### Using Make (Recommended)

1. **Start the application**
   ```bash
   make start
   ```
   This command will:
   - Start PostgreSQL database using Docker Compose
   - Run the Go application with hot reload using Air

2. **Stop the application**
   ```bash
   make stop
   ```

### Manual Setup

1. **Start PostgreSQL database**
   ```bash
   docker compose up -d
   ```

2. **Run the application**
   ```bash
   go run main.go
   ```

3. **Or run with custom port**
   ```bash
   go run main.go -port=3000
   ```

## 🏗️ Project Structure

```
.
├── main.go                 # Application entry point
├── go.mod                  # Go module dependencies
├── docker-compose.yml      # Docker services configuration
├── Makefile               # Build and development commands
├── internal/              # Private application code
│   ├── app/
│   │   └── app.go         # Application setup and configuration
│   ├── routes/
│   │   └── routes.go      # HTTP route definitions
│   ├── store/
│   │   └── database.go    # Database connection and operations
│   └── utils/
│       └── utils.go       # Utility functions (JSON responses, etc.)
└── migrations/            # Database schema migrations
    ├── fs.go              # Embedded filesystem for migrations
    └── 00001_users.sql    # Initial user table migration
```

## 📁 Code Explanation

### Main Application (`main.go`)
The entry point of the application that:
- Parses command-line flags for port configuration
- Initializes the application with database connection
- Sets up HTTP server with proper timeouts
- Starts the web server

### Application Setup (`internal/app/app.go`)
Contains the core application structure:
- **Application struct**: Holds logger and database connection
- **NewApplication()**: Initializes database connection and runs migrations
- **HealthCheck()**: Endpoint to verify database connectivity

### Database Layer (`internal/store/database.go`)
Manages database operations:
- **Open()**: Establishes PostgreSQL connection
- **MigrateFS()**: Applies database migrations using embedded files
- Uses `pgx` driver for efficient PostgreSQL connectivity

### Routes (`internal/routes/routes.go`)
Defines HTTP endpoints using Chi router:
- Currently includes health check endpoint
- Easily extensible for additional API endpoints

### Utilities (`internal/utils/utils.go`)
Common helper functions:
- **WriteJson()**: Standardized JSON response formatting
- **Envelope**: Type for consistent API response structure

### Migrations (`migrations/`)
Database schema management:
- **fs.go**: Embeds SQL files into the binary
- **00001_users.sql**: Creates initial users table with proper indexing

## 🔧 API Endpoints

| Method | Endpoint  | Description                                 |
| ------ | --------- | ------------------------------------------- |
| GET    | `/health` | Health check - verifies database connection |

### Health Check Response
```json
{
  "status": "Healthy"
}
```

## 🗄️ Database Schema

### Users Table
```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    bio TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

## 🐳 Docker Services

The `docker-compose.yml` sets up:

1. **PostgreSQL Database**
   - Container: `postgres-project`
   - Port: `5432`
   - Credentials: `postgres/postgres`
   - Persistent data storage

2. **Adminer** (Database Admin UI)
   - Container: `adminer`
   - Port: `9090`
   - Access: http://localhost:9090

## 🔧 Configuration

### Environment Variables
The application uses the following default configuration:
- **Database Host**: `localhost:5432`
- **Database Name**: `postgres`
- **Database User**: `postgres`
- **Database Password**: `postgres`
- **Application Port**: `8080` (configurable via `-port` flag)

### Server Configuration
- **Idle Timeout**: 1 minute
- **Read Timeout**: 10 seconds
- **Write Timeout**: 10 seconds

## 🧪 Testing the Application

1. **Check if the server is running**
   ```bash
   curl http://localhost:8080/health
   ```

2. **Expected response**
   ```json
   {
     "status": "Healthy"
   }
   ```

3. **Access database admin** (optional)
   - Visit: http://localhost:9090
   - Server: `db`
   - Username: `postgres`
   - Password: `postgres`
   - Database: `postgres`

## 🚀 Extending the Application

### Adding New Endpoints
1. Create handler functions in `internal/app/app.go`
2. Add routes in `internal/routes/routes.go`
3. Use the `utils.WriteJson()` function for consistent responses

### Adding Database Models
1. Create new migration files in `migrations/`
2. Add corresponding Go structs
3. Implement CRUD operations in the store package

### Example: Adding a new endpoint
```go
// In internal/app/app.go
func (app *Application) GetUsers(w http.ResponseWriter, r *http.Request) {
    // Implementation here
    utils.WriteJson(w, http.StatusOK, utils.Envelope{"users": users})
}

// In internal/routes/routes.go
func SetupRoutes(app *app.Application) *chi.Mux {
    r := chi.NewRouter()
    r.Get("/health", app.HealthCheck)
    r.Get("/users", app.GetUsers) // New endpoint
    return r
}
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is open source and available under the [MIT License](LICENSE).

## 🆘 Troubleshooting

### Common Issues

1. **Database connection failed**
   - Ensure Docker is running
   - Check if PostgreSQL container is up: `docker compose ps`
   - Verify database credentials in connection string

2. **Port already in use**
   - Change the port using: `go run main.go -port=3001`
   - Or stop the service using the conflicting port

3. **Migration errors**
   - Check migration file syntax
   - Ensure database is accessible
   - Review migration logs in application output

### Getting Help

If you encounter any issues:
1. Check the application logs
2. Verify Docker containers are running
3. Ensure all dependencies are installed
4. Review the database connection settings

---

**Happy coding! 🎉**