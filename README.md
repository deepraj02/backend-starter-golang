# Go PostgreSQL Authentication Service

A production-ready Go web application with complete authentication system, featuring user registration, login, forgot password functionality with email verification, Redis caching, and PostgreSQL database integration.

## 🚀 Features

### Core Features
- **Clean Architecture**: Well-organized project structure following Go best practices
- **PostgreSQL Integration**: Robust database connectivity with connection pooling
- **Redis Caching**: Fast in-memory caching for reset codes and session management
- **Database Migrations**: Automated schema management using Goose
- **Docker Support**: Easy development setup with Docker Compose
- **JSON API**: RESTful API with proper JSON response handling
- **Logging**: Structured logging throughout the application
- **Hot Reload**: Development server with Air for automatic reloading

### Authentication Features
- **User Registration**: Secure user account creation with password hashing
- **User Login**: JWT-based authentication with access and refresh tokens
- **User Profile**: Protected endpoint to retrieve user information
- **Forgot Password**: Email-based password reset with 6-digit verification codes
- **Password Reset**: Secure password update with code verification
- **Email Service**: Integrated Resend email service for transactional emails
- **Rate Limiting**: Redis-based reset code expiration (15 minutes)

### Security Features
- **Password Hashing**: bcrypt-based secure password storage
- **JWT Tokens**: Access and refresh token implementation
- **Input Validation**: Comprehensive request validation
- **Protected Routes**: Middleware-based route protection
- **Code Expiration**: Automatic cleanup of expired reset codes

## 📋 Prerequisites

Before running this project, make sure you have the following installed:

- [Go](https://golang.org/dl/) (version 1.24.0 or later)
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Air](https://github.com/air-verse/air) (for hot reload during development)
- [Resend Account](https://resend.com/) (for email service)

## 🛠️ Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/deepraj02/go-postgres-starter.git
   cd backend-auth
   ```

2. **Install Go dependencies**
   ```bash
   go mod download
   ```

3. **Create environment configuration**
   ```bash
   cp .env.example .env
   ```
   
4. **Configure environment variables in `.env`**
   ```bash
   # Database Configuration
   DB_HOST=localhost
   DB_PORT=5432
   DB_NAME=postgres
   DB_USER=postgres
   DB_PASSWORD=postgres
   
   # Redis Configuration
   REDIS_HOST=localhost
   REDIS_PORT=6379
   
   # Email Service
   RESEND_API_KEY=your_resend_api_key_here
   EMAIL_FROM=noreply@yourdomain.com
   
   # JWT Configuration
   JWT_SECRET=your_jwt_secret_here
   ```

5. **Install Air for hot reload (optional but recommended)**
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
   - Start PostgreSQL and Redis using Docker Compose
   - Load environment variables from `.env`
   - Run the Go application with hot reload using Air

2. **Stop the application**
   ```bash
   make stop
   ```

### Manual Setup

1. **Start services**
   ```bash
   docker compose up -d
   ```

2. **Export environment variables**
   ```bash
   export $(cat .env | xargs)
   ```

3. **Run the application**
   ```bash
   go run main.go
   ```

## 🏗️ Project Structure

```
.
├── main.go                     # Application entry point
├── go.mod                      # Go module dependencies
├── docker-compose.yml          # Docker services configuration
├── Makefile                    # Build and development commands
├── .env.example                # Environment variables template
├── static/                     # Frontend HTML pages
│   ├── index.html             # Landing page
│   ├── register.html          # User registration form
│   ├── login.html             # User login form
│   ├── dashboard.html         # User dashboard (protected)
│   ├── forgot-password.html   # Forgot password form
│   └── reset-password.html    # Password reset form
├── internal/                   # Private application code
│   ├── app/
│   │   └── app.go             # Application setup and configuration
│   ├── api/
│   │   └── auth_handler.go    # Authentication HTTP handlers
│   ├── middleware/
│   │   └── auth_middleware.go # JWT authentication middleware
│   ├── routes/
│   │   └── routes.go          # HTTP route definitions
│   ├── store/
│   │   ├── database.go        # Database connection
│   │   ├── auth_store.go      # User data operations
│   │   ├── cache_store.go     # Redis operations
│   │   └── email_store.go     # Email service implementation
│   └── utils/
│       ├── auth/
│       │   └── auth.go        # JWT and password utilities
│       ├── json/
│       │   └── json.go        # JSON response utilities
│       └── logger/
│           └── logger.go      # Structured logging
└── migrations/                 # Database schema migrations
    ├── fs.go                  # Embedded filesystem for migrations
    └── 00001_users.sql        # User table creation
```

## 📁 API Endpoints

### Public Endpoints

| Method | Endpoint             | Description                    |
|--------|---------------------|--------------------------------|
| GET    | `/health`           | Health check                   |
| POST   | `/auth/register`    | User registration              |
| POST   | `/auth/login`       | User login                     |
| POST   | `/auth/forgot-password` | Request password reset     |
| POST   | `/auth/reset-password`  | Reset password with code   |

### Protected Endpoints

| Method | Endpoint         | Description              |
|--------|-----------------|--------------------------|
| GET    | `/auth/profile` | Get user profile         |

### Frontend Pages

| Route                    | Description                    |
|--------------------------|--------------------------------|
| `/`                      | Landing page                   |
| `/register.html`         | User registration form         |
| `/login.html`            | User login form                |
| `/dashboard.html`        | User dashboard (requires auth) |
| `/forgot-password.html`  | Forgot password form           |
| `/reset-password.html`   | Password reset form            |

## 🔧 API Usage Examples

### User Registration
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "secure123",
    "bio": "Software developer"
  }'
```

**Response:**
```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "token_type": "Bearer",
    "expires_in": 86400,
    "user": {
      "id": 1,
      "username": "johndoe",
      "email": "john@example.com",
      "bio": "Software developer",
      "created_at": "2025-09-27T10:00:00Z",
      "updated_at": "2025-09-27T10:00:00Z"
    }
  }
}
```

### User Login
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "secure123"
  }'
```

### Forgot Password
```bash
curl -X POST http://localhost:8080/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com"
  }'
```

**Response:**
```json
{
  "message": "Reset code sent to your email"
}
```

### Reset Password
```bash
curl -X POST http://localhost:8080/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "code": "123456",
    "new_password": "newsecure123"
  }'
```

### Get User Profile (Protected)
```bash
curl -X GET http://localhost:8080/auth/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
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

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
```

## 🐳 Docker Services

The `docker-compose.yml` sets up:

1. **PostgreSQL Database**
   - Container: `postgres-project`
   - Port: `5432`
   - Persistent data storage
   - Health checks enabled

2. **Redis Cache**
   - Container: `redis-project`
   - Port: `6379`
   - Used for reset code storage
   - 15-minute code expiration

3. **Adminer** (Database Admin UI)
   - Container: `adminer`
   - Port: `9090`
   - Access: http://localhost:9090

## 🔒 Security Features

### Password Security
- **bcrypt Hashing**: All passwords are hashed using bcrypt with salt
- **Minimum Length**: 6 characters minimum (configurable)
- **No Plain Text**: Passwords are never stored in plain text

### JWT Authentication
- **Access Tokens**: Short-lived tokens (24 hours)
- **Refresh Tokens**: Long-lived tokens for token renewal
- **Secure Headers**: Proper JWT header validation
- **Claims Validation**: User ID and email verification

### Reset Code Security
- **6-Digit Codes**: Cryptographically secure random generation
- **Time Expiration**: 15-minute automatic expiration
- **Single Use**: Codes are deleted after successful use
- **Email Verification**: Codes sent only to registered emails

## 📧 Email Templates

The system includes professional HTML email templates for:

### Password Reset Email
- Clean, responsive design
- Large, readable reset code
- Clear expiration notice
- Security warning for unauthorized requests

## 🎨 Frontend Features

### Responsive Design
- Mobile-first responsive design
- Modern CSS with CSS Grid and Flexbox
- Professional color scheme and typography
- Loading states and animations

### User Experience
- **Form Validation**: Client-side and server-side validation
- **Loading Indicators**: Visual feedback during requests
- **Error Handling**: Clear error messages and recovery options
- **Success States**: Confirmation messages and redirects
- **Auto-Fill**: Email persistence during password reset flow

### Pages Features
- **Landing Page**: Clean introduction with navigation
- **Registration**: Complete signup with bio field
- **Login**: Secure authentication with forgot password link
- **Dashboard**: Protected user profile display
- **Forgot Password**: Email submission with validation
- **Reset Password**: Code entry with password confirmation

## 🧪 Testing the Application

### Backend Health Check
```bash
curl http://localhost:8080/health
```

### Frontend Access
- **Main Application**: http://localhost:8080
- **Database Admin**: http://localhost:9090
- **Registration**: http://localhost:8080/register.html
- **Login**: http://localhost:8080/login.html

### Complete Authentication Flow Test
1. Register a new user via frontend or API
2. Log in to receive JWT tokens
3. Access protected profile endpoint
4. Test forgot password flow
5. Verify email reset code functionality
6. Complete password reset process

<!-- 
## 🚀 Production Deployment

### Environment Setup
1. **Secure JWT Secret**: Use a cryptographically secure random string
2. **Database Credentials**: Use strong, unique credentials
3. **SSL/TLS**: Enable HTTPS for all endpoints
4. **Environment Variables**: Never commit secrets to version control

### Performance Optimizations
1. **Connection Pooling**: Configure appropriate database pool sizes
2. **Redis Persistence**: Configure Redis persistence for production
3. **Logging**: Set appropriate log levels
4. **Rate Limiting**: Implement API rate limiting
5. **CORS**: Configure CORS for your domain

### Monitoring
- Health check endpoint for load balancers
- Structured logging for monitoring tools
- Database connection monitoring
- Redis connection monitoring
- Email service monitoring -->

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines
- Follow Go best practices and conventions
- Add tests for new features
- Update documentation for API changes
- Ensure security best practices
- Test email functionality with real SMTP service

## 📝 License

This project is open source and available under the [MIT License](LICENSE).

## 🆘 Troubleshooting

### Common Issues

1. **Database connection failed**
   - Ensure Docker is running: `docker compose ps`
   - Check database credentials in `.env`
   - Verify PostgreSQL container health

2. **Redis connection failed**
   - Ensure Redis container is running
   - Check Redis configuration in `.env`
   - Verify Redis container health

3. **Email not sending**
   - Verify `RESEND_API_KEY` in `.env`
   - Check Resend account status and credits
   - Verify `EMAIL_FROM` domain is configured in Resend

4. **JWT token issues**
   - Ensure `JWT_SECRET` is set in `.env`
   - Check token expiration settings
   - Verify Authorization header format: `Bearer <token>`

5. **Frontend not loading**
   - Check if static files are being served correctly
   - Verify port configuration
   - Check browser console for JavaScript errors

6. **Password reset code not working**
   - Check Redis connection and data persistence
   - Verify code hasn't expired (15 minutes)
   - Ensure email matches exactly

### Logging and Debugging
- Application logs show detailed error information
- Use `make logs` to view Docker container logs
- Enable debug logging by setting log level in environment

### Getting Help

If you encounter any issues:
1. Check the application logs for detailed error messages
2. Verify all environment variables are properly set
3. Ensure all Docker containers are running and healthy
4. Test individual components (database, Redis, email service)
5. Review the API documentation and test with curl commands

----------