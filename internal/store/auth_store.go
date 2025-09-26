package store

import (
	"database/sql"
	"time"
)

type User struct {
    ID           int64     `json:"id" db:"id"`
    Username     string    `json:"username" db:"username"`
    Email        string    `json:"email" db:"email"`
    PasswordHash string    `json:"-" db:"password_hash"`
    Bio          *string   `json:"bio" db:"bio"`
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
    UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type RegisterRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
    Bio      string `json:"bio,omitempty"`
}

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type AuthResponse struct {
    Token string `json:"token"`
    User  User   `json:"user"`
}

type PostgresAuthStore struct {
	db *sql.DB
}

func NewPostgresAuthStore(db *sql.DB) *PostgresAuthStore {
	return &PostgresAuthStore{db: db}
}


type AuthStore interface {
	CreateUser(user *User) error
	GetUserByUsername(username string) (*User, error)
	GetUserByID(id int64) (*User, error)
}

func (s *PostgresAuthStore) CreateUser(user *User) error {
	query := `INSERT INTO users (username, email, password_hash, bio)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, updated_at`
	err := s.db.QueryRow(query, user.Username, user.Email, user.PasswordHash, user.Bio).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	
	if err != nil{
		return err
	}
	return nil
}

func (s *PostgresAuthStore) GetUserByUsername(username string) (*User, error) {
    user := User{}
    query := `
        SELECT id, username, email, password_hash, bio, created_at, updated_at
        FROM users
        WHERE username = $1`

    err := s.db.QueryRow(query, username).Scan(
        &user.ID, &user.Username, &user.Email, &user.PasswordHash,
        &user.Bio, &user.CreatedAt, &user.UpdatedAt,
    )

    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        return nil, err
    }

    return &user, nil
}

func (s *PostgresAuthStore) GetUserByID(id int64) (*User, error) {
    user := User{}
    query := `
        SELECT id, username, email, password_hash, bio, created_at, updated_at
        FROM users
        WHERE id = $1`

    err := s.db.QueryRow(query, id).Scan(
        &user.ID, &user.Username, &user.Email, &user.PasswordHash,
        &user.Bio, &user.CreatedAt, &user.UpdatedAt,
    )

    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        return nil, err
    }

    return &user, nil
}